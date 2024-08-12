/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/apenella/go-ansible/v2/pkg/execute"
	"github.com/apenella/go-ansible/v2/pkg/inventory"
	"github.com/c-robinson/iplib"
	mo "github.com/rogerscuall/crispy-enigma/models"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

/*
1. Run the ansible-inventory command to get the inventory data. We dont need this anymore as we are reading from the structured configuration files.
2. Parse the JSON data
3. Extract the hostnames
4. Load the ACT Topology YAML data from a file. This ACT Topology has some data that needs to be kept, we are just adding new nodes and links to it.
5. Add the hostnames to the ACT Topology data.
6. Parse the information from the structured configuration files.
7. From the structured configuration data, parse the interface ethernet configuration as links struct. (DONE)
8. Identify what ports are available in the test environment, if the port does not connect to another switch is available to be used.
8.5. The available ports should be added to the test node as links in the ACT Topology data.
9. Add the links struct to the ACT Topology data for the "no available" ports. These are the ports that connect the switches.
10. Output the ACT Topology data to a file
*/

// actTopologyCmd represents the actTopology command
var actTopologyCmd = &cobra.Command{
	Use:   "actTopology",
	Short: "From AVD project creates an ACT topology",
	Long: `This command will create an ACT topology from and AVD project. 
It can use as based the intended configuration in structured format and add test nodes to the topology for data plane testing`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("actTopology called")
		actTopology()
	},
}

