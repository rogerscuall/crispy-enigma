package act

import (
	"fmt"
	"net"

	"github.com/c-robinson/iplib"
	mo "github.com/rogerscuall/crispy-enigma/models"
)

// Defines an ACT topology configuration in yaml format
type TopologyConfig struct {
	Veos       DeviceDetails `yaml:"veos"`
	Cvp        DeviceDetails `yaml:"cvp"`
	Generic    DeviceDetails `yaml:"generic"`
	ThirdParty DeviceDetails `yaml:"third-party"`
	Nodes      []*Node       `yaml:"nodes"`
	Links      []*Link       `yaml:"links"`
}

type Link struct {
	Connection []string `yaml:"connection"`
}

// AddNode adds a new node to the TopologyConfig
func (a *TopologyConfig) AddNode(nodeName, ipAddr string) error {
	// Check if the node already exists
	for _, n := range a.Nodes {
		if n.Name == nodeName {
			return fmt.Errorf("Node %s already exists", nodeName)
		}
	}
	a.Nodes = append(a.Nodes, &Node{Name: nodeName, IPAddr: ipAddr, NodeType: "veos"})
	return nil
}

/*
	 AddLinksToNodes adds links to the nodes in the ACTTopologyConfig
		linksMap should have this format:
		{"node1": ["node2:port1", "node3:port2"], "node2": ["node1:port1"]}
*/
func (c *TopologyConfig) AddLinksToNodes(network mo.Network) {
	networkConnections := network.CleanInNetworkConnections()
	for _, connection := range networkConnections {
		side1 := fmt.Sprint(connection.SideA, ":", connection.PortA)
		side2 := fmt.Sprint(connection.SideB, ":", connection.PortB)
		connection := Link{
			Connection: []string{side1, side2},
		}
		c.Links = append(c.Links, &connection)
	}
}

// AddIPToHosts adds a new node to the ACTTopologyConfig
func (c *TopologyConfig) AddIPToHosts(hostnames []string, firstIP string) {
	ipInit := net.ParseIP(firstIP)
	ipNext := iplib.NextIP(ipInit)
	for _, hostname := range hostnames {
		if err := c.AddNode(hostname, ipNext.String()); err != nil {
			fmt.Println(err)
		}
		ipNext = iplib.NextIP(ipNext)
	}
}

func (c *TopologyConfig) AddPortsToNodes(network mo.Network) {
	for _, node := range c.Nodes {
		for _, config := range network.Configs {
			if node.Name == config.Hostname {
				for _, e := range config.EthernetInterfaces {
					node.Ports = append(node.Ports, e.Name)
				}
			}
		}
	}
}

type MetaData struct {
	HostVars map[string]HostVars `json:"hostvars"`
}

// Define struct to represent the JSON data
type InventoryData struct {
	Meta MetaData `json:"_meta"`
}

func (i *InventoryData) GetHost() []string {
	var hostnames []string
	fmt.Println("Hosts information:")
	for hostname, vars := range i.Meta.HostVars {
		fmt.Printf("Hostname: %s\n", hostname)
		for key, value := range vars {
			fmt.Printf("  %s: %v\n", key, value)
		}
		hostnames = append(hostnames, hostname)
	}
	return hostnames
}

type HostVars map[string]interface{}

type DeviceDetails struct {
	Password string `yaml:"password"`
	Username string `yaml:"username"`
	Version  string `yaml:"version"`
}



func (n Node) MarshalYAML() (interface{}, error) {
	nodeDetails := map[string]interface{}{}

	// Only add non-empty fields to the nodeDetails map
	if n.IPAddr != "" {
		nodeDetails["ip_addr"] = n.IPAddr
	}
	if n.NodeType != "" {
		nodeDetails["node_type"] = n.NodeType
	}
	if n.AutoConfiguration {
		nodeDetails["auto_configuration"] = n.AutoConfiguration
	}
	if len(n.Ports) > 0 {
		nodeDetails["ports"] = n.Ports
	}

	// Return a map with the node name as the key, only if nodeDetails is not empty
	if len(nodeDetails) > 0 {
		return map[string]interface{}{
			n.Name: nodeDetails,
		}, nil
	}

	return nil, nil // Avoid including an empty map
}