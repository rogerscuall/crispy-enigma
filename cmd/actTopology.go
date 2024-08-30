/*
Copyright © 2024 Roger Gomez rogerscuall@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rogerscuall/crispy-enigma/internal/act"
	mo "github.com/rogerscuall/crispy-enigma/models"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

/*
- Parse the JSON data
- Extract the hostnames
- Load the ACT Topology YAML data from a file. This ACT Topology has some data that needs to be kept, we are just adding new nodes and links to it.
- Add the hostnames to the ACT Topology data.
- Parse the information from the structured configuration files.
- From the structured configuration data, parse the interface ethernet configuration as links struct.
- Identify what ports are available in the test environment, if the port does not connect to another switch is available to be used. (Done)
- The available ports should be added to the test node as links in the ACT Topology data.
- Add the links struct to the ACT Topology data for the "no available" ports. These are the ports that connect the switches.
- Output the ACT Topology data to a file
*/

// actTopologyCmd represents the actTopology command
var actTopologyCmd = &cobra.Command{
	Use:   "actTopology",
	Short: "From AVD project creates an ACT topology",
	Long: `This command will create an ACT topology from and AVD project
It needs an input ACT topology (default: topology.yml) 
and a folder with the structured configuration files (default: intended/structured_configs)
It will output the ACT topology to a file (default: act-topology.yml)
For an example input file run the command with the -e flag`,
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
		input := cmd.Flag("input").Value.String()
		actTopology(folder, input, output)
	},
}

func init() {
	rootCmd.AddCommand(actTopologyCmd)
	// Prints an example input file
	actTopologyCmd.Flags().StringP("folder", "f", "intended/structured_configs", "Folder with the structured configuration files")
	actTopologyCmd.Flags().StringP("input", "i", "topology.yml", "ACT Topology file")
	actTopologyCmd.Flags().BoolP("example", "e", false, "Prints an example input file")
	actTopologyCmd.Flags().StringP("output", "O", "act-topology.yml", "Output file")
}

var (
	// defaultPorts = []string{"Ethernet1-32"}
	debug = false
)

func actTopology(folder, inputActTopology, actTopology string) {
	files, err := getYmlFiles(folder)
	if err != nil {
		fmt.Println(err)
	}
	var network mo.Network
	var networkConfigs []*mo.Config
	for _, file := range files {
		// Ignore files inside the "cvp" folder with a file name that begins with "cvp"
		// This file belong to CVP and not to the network devices
		if strings.Contains(file, "cvp/cvp") {
			continue
		}
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
	yamlData, err := os.ReadFile(inputActTopology)
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
	err = actConfig.AddNodes(network)
	if err != nil {
		cobra.CheckErr(err)
	}
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
# Update the IP addresses of the test nodes and CVP to be in the same subnet as the devices.
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
