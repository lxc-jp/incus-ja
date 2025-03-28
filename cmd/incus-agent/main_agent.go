package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"

	"github.com/lxc/incus/v6/internal/linux"
	"github.com/lxc/incus/v6/internal/server/instance/instancetype"
	"github.com/lxc/incus/v6/shared/logger"
	"github.com/lxc/incus/v6/shared/subprocess"
	"github.com/lxc/incus/v6/shared/util"
)

var (
	servers = make(map[string]*http.Server, 2)
	errChan = make(chan error)
)

type cmdAgent struct {
	global *cmdGlobal
}

func (c *cmdAgent) Command() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = "incus-agent [--debug]"
	cmd.Short = "Incus virtual machine agent"
	cmd.Long = `Description:
  Incus virtual machine agent

  This daemon is to be run inside virtual machines managed by Incus.
  It will normally be started through init scripts present or injected
  into the virtual machine.
`
	cmd.RunE = c.Run

	return cmd
}

func (c *cmdAgent) Run(cmd *cobra.Command, args []string) error {
	// Setup logger.
	err := logger.InitLogger("", "", c.global.flagLogVerbose, c.global.flagLogDebug, nil)
	if err != nil {
		os.Exit(1)
	}

	logger.Info("Starting")
	defer logger.Info("Stopped")

	// Apply the templated files.
	files, err := templatesApply("files/")
	if err != nil {
		return err
	}

	// Sync the hostname.
	if util.PathExists("/proc/sys/kernel/hostname") && slices.Contains(files, "/etc/hostname") {
		// Open the two files.
		src, err := os.Open("/etc/hostname")
		if err != nil {
			return err
		}

		dst, err := os.Create("/proc/sys/kernel/hostname")
		if err != nil {
			return err
		}

		// Copy the data.
		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}

		// Close the files.
		_ = src.Close()
		err = dst.Close()
		if err != nil {
			return err
		}
	}

	// Run cloud-init.
	if util.PathExists("/etc/cloud") && slices.Contains(files, "/var/lib/cloud/seed/nocloud-net/meta-data") {
		logger.Info("Seeding cloud-init")

		cloudInitPath := "/run/cloud-init"
		if util.PathExists(cloudInitPath) {
			logger.Info(fmt.Sprintf("Removing %q", cloudInitPath))
			err = os.RemoveAll(cloudInitPath)
			if err != nil {
				return err
			}
		}

		logger.Info("Rebooting")
		_, _ = subprocess.RunCommand("reboot")

		// Wait up to 5min for the reboot to actually happen, if it doesn't, then move on to allowing connections.
		time.Sleep(300 * time.Second)
	}

	reconfigureNetworkInterfaces()

	// Load the kernel driver.
	if !util.PathExists("/dev/vsock") {
		logger.Info("Loading vsock module")

		err = linux.LoadModule("vsock")
		if err != nil {
			return fmt.Errorf("Unable to load the vsock kernel module: %w", err)
		}

		// Wait for vsock device to appear.
		for i := 0; i < 5; i++ {
			if !util.PathExists("/dev/vsock") {
				time.Sleep(1 * time.Second)
			}
		}
	}

	// Mount shares from host.
	c.mountHostShares()

	d := newDaemon(c.global.flagLogDebug, c.global.flagLogVerbose)

	// Start the server.
	err = startHTTPServer(d, c.global.flagLogDebug)
	if err != nil {
		return fmt.Errorf("Failed to start HTTP server: %w", err)
	}

	// Check whether we should start the DevIncus server in the early setup. This way, /dev/incus/sock
	// will be available for any systemd services starting after the agent.
	if util.PathExists("agent.conf") {
		f, err := os.Open("agent.conf")
		if err != nil {
			return err
		}

		err = setConnectionInfo(d, f)
		if err != nil {
			_ = f.Close()
			return err
		}

		_ = f.Close()

		if d.DevIncusEnabled {
			err = startDevIncusServer(d)
			if err != nil {
				return err
			}
		}
	}

	// Create a cancellation context.
	ctx, cancelFunc := context.WithCancel(context.Background())

	// Start status notifier in background.
	cancelStatusNotifier := c.startStatusNotifier(ctx, d.chConnected)

	// Done with early setup, tell systemd to continue boot.
	// Allows a service that needs a file that's generated by the agent to be able to declare After=incus-agent
	// and know the file will have been created by the time the service is started.
	if os.Getenv("NOTIFY_SOCKET") != "" {
		_, err := subprocess.RunCommand("systemd-notify", "READY=1")
		if err != nil {
			cancelStatusNotifier() // Ensure STOPPED status is written to QEMU status ringbuffer.
			cancelFunc()

			return fmt.Errorf("Failed to notify systemd of readiness: %w", err)
		}
	}

	// Cancel context when SIGTEM is received.
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, unix.SIGTERM)

	exitStatus := 0

	select {
	case <-chSignal:
	case err := <-errChan:
		fmt.Fprintln(os.Stderr, err)
		exitStatus = 1
	}

	cancelStatusNotifier() // Ensure STOPPED status is written to QEMU status ringbuffer.
	cancelFunc()

	os.Exit(exitStatus)

	return nil
}

