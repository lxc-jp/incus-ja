package ovs

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"

	ovsdbClient "github.com/ovn-org/libovsdb/client"
	ovsdbModel "github.com/ovn-org/libovsdb/model"
	"github.com/ovn-org/libovsdb/ovsdb"

	"github.com/lxc/incus/internal/server/ip"
	ovsSwitch "github.com/lxc/incus/internal/server/network/ovs/schema/ovs"
	"github.com/lxc/incus/shared/subprocess"
	"github.com/lxc/incus/shared/util"
)

// ovnBridgeMappingMutex locks access to read/write external-ids:ovn-bridge-mappings.
var ovnBridgeMappingMutex sync.Mutex

// Installed returns true if the OVS tools are installed.
func (o *VSwitch) Installed() bool {
	_, err := exec.LookPath("ovs-vsctl")
	return err == nil
}

// BridgeExists returns true if the bridge exists.
func (o *VSwitch) BridgeExists(bridgeName string) (bool, error) {
	ctx := context.TODO()
	bridge := &ovsSwitch.Bridge{Name: bridgeName}

	err := o.client.Get(ctx, bridge)
	if err != nil {
		if err == ovsdbClient.ErrNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// BridgeAdd adds a new bridge.
func (o *VSwitch) BridgeAdd(bridgeName string, mayExist bool, hwaddr net.HardwareAddr, mtu uint32) error {
	ctx := context.TODO()

	// Create interface.
	iface := ovsSwitch.Interface{
		UUID: "interface",
		Name: bridgeName,
	}

	if mtu > 0 {
		mtu := int(mtu)
		iface.MTURequest = &mtu
	}

	interfaceOps, err := o.client.Create(&iface)
	if err != nil {
		return err
	}

	// Create port.
	port := ovsSwitch.Port{
		UUID:       "port",
		Name:       bridgeName,
		Interfaces: []string{iface.UUID},
	}

	portOps, err := o.client.Create(&port)
	if err != nil {
		return err
	}

	// Create bridge.
	bridge := ovsSwitch.Bridge{
		UUID:  "bridge",
		Name:  bridgeName,
		Ports: []string{port.UUID},
	}

	if hwaddr != nil {
		bridge.OtherConfig = map[string]string{"hwaddr": hwaddr.String()}
	}

	bridgeOps, err := o.client.Create(&bridge)
	if err != nil {
		return err
	}

	if mayExist {
		err = o.client.Get(ctx, &bridge)
		if err != nil && err != ovsdbClient.ErrNotFound {
			return err
		}

		if bridge.UUID != "bridge" {
			// Bridge already exists.
			return nil
		}
	}

	// Create switch entry.
	ovsRow := ovsSwitch.OpenvSwitch{
		UUID: o.rootUUID,
	}

	mutateOps, err := o.client.Where(&ovsRow).Mutate(&ovsRow, ovsdbModel.Mutation{
		Field:   &ovsRow.Bridges,
		Mutator: ovsdb.MutateOperationInsert,
		Value:   []string{bridge.UUID},
	})
	if err != nil {
		return err
	}

	operations := append(interfaceOps, portOps...)
	operations = append(operations, bridgeOps...)
	operations = append(operations, mutateOps...)

	resp, err := o.client.Transact(ctx, operations...)
	if err != nil {
		return err
	}

	_, err = ovsdb.CheckOperationResults(resp, operations)
	if err != nil {
		return err
	}

	// Wait for kernel interface to appear.
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)

		if util.PathExists(fmt.Sprintf("/sys/class/net/%s", bridgeName)) {
			return nil
		}
	}

	return fmt.Errorf("Bridge interface failed to appear")
}

// BridgeDelete deletes a bridge.
func (o *VSwitch) BridgeDelete(bridgeName string) error {
	ctx := context.TODO()

	bridge := ovsSwitch.Bridge{
		Name: bridgeName,
	}

	err := o.client.Get(ctx, &bridge)
	if err != nil {
		return err
	}

	ovsRow := ovsSwitch.OpenvSwitch{
		UUID: o.rootUUID,
	}

	operations, err := o.client.Where(&ovsRow).Mutate(&ovsRow, ovsdbModel.Mutation{
		Field:   &ovsRow.Bridges,
		Mutator: "delete",
		Value:   []string{bridge.UUID},
	})
	if err != nil {
		return err
	}

	resp, err := o.client.Transact(ctx, operations...)
	if err != nil {
		return err
	}

	_, err = ovsdb.CheckOperationResults(resp, operations)
	if err != nil {
		return err
	}

	return nil
}

