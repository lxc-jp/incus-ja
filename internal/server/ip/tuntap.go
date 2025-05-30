package ip

import (
	"github.com/lxc/incus/v6/shared/subprocess"
)

// Tuntap represents arguments for tuntap manipulation.
type Tuntap struct {
	Name       string
	Mode       string
	MultiQueue bool
	Master     string
}

// Add adds new tuntap interface.
func (t *Tuntap) Add() error {
	cmd := []string{"tuntap", "add", "name", t.Name, "mode", t.Mode}
	if t.MultiQueue {
		cmd = append(cmd, "multi_queue")
	}

	_, err := subprocess.RunCommand("ip", cmd...)
	if err != nil {
		return err
	}

	if t.Master != "" {
		_, err := subprocess.RunCommand("ip", "link", "set", t.Name, "master", t.Master)
		if err != nil {
			return err
		}
	}

	return nil
}
