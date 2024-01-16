package models

type Tenant struct {
	Name           string `yaml:"name"`
	MacVrfVniBase  int    `yaml:"mac_vrf_vni_base"`
	Vrfs           []VRF  `yaml:"vrfs"`
	L2Vlans        []L2Vlan `yaml:"l2vlans"`
}

type VRF struct {
	Name            string `yaml:"name"`
	VrfVni          int    `yaml:"vrf_vni"`
	VtepDiagnostic  VtepDiagnostic `yaml:"vtep_diagnostic"`
	Svis            []SVI  `yaml:"svis"`
}

type VtepDiagnostic struct {
	Loopback        int    `yaml:"loopback"`
	LoopbackIpRange string `yaml:"loopback_ip_range"`
}

type SVI struct {
	Id               int    `yaml:"id"`
	Name             string `yaml:"name"`
	Enabled          bool   `yaml:"enabled"`
	IpAddressVirtual string `yaml:"ip_address_virtual"`
}

type L2Vlan struct {
	Id   int    `yaml:"id"`
	Name string `yaml:"name"`
}

type NetworkService struct {
	Tenants []Tenant `yaml:"tenants"`
}
