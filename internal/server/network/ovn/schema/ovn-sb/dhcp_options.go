// Code generated by "libovsdb.modelgen"
// DO NOT EDIT.

package ovsmodel

const DHCPOptionsTable = "DHCP_Options"

type (
	DHCPOptionsType = string
)

var (
	DHCPOptionsTypeBool         DHCPOptionsType = "bool"
	DHCPOptionsTypeUint8        DHCPOptionsType = "uint8"
	DHCPOptionsTypeUint16       DHCPOptionsType = "uint16"
	DHCPOptionsTypeUint32       DHCPOptionsType = "uint32"
	DHCPOptionsTypeIpv4         DHCPOptionsType = "ipv4"
	DHCPOptionsTypeStaticRoutes DHCPOptionsType = "static_routes"
	DHCPOptionsTypeStr          DHCPOptionsType = "str"
	DHCPOptionsTypeHostID       DHCPOptionsType = "host_id"
	DHCPOptionsTypeDomains      DHCPOptionsType = "domains"
)

// DHCPOptions defines an object in DHCP_Options table
type DHCPOptions struct {
	UUID string          `ovsdb:"_uuid"`
	Code int             `ovsdb:"code"`
	Name string          `ovsdb:"name"`
	Type DHCPOptionsType `ovsdb:"type"`
}
