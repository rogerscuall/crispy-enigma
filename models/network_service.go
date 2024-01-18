package models

import (
	"net/netip"
	"strconv"
	"strings"
)

type Tenant struct {
	Name          string   `yaml:"name"`
	MacVrfVniBase int      `yaml:"mac_vrf_vni_base"`
	VRFs          []VRF    `yaml:"vrfs"`
	L2Vlans       []L2Vlan `yaml:"l2vlans"`
}

type VRF struct {
	Name           string         `yaml:"name"`
	VrfVni         int            `yaml:"vrf_vni"`
	VtepDiagnostic VtepDiagnostic `yaml:"vtep_diagnostic"`
	SVIs           []SVI          `yaml:"svis"`
}

type VtepDiagnostic struct {
	Loopback        int    `yaml:"loopback"`
	LoopbackIpRange string `yaml:"loopback_ip_range"`
}

type SVI struct {
	ID               int    `yaml:"id"`
	Name             string `yaml:"name"`
	Enabled          bool   `yaml:"enabled"`
	IpAddressVirtual string `yaml:"ip_address_virtual"`
}

type L2Vlan struct {
	ID   int    `yaml:"id"`
	Name string `yaml:"name"`
}

type NetworkService struct {
	Tenants []Tenant `yaml:"tenants"`
}

func NewSVI(subnet string) *SVI {
	// get the first address
	prefix := netip.MustParsePrefix(subnet)
	first := prefix.Addr()
	octects := strings.Split(subnet, ".")
	id, _ := strconv.Atoi(octects[2])
	name := "VLAN_" + octects[2]
	return &SVI{
		IpAddressVirtual: first.Next().String() + "/" + strconv.Itoa(prefix.Bits()),
		Enabled:          true,
		Name:             name,
		ID:               id,
	}
}
