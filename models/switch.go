package models

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)


func NewConfigFromYaml(yamlFile string) (*Config, error) {
	config := &Config{}
	f, err := os.Open(yamlFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

type Config struct {
	Hostname   string   `yaml:"hostname"`
	Metadata   Metadata `yaml:"metadata"`
	IsDeployed bool     `yaml:"is_deployed"`
	// RouterBgp                    RouterBgp              `yaml:"router_bgp"`
	// AddressFamilyIpv4            AddressFamilyIpv4      `yaml:"address_family_ipv4"`
	// AddressFamilyEvpn            AddressFamilyEvpn      `yaml:"address_family_evpn"`
	// Vrfs                         []Vrf                  `yaml:"vrfs"`
	Vlans []Vlan `yaml:"vlans"`
	// StaticRoutes                 []StaticRoute          `yaml:"static_routes"`
	ServiceRoutingProtocolsModel string `yaml:"service_routing_protocols_model"`
	IpRouting                    bool   `yaml:"ip_routing"`
	// DaemonTerminattr             DaemonTerminattr       `yaml:"daemon_terminattr"`
	// VlanInternalOrder            VlanInternalOrder      `yaml:"vlan_internal_order"`
	// IpNameServers                []IpNameServer         `yaml:"ip_name_servers"`
	// SpanningTree                 SpanningTree           `yaml:"spanning_tree"`
	// LocalUsers                   []LocalUser            `yaml:"local_users"`
	// ManagementInterfaces         []ManagementInterface  `yaml:"management_interfaces"`
	// ManagementApiHttp            ManagementApiHttp      `yaml:"management_api_http"`
	VlanInterfaces               []VlanInterface        `yaml:"vlan_interfaces"`
	// PortChannelInterfaces        []PortChannelInterface `yaml:"port_channel_interfaces"`
	// EthernetInterfaces           []EthernetInterface    `yaml:"ethernet_interfaces"`
	// MlagConfiguration            MlagConfiguration      `yaml:"mlag_configuration"`
	// RouteMaps                    []RouteMap             `yaml:"route_maps"`
	// LoopbackInterfaces           []LoopbackInterface    `yaml:"loopback_interfaces"`
	// PrefixLists                  []PrefixList           `yaml:"prefix_lists"`
	// RouterBfd                    RouterBfd              `yaml:"router_bfd"`
	// IpIgmpSnooping               IpIgmpSnooping         `yaml:"ip_igmp_snooping"`
	IpVirtualRouterMacAddress string `yaml:"ip_virtual_router_mac_address"`
	// VxlanInterface               VxlanInterface         `yaml:"vxlan_interface"`
	// VirtualSourceNatVrfs         []VirtualSourceNatVrf  `yaml:"virtual_source_nat_vrfs"`
	// Ntp                          Ntp                    `yaml:"ntp"`
}

// Define other structs like Metadata, RouterBgp, AddressFamilyIpv4, etc., here.
// Each of these structs should correspond to the structure of the YAML fields.

// Example of a nested struct
type Metadata struct {
	Platform string `yaml:"platform"`
}

type Vlan struct {
	ID                 int          `yaml:"id"`
	Tenant             string       `yaml:"tenant"`
	Rd                 string       `yaml:"rd"`
	RouteTargets       RouteTargets `yaml:"route_targets"`
	RedistributeRoutes []string     `yaml:"redistribute_routes"`
	Name               string       `yaml:"name"`
	TrunkGroups        []string     `yaml:"trunk_groups"`
}

type RouteTargets struct {
	Both []RouteTarget `yaml:"both"`
}

type RouteTarget struct {
	RouteTarget string `yaml:"route_target"`
}

type VlanInterface struct {
	Name             string `yaml:"name"`
	Description      string `yaml:"description"`
	Shutdown         bool   `yaml:"shutdown"`
	Mtu              int    `yaml:"mtu"`
	IpAddress        string `yaml:"ip_address"`
	IpAddressVirtual string `yaml:"ip_address_virtual"`
	Vrf              string `yaml:"vrf"`
	NoAutostate      bool   `yaml:"no_autostate"`
	Tenant           string `yaml:"tenant"`
	Type             string `yaml:"type"`
}