// startStatusNotifier sends status of agent to vserial ring buffer every 10s or when context is done.
// Returns a function that can be used to update the running status to STOPPED in the ring buffer.
func (c *cmdAgent) startStatusNotifier(ctx context.Context, chConnected <-chan struct{}) context.CancelFunc {
	// Write initial started status.
	_ = c.writeStatus("STARTED")

	wg := sync.WaitGroup{}
	exitCtx, exit := context.WithCancel(ctx) // Allows manual synchronous cancellation via cancel function.
	cancel := func() {
		exit()    // Signal for the go routine to end.
		wg.Wait() // Wait for the go routine to actually finish.
	}

	wg.Add(1)
	go func() {
		defer wg.Done() // Signal to cancel function that we are done.

		ticker := time.NewTicker(time.Duration(time.Second) * 5)
		defer ticker.Stop()

		for {
			select {
			case <-chConnected:
				_ = c.writeStatus("CONNECTED") // Indicate we were able to connect.
			case <-ticker.C:
				_ = c.writeStatus("STARTED") // Re-populate status periodically in case the daemon restarts.
			case <-exitCtx.Done():
				_ = c.writeStatus("STOPPED") // Indicate we are stopping and exit go routine.
				return
			}
		}
	}()

	return cancel
}

// writeStatus writes a status code to the vserial ring buffer used to detect agent status on host.
func (c *cmdAgent) writeStatus(status string) error {
	if util.PathExists("/dev/virtio-ports/org.linuxcontainers.incus") {
		vSerial, err := os.OpenFile("/dev/virtio-ports/org.linuxcontainers.incus", os.O_RDWR, 0o600)
		if err != nil {
			return err
		}

		defer vSerial.Close()

		_, err = vSerial.Write([]byte(fmt.Sprintf("%s\n", status)))
		if err != nil {
			return err
		}
	}

	return nil
}

// mountHostShares reads the agent-mounts.json file from config share and mounts the shares requested.
func (c *cmdAgent) mountHostShares() {
	agentMountsFile := "./agent-mounts.json"
	if !util.PathExists(agentMountsFile) {
		return
	}

	b, err := os.ReadFile(agentMountsFile)
	if err != nil {
		logger.Errorf("Failed to load agent mounts file %q: %v", agentMountsFile, err)
	}

	var agentMounts []instancetype.VMAgentMount
	err = json.Unmarshal(b, &agentMounts)
	if err != nil {
		logger.Errorf("Failed to parse agent mounts file %q: %v", agentMountsFile, err)
		return
	}

	for _, mount := range agentMounts {
		if !slices.Contains([]string{"9p", "virtiofs"}, mount.FSType) {
			logger.Infof("Unsupported mount fstype %q", mount.FSType)
			continue
		}

		err = tryMountShared(mount.Source, mount.Target, mount.FSType, mount.Options)
		if err != nil {
			logger.Infof("Failed to mount %q (Type: %q, Options: %v) to %q: %v", mount.Source, "virtiofs", mount.Options, mount.Target, err)
			continue
		}

		logger.Infof("Mounted %q (Type: %q, Options: %v) to %q", mount.Source, mount.FSType, mount.Options, mount.Target)
	}
}

func tryMountShared(src string, dst string, fstype string, opts []string) error {
	// Convert relative mounts to absolute from / otherwise dir creation fails or mount fails.
	if !strings.HasPrefix(dst, "/") {
		dst = fmt.Sprintf("/%s", dst)
	}

	// Check mount path.
	if !util.PathExists(dst) {
		// Create the mount path.
		err := os.MkdirAll(dst, 0o755)
		if err != nil {
			return fmt.Errorf("Failed to create mount target %q", dst)
		}
	} else if linux.IsMountPoint(dst) {
		// Already mounted.
		return nil
	}

	// Prepare the arguments.
	sharedArgs := []string{}
	p9Args := []string{}

	for _, opt := range opts {
		// transport and msize mount option are specific to 9p.
		if strings.HasPrefix(opt, "trans=") || strings.HasPrefix(opt, "msize=") {
			p9Args = append(p9Args, "-o", opt)
			continue
		}

		sharedArgs = append(sharedArgs, "-o", opt)
	}

	// Always try virtiofs first.
	args := []string{"-t", "virtiofs", src, dst}
	args = append(args, sharedArgs...)

	_, err := subprocess.RunCommand("mount", args...)
	if err == nil {
		return nil
	} else if fstype == "virtiofs" {
		return err
	}

	// Then fallback to 9p.
	args = []string{"-t", "9p", src, dst}
	args = append(args, sharedArgs...)
	args = append(args, p9Args...)

	_, err = subprocess.RunCommand("mount", args...)
	if err != nil {
		return err
	}

	return nil
}
