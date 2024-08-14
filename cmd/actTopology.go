/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/rogerscuall/crispy-enigma/internal/act"
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

var (
	// defaultPorts = []string{"Ethernet1-32"}
	debug = false
)

func actTopology() {

	files, err := getYmlFiles("intended/structured_configs")
	if err != nil {
		fmt.Println(err)
	}
	var network mo.Network
	var networkConfigs []*mo.Config
	for _, file := range files {
		fmt.Println("Working on file:", file)
		c, err := mo.NewConfigFromYaml(file)
		if err != nil {
			cobra.CheckErr(err)
		}
		networkConfigs = append(networkConfigs, c)
	}
	network.Configs = networkConfigs
	hostnames := network.GetHostnames()
	if debug {
		fmt.Println("Hostnames:", hostnames)
	}
	networkConnections := network.CleanNetworkConnections()
	interfaceMap := make(map[string][]string)
	for _, c := range networkConfigs {
		var interfaces []string
		if _, ok := interfaceMap[c.Hostname]; !ok {
			interfaceMap[c.Hostname] = interfaces
		}
		for _, e := range c.EthernetInterfaces {
			interfaces = append(interfaces, e.Name)
		}
		interfaceMap[c.Hostname] = interfaces
	}
	// Load the YAML data from a file
	yamlData, err := os.ReadFile("topology.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var actConfig act.TopologyConfig
	err = yaml.Unmarshal([]byte(yamlData), &actConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Print the unmarshalled data
	if debug {
		fmt.Printf("%+v\n", actConfig)
	}
	// Add a new node to the config with an unique IP address
	actConfig.AddIPToHosts(hostnames, "192.168.0.20")

	// Get all the ethernet interfaces from the structured configuration
	

	// networkInterfaces interfaces connected to other network devices in this fabric

	actConfig.AddPortsToNodes(interfaceMap)
	actConfig.AddLinksToNodes(networkConnections)

	// Output to a file
	yamlData, err = yaml.Marshal(&actConfig)
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
// func getInventoryData(debug bool, inventoryPath string) (InventoryData, error) {
// 	ansibleInventoryOptions := inventory.AnsibleInventoryOptions{
// 		// Graph:     true,
// 		List:      true,
// 		Inventory: inventoryPath,
// 		Vars:      true,
// 		// Yaml:      true,
// 	}

// 	buff := new(bytes.Buffer)

// 	// Run the ansible-inventory command with options
// 	inventoryCmd := inventory.NewAnsibleInventoryCmd(
// 		inventory.WithPattern("all"),
// 		inventory.WithInventoryOptions(&ansibleInventoryOptions),
// 	)

// 	if debug {
// 		fmt.Println("Test strings:", inventoryCmd.String())
// 	}

// 	// Execute the command and write the output to the buffer
// 	exec := execute.NewDefaultExecute(
// 		execute.WithCmd(inventoryCmd),
// 		execute.WithWrite(io.Writer(buff)),
// 	)

// 	err := exec.Execute(context.TODO())
// 	if err != nil {
// 		return InventoryData{}, err
// 	}

// 	// Parse the JSON data
// 	var inventoryData InventoryData
// 	err = json.Unmarshal(buff.Bytes(), &inventoryData)
// 	if err != nil {
// 		return InventoryData{}, err
// 	}

// 	if debug {
// 		prettyJSON, err := json.MarshalIndent(inventoryData, "", "  ")
// 		if err != nil {
// 			return InventoryData{}, err
// 		}

// 		fmt.Println("Parsed JSON data:")
// 		fmt.Println(string(prettyJSON))
// 	}

// 	return inventoryData, nil
// }

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
// func CreateLinkMap(ethernets map[string][]*mo.EthernetInterface) [][]string {
// 	linksMap := make(map[string][]string)
// 	for _, interfaces := range ethernets {
// 		for _, e := range interfaces {
// 			// Add the port to the node
// 			if _, ok := linksMap[e.Peer]; !ok {
// 				linksMap[e.Peer] = []string{}
// 			}
// 			linksMap[e.Peer] = append(linksMap[e.Peer], fmt.Sprintf("%s:%s", e.Peer, e.Name))
// 		}
// 	}
// 	return linksMap
// }

// contains checks if a string is in a slice of strings
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