func init() {
	rootCmd.AddCommand(actTopologyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// actTopologyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// actTopologyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

type MetaData struct {
	HostVars map[string]HostVars `json:"hostvars"`
}

type HostVars map[string]interface{}

// Struct definitions
type ACTTopologyConfig struct {
	Veos       DeviceDetails `yaml:"veos"`
	Cvp        DeviceDetails `yaml:"cvp"`
	Generic    DeviceDetails `yaml:"generic"`
	ThirdParty DeviceDetails `yaml:"third-party"`
	Nodes      []*Node       `yaml:"nodes"`
	Links      []*Link       `yaml:"links"`
}

type DeviceDetails struct {
	Password string `yaml:"password"`
	Username string `yaml:"username"`
	Version  string `yaml:"version"`
}

type Node struct {
	Name     string   `yaml:"-"`
	IPAddr   string   `yaml:"ip_addr"`
	NodeType string   `yaml:"node_type"`
	Ports    []string `yaml:"ports"`
}

type Link struct {
	Connection []string `yaml:"connection"`
}

var (
	// defaultPorts = []string{"Ethernet1-32"}
	debug = false
)

// AddNode adds a new node to the ACTTopologyConfig
func (a *ACTTopologyConfig) AddNode(nodeName, ipAddr string) error {
	// Check if the node already exists
	for _, n := range a.Nodes {
		if n.Name == nodeName {
			return fmt.Errorf("Node %s already exists", nodeName)
		}
	}
	a.Nodes = append(a.Nodes, &Node{Name: nodeName, IPAddr: ipAddr, NodeType: "veos"})
	return nil
}

func (n *Node) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var nodeMap map[string]interface{}
	if err := unmarshal(&nodeMap); err != nil {
		return err
	}

	for key, value := range nodeMap {
		n.Name = key

		nodeDetails := value.(map[interface{}]interface{})
		n.IPAddr = nodeDetails["ip_addr"].(string)
		n.NodeType = nodeDetails["node_type"].(string)

		ports := nodeDetails["ports"].([]interface{})
		for _, port := range ports {
			portName := port.(string)
			n.Ports = append(n.Ports, portName)
		}
	}

	return nil
}

func (n Node) MarshalYAML() (interface{}, error) {
	nodeDetails := map[string]interface{}{
		"ip_addr":   n.IPAddr,
		"node_type": n.NodeType,
		"ports":     n.Ports,
	}

	return map[string]interface{}{
		n.Name: nodeDetails,
	}, nil
}

func actTopology() {

	/*
		TODO: We are reading info from the structured configuration files.
		TODO: So we dont need to run the inventory command anymore.
	*/
	inventoryData, err := getInventoryData(debug, "../inventory.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	hostnames := inventoryData.GetHost()
	if debug {
		fmt.Println("Hostnames:", hostnames)
	}
	// Load the YAML data from a file
	yamlData, err := os.ReadFile("topology.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config ACTTopologyConfig
	err = yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Print the unmarshalled data
	if debug {
		fmt.Printf("%+v\n", config)
	}
	// Add a new node to the config with an unique IP address
	config.AddIPToHosts(hostnames, "192.168.0.20")
	files, err := getYmlFiles("intended/structured_configs")
	if err != nil {
		fmt.Println(err)
	}
	var configs []*mo.Config
	for _, file := range files {
		fmt.Println("Working on file:", file)
		c, err := mo.NewConfigFromYaml(file)
		if err != nil {
			cobra.CheckErr(err)
		}
		configs = append(configs, c)
	}
	// Get all the ethernet interfaces from the structured configuration
	interfaceMap := make(map[string][]string)
	for _, c := range configs {
		var interfaces []string
		if _, ok := interfaceMap[c.Hostname]; !ok {
			interfaceMap[c.Hostname] = interfaces
		}
		for _, e := range c.EthernetInterfaces {
			interfaces = append(interfaces, e.Name)
		}
		interfaceMap[c.Hostname] = interfaces
	}

	// networkInterfaces interfaces connected to other network devices in this fabric
	networkInterfaces := make(map[string][]*mo.EthernetInterface)
	for _, c := range configs {
		networkInterfaces[c.Hostname] = NetworkInterfaces(c, hostnames)
	}
	linksMap := CreateLinkMap(networkInterfaces)
	fmt.Println("Links Map")
	for h, links := range linksMap {
		fmt.Println("Hosts", h)
		for _, l := range links {
			fmt.Println(l)
		}
	}
	config.AddPortsToNodes(interfaceMap)
	config.AddLinksToNodes(linksMap)

	// Output to a file
	yamlData, err = yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = os.WriteFile("topology-out.yml", yamlData, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// for h, host := range networkInterfaces {
	// 	fmt.Println("Hosts", h)
	// 	for _, e := range host {
	// 		fmt.Println(e.Peer)
	// 	}
	// }

	// fmt.Println("Link Map")
	// linksMap := CreateLinkMap(networkInterfaces)
	// for h, links := range linksMap {
	// 	fmt.Println("Hosts", h)
	// 	for _, l := range links {
	// 		fmt.Println(l)
	// 	}
	// }

}

// getInventoryData runs the ansible-inventory command and returns the InventoryData struct
// Ansible is needed to run this command
func getInventoryData(debug bool, inventoryPath string) (InventoryData, error) {
	ansibleInventoryOptions := inventory.AnsibleInventoryOptions{
		// Graph:     true,
		List:      true,
		Inventory: inventoryPath,
		Vars:      true,
		// Yaml:      true,
	}

	buff := new(bytes.Buffer)

	// Run the ansible-inventory command with options
	inventoryCmd := inventory.NewAnsibleInventoryCmd(
		inventory.WithPattern("all"),
		inventory.WithInventoryOptions(&ansibleInventoryOptions),
	)

	if debug {
		fmt.Println("Test strings:", inventoryCmd.String())
	}

	// Execute the command and write the output to the buffer
	exec := execute.NewDefaultExecute(
		execute.WithCmd(inventoryCmd),
		execute.WithWrite(io.Writer(buff)),
	)

	err := exec.Execute(context.TODO())
	if err != nil {
		return InventoryData{}, err
	}

	// Parse the JSON data
	var inventoryData InventoryData
	err = json.Unmarshal(buff.Bytes(), &inventoryData)
	if err != nil {
		return InventoryData{}, err
	}

	if debug {
		prettyJSON, err := json.MarshalIndent(inventoryData, "", "  ")
		if err != nil {
			return InventoryData{}, err
		}

		fmt.Println("Parsed JSON data:")
		fmt.Println(string(prettyJSON))
	}

	return inventoryData, nil
}

// AddIPToHosts adds a new node to the ACTTopologyConfig
func (c *ACTTopologyConfig) AddIPToHosts(hostnames []string, firstIP string) {
	ipInit := net.ParseIP(firstIP)
	ipNext := iplib.NextIP(ipInit)
	for _, hostname := range hostnames {
		if err := c.AddNode(hostname, ipNext.String()); err != nil {
			fmt.Println(err)
		}
		ipNext = iplib.NextIP(ipNext)
	}
}

func (c *ACTTopologyConfig) AddPortsToNodes(interfaceMap map[string][]string) {
	for _, node := range c.Nodes {
		// avoid overwriting the ports if the node is not in the interfaceMap
		if _, ok := interfaceMap[node.Name]; !ok {
			continue
		}
		node.Ports = interfaceMap[node.Name]
	}
}

/*
	 AddLinksToNodes adds links to the nodes in the ACTTopologyConfig
		linksMap should have this format:
		{"node1": ["node2:port1", "node3:port2"], "node2": ["node1:port1"]}
*/
func (c *ACTTopologyConfig) AddLinksToNodes(linksMap map[string][]string) {
	for hostname := range linksMap {
		c.Links = append(c.Links, &Link{Connection: linksMap[hostname]})
	}
}

// getYmlFiles returns a slice of all the .yml files in the given path
func getYmlFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".yml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// NonNetworkInterfaces returns a list list of interfaces that are not connected to a network device
func NonNetworkInterfaces(m *mo.Config, hostnames []string) []*mo.EthernetInterface {
	var nonNetworkInterfaces []*mo.EthernetInterface
	for _, e := range m.EthernetInterfaces {
		if !contains(hostnames, e.Peer) {
			nonNetworkInterfaces = append(nonNetworkInterfaces, &e)
		}
	}
	return nonNetworkInterfaces
}

func NetworkInterfaces(m *mo.Config, hostnames []string) []*mo.EthernetInterface {
	var networkInterfaces []*mo.EthernetInterface
	for _, e := range m.EthernetInterfaces {
		if contains(hostnames, e.Peer) {
			networkInterfaces = append(networkInterfaces, &e)
		}
	}
	return networkInterfaces
}

/*
	 CreateLinkMap creates a map of links between nodes
	 	ethernet should have this format:
		{"node1": [{"name": "Ethernet1", "peer": "node2"}, {"name": "Ethernet2", "peer": "node3"}], "node2": [{"name": "Ethernet1", "peer": "node1"}]}
		linksMap will have this format:
		[["node1:port1", "node1:port2"], ["node2:port1", "node3:port2"]]}
		notice that the key value is also infront of the port
*/
func CreateLinkMap(ethernets map[string][]*mo.EthernetInterface) [][]string {
	linksMap := make(map[string][]string)
	for _, interfaces := range ethernets {
		for _, e := range interfaces {
			// Add the port to the node
			if _, ok := linksMap[e.Peer]; !ok {
				linksMap[e.Peer] = []string{}
			}
			linksMap[e.Peer] = append(linksMap[e.Peer], fmt.Sprintf("%s:%s", e.Peer, e.Name))
		}
	}
	return linksMap
}

// contains checks if a string is in a slice of strings
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