// BridgePortAdd adds a port to the bridge (if already attached does nothing).
func (o *VSwitch) BridgePortAdd(bridgeName string, portName string, mayExist bool) error {
	ctx := context.TODO()

	// Get the bridge.
	bridge := ovsSwitch.Bridge{
		Name: bridgeName,
	}

	err := o.client.Get(ctx, &bridge)
	if err != nil {
		return err
	}

	// Create the interface.
	iface := ovsSwitch.Interface{
		UUID: "interface",
		Name: portName,
	}

	interfaceOps, err := o.client.Create(&iface)
	if err != nil {
		return err
	}

	// Create the port.
	port := ovsSwitch.Port{
		Name: portName,
	}

	err = o.client.Get(ctx, &port)
	if err != nil && err != ovsdbClient.ErrNotFound {
		return err
	}

	if port.UUID != "" {
		if mayExist {
			// Already exists.
			return nil
		}

		return fmt.Errorf("OVS port %q already exists on %q", portName, bridgeName)
	}

	port.UUID = "port"
	port.Interfaces = []string{iface.UUID}
	portOps, err := o.client.Create(&port)
	if err != nil {
		return err
	}

	// Create the bridge port entry.
	mutateOps, err := o.client.Where(&bridge).Mutate(&bridge, ovsdbModel.Mutation{
		Field:   &bridge.Ports,
		Mutator: ovsdb.MutateOperationInsert,
		Value:   []string{port.UUID},
	})
	if err != nil {
		return err
	}

	operations := append(interfaceOps, portOps...)
	operations = append(operations, mutateOps...)

	resp, err := o.client.Transact(ctx, operations...)
	if err != nil {
		return err
	}

	_, err = ovsdb.CheckOperationResults(resp, operations)
	if err != nil {
		return err
	}

	return nil
}

// BridgePortDelete deletes a port from the bridge (if already detached does nothing).
func (o *VSwitch) BridgePortDelete(bridgeName string, portName string) error {
	_, err := subprocess.RunCommand("ovs-vsctl", "--if-exists", "del-port", bridgeName, portName)
	if err != nil {
		return err
	}

	return nil
}

// BridgePortSet sets port options.
func (o *VSwitch) BridgePortSet(portName string, options ...string) error {
	_, err := subprocess.RunCommand("ovs-vsctl", append([]string{"set", "port", portName}, options...)...)
	if err != nil {
		return err
	}

	return nil
}

// InterfaceAssociateOVNSwitchPort removes any existing switch ports associated to the specified ovnSwitchPortName
// and then associates the specified interfaceName to the OVN switch port.
func (o *VSwitch) InterfaceAssociateOVNSwitchPort(interfaceName string, ovnSwitchPortName string) error {
	// Clear existing ports that were formerly associated to ovnSwitchPortName.
	existingPorts, err := subprocess.RunCommand("ovs-vsctl", "--format=csv", "--no-headings", "--data=bare", "--colum=name", "find", "interface", fmt.Sprintf("external-ids:iface-id=%s", string(ovnSwitchPortName)))
	if err != nil {
		return err
	}

	existingPorts = strings.TrimSpace(existingPorts)
	if existingPorts != "" {
		for _, port := range strings.Split(existingPorts, "\n") {
			_, err = subprocess.RunCommand("ovs-vsctl", "del-port", port)
			if err != nil {
				return err
			}

			// Atempt to remove port, but don't fail if doesn't exist or can't be removed, at least
			// the switch association has been successfully removed, so the new port being added next
			// won't fail to work properly.
			link := &ip.Link{Name: port}
			_ = link.Delete()
		}
	}

	_, err = subprocess.RunCommand("ovs-vsctl", "set", "interface", interfaceName, fmt.Sprintf("external_ids:iface-id=%s", string(ovnSwitchPortName)))
	if err != nil {
		return err
	}

	return nil
}

// InterfaceAssociatedOVNSwitchPort returns the OVN switch port associated to the interface.
func (o *VSwitch) InterfaceAssociatedOVNSwitchPort(interfaceName string) (string, error) {
	ovnSwitchPort, err := subprocess.RunCommand("ovs-vsctl", "get", "interface", interfaceName, "external_ids:iface-id")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(ovnSwitchPort), nil
}

// ChassisID returns the local chassis ID.
func (o *VSwitch) ChassisID() (string, error) {
	ctx := context.TODO()

	vSwitch := &ovsSwitch.OpenvSwitch{
		UUID: o.rootUUID,
	}

	err := o.client.Get(ctx, vSwitch)
	if err != nil {
		return "", err
	}

	val := vSwitch.ExternalIDs["system-id"]
	return val, nil
}

// OVNEncapIP returns the enscapsulation IP used for OVN underlay tunnels.
func (o *VSwitch) OVNEncapIP() (net.IP, error) {
	// ovs-vsctl's get command doesn't support its --format flag, so we always get the output quoted.
	// However ovs-vsctl's find and list commands don't support retrieving a single column's map field.
	// And ovs-vsctl's JSON output is unfriendly towards statically typed languages as it mixes data types
	// in a slice. So stick with "get" command and use Go's strconv.Unquote to return the actual values.
	encapIPStr, err := subprocess.RunCommand("ovs-vsctl", "get", "open_vswitch", ".", "external_ids:ovn-encap-ip")
	if err != nil {
		return nil, err
	}

	encapIPStr = strings.TrimSpace(encapIPStr)
	encapIPStr, err = unquote(encapIPStr)
	if err != nil {
		return nil, fmt.Errorf("Failed unquoting: %w", err)
	}

	encapIP := net.ParseIP(encapIPStr)
	if encapIP == nil {
		return nil, fmt.Errorf("Invalid ovn-encap-ip address")
	}

	return encapIP, nil
}

