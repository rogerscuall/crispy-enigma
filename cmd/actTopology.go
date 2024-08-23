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
7. From the structured configuration data, parse the interface ethernet configuration as links struct.
8. Identify what ports are available in the test environment, if the port does not connect to another switch is available to be used. (Done)
8.5. The available ports should be added to the test node as links in the ACT Topology data.
9. Add the links struct to the ACT Topology data for the "no available" ports. These are the ports that connect the switches.
10. Output the ACT Topology data to a file
*/

// actTopologyCmd represents the actTopology command
var actTopologyCmd = &cobra.Command{
	Use:   "actTopology",
	Short: "From AVD project creates an ACT topology",
	Long:  `This command will create an ACT topology from and AVD project. It needs the intended configuration in structured format and add test nodes to the topology for data plane testing`,
	Run: func(cmd *cobra.Command, args []string) {
		example, _ := cmd.Flags().GetBool("example")
		if example {
			fmt.Println("Example:")
			fmt.Println(actExampleConfig)
			os.Exit(0)
		}
		fmt.Println("actTopology called")
		folder := cmd.Flag("folder").Value.String()
		output := cmd.Flag("output").Value.String()
		actTopology(folder, output)
	},
}

func init() {
	rootCmd.AddCommand(actTopologyCmd)
	// Prints an example input file
	actTopologyCmd.Flags().StringP("folder", "f", "intended/structured_configs", "Folder with the structured configuration files")
	actTopologyCmd.Flags().BoolP("example", "e", false, "Prints an example input file")
	actTopologyCmd.Flags().StringP("output", "O", "act-topology.yml", "Output file")
}

var (
	// defaultPorts = []string{"Ethernet1-32"}
	debug = false
)

func actTopology(folder, actTopology string) {
	files, err := getYmlFiles(folder)
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
	if verbose {
		fmt.Printf("%+v\n", actConfig)
	}
	// Add a new node to the config with an unique IP address
	// actConfig.AddIPToHosts(hostnames, "192.168.0.20")
	actConfig.AddNodes(network)

	// Get all the ethernet interfaces from the structured configuration files
	// networkInterfaces interfaces connected to other network devices in this fabric

	actConfig.AddPortsToNodes(network)
	actConfig.AddLinksToNodes(network)

	// Output to a file
	yamlData, err = yaml.Marshal(&actConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = os.WriteFile(actTopology, yamlData, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("ACT Topology data written to", actTopology)
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

var actExampleConfig = `
veos:
  password: cvp123!
  username: cvpadmin
  version: 4.27.0F

cvp:
  password: cvproot
  username: root
  version: 2022.2.2

generic:
  password: ansible
  username: ansible
  version: Rocky-8.5

third-party:
  password: ansible
  username: ansible
  version: infoblox

nodes:
  - CVP:
      ip_addr: 192.168.0.5
      node_type: cvp
      auto_configuration: true
  - INTERNAL-TEST:
      ip_addr: 192.168.0.11
      node_type: veos
      ports:
        - Ethernet1-32
  - EXTERNAL-TEST:
      ip_addr: 192.168.0.11
      node_type: veos
      ports:
        - Ethernet1-32
links:
  - connection:
      - INTERNAL-TEST:Ethernet1
      # Update this next line with the name of an existing switch and a port that is configured but not connected to another switch in this network
      - DC1-L2LEAF1A:Ethernet20
  - connection:
      - EXTERNAL-TEST:Ethernet1
      # Update this next line with the name of an existing switch and a port that is configured but not connected to another switch in this network
      - DC1-L2LEAF1A:Ethernet21
`
