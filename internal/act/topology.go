package act

import (
	"fmt"
	"net"
	"strings"

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

// AddNodeWithNewIP adds a new node to the TopologyConfig
func (a *TopologyConfig) AddNodeWithNewIP(nodeName, ipAddr string) error {
	// Check if the node already exists
	for _, n := range a.Nodes {
		if n.Name == nodeName {
			return fmt.Errorf("Node %s already exists", nodeName)
		}
	}
	a.Nodes = append(a.Nodes, &Node{Name: nodeName, IPAddr: ipAddr, NodeType: "veos"})
	return nil
}

func (a *TopologyConfig) AddNodes(network mo.Network) error {
	for _, config := range network.Configs {
		if len(config.ManagementInterfaces) == 0 {
			return fmt.Errorf("no management interface found for node %s", config.Hostname)
		}
		ip := config.ManagementInterfaces[0].IPAddress
		mgmtIP := strings.Split(ip, "/")
		node := Node{
			Name:     config.Hostname,
			IPAddr:   mgmtIP[0],
			NodeType: "veos",
		}
		a.Nodes = append(a.Nodes, &node)
	}
	return nil
}

/*
AddLinksToNodes adds links to the nodes in the ACTTopologyConfig.
Guarantees that the connections are unique, and terminates in another network device (InNetworkConnections).
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
		if err := c.AddNodeWithNewIP(hostname, ipNext.String()); err != nil {
			fmt.Println(err)
		}
		ipNext = iplib.NextIP(ipNext)
	}
}

// AddPortsToNodes adds the ports to the nodes in the ACTTopologyConfig
// All interfaces configured in the network.Configs are added to the nodes regardless of state or type.
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