// OVNBridgeMappings gets the current OVN bridge mappings.
func (o *VSwitch) OVNBridgeMappings(bridgeName string) ([]string, error) {
	ctx := context.TODO()

	vSwitch := &ovsSwitch.OpenvSwitch{
		UUID: o.rootUUID,
	}

	err := o.client.Get(ctx, vSwitch)
	if err != nil {
		return nil, err
	}

	val := vSwitch.ExternalIDs["ovn-bridge-mappings"]
	if val == "" {
		return []string{}, nil
	}

	return strings.SplitN(val, ",", -1), nil
}

// OVNBridgeMappingAdd appends an OVN bridge mapping between a bridge and the logical provider name.
func (o *VSwitch) OVNBridgeMappingAdd(bridgeName string, providerName string) error {
	ovnBridgeMappingMutex.Lock()
	defer ovnBridgeMappingMutex.Unlock()

	mappings, err := o.OVNBridgeMappings(bridgeName)
	if err != nil {
		return err
	}

	newMapping := fmt.Sprintf("%s:%s", providerName, bridgeName)
	for _, mapping := range mappings {
		if mapping == newMapping {
			return nil // Mapping is already present, nothing to do.
		}
	}

	mappings = append(mappings, newMapping)

	// Set new mapping string back into the database.
	_, err = subprocess.RunCommand("ovs-vsctl", "set", "open_vswitch", ".", fmt.Sprintf("external-ids:ovn-bridge-mappings=%s", strings.Join(mappings, ",")))
	if err != nil {
		return err
	}

	return nil
}

// OVNBridgeMappingDelete deletes an OVN bridge mapping between a bridge and the logical provider name.
func (o *VSwitch) OVNBridgeMappingDelete(bridgeName string, providerName string) error {
	ovnBridgeMappingMutex.Lock()
	defer ovnBridgeMappingMutex.Unlock()

	mappings, err := o.OVNBridgeMappings(bridgeName)
	if err != nil {
		return err
	}

	changed := false
	newMappings := make([]string, 0, len(mappings))
	matchMapping := fmt.Sprintf("%s:%s", providerName, bridgeName)
	for _, mapping := range mappings {
		if mapping != matchMapping {
			newMappings = append(newMappings, mapping)
		} else {
			changed = true
		}
	}

	if changed {
		if len(newMappings) < 1 {
			// Remove mapping key in the database.
			_, err = subprocess.RunCommand("ovs-vsctl", "remove", "open_vswitch", ".", "external-ids", "ovn-bridge-mappings")
			if err != nil {
				return err
			}
		} else {
			// Set updated mapping string back into the database.
			_, err = subprocess.RunCommand("ovs-vsctl", "set", "open_vswitch", ".", fmt.Sprintf("external-ids:ovn-bridge-mappings=%s", strings.Join(newMappings, ",")))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// BridgePortList returns a list of ports that are connected to the bridge.
func (o *VSwitch) BridgePortList(bridgeName string) ([]string, error) {
	// Clear existing ports that were formerly associated to ovnSwitchPortName.
	portString, err := subprocess.RunCommand("ovs-vsctl", "list-ports", bridgeName)
	if err != nil {
		return nil, err
	}

	ports := []string{}

	portString = strings.TrimSpace(portString)
	if portString != "" {
		for _, port := range strings.Split(portString, "\n") {
			ports = append(ports, strings.TrimSpace(port))
		}
	}

	return ports, nil
}

// HardwareOffloadingEnabled returns true if hardware offloading is enabled.
func (o *VSwitch) HardwareOffloadingEnabled() bool {
	// ovs-vsctl's get command doesn't support its --format flag, so we always get the output quoted.
	// However ovs-vsctl's find and list commands don't support retrieving a single column's map field.
	// And ovs-vsctl's JSON output is unfriendly towards statically typed languages as it mixes data types
	// in a slice. So stick with "get" command and use Go's strconv.Unquote to return the actual values.
	offload, err := subprocess.RunCommand("ovs-vsctl", "--if-exists", "get", "open_vswitch", ".", "other_config:hw-offload")
	if err != nil {
		return false
	}

	offload = strings.TrimSpace(offload)
	if offload == "" {
		return false
	}

	offload, err = unquote(offload)
	if err != nil {
		return false
	}

	return offload == "true"
}

// OVNSouthboundDBRemoteAddress gets the address of the southbound ovn database.
func (o *VSwitch) OVNSouthboundDBRemoteAddress() (string, error) {
	result, err := subprocess.RunCommand("ovs-vsctl", "get", "open_vswitch", ".", "external_ids:ovn-remote")
	if err != nil {
		return "", err
	}

	addr, err := unquote(strings.TrimSuffix(result, "\n"))
	if err != nil {
		return "", err
	}

	return addr, nil
}
