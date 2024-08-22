/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/rogerscuall/crispy-enigma/internal/act"
	"github.com/rogerscuall/crispy-enigma/internal/ansible"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// actInventoryCmd represents the actInventory command
var actInventoryCmd = &cobra.Command{
	Use:   "actInventory",
	Short: "Takes an Ansible AVD inventory and updates it with an ACT topology",
	Long: `ACT has a specific management interface, we need our ACT devices to use that one.
Most of the time it does not match the management interface in the Ansible AVD inventory.
This script will update the Ansible AVD inventory with the ACT topology, and remove other fields like:
- serial_number
To create an ACT topology use the command actTopology`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("actInventory called")
		original := cmd.Flag("original").Value.String()
		act := cmd.Flag("act").Value.String()
		outputFile := cmd.Flag("output").Value.String()
		inventoryUpdate(original, act, outputFile)
	},
}

func init() {
	rootCmd.AddCommand(actInventoryCmd)
	actInventoryCmd.Flags().StringP("act", "a", "act-topology.yml", "ACT topology file")
	actInventoryCmd.Flags().StringP("original", "o", "inventory.yml", "Original inventory file")
	actInventoryCmd.Flags().StringP("output", "O", "updated-inventory.yml", "Output file")
}

/*
inventoryUpdate updates the original inventory with the ACT inventory
It will update the IP address of all the hosts in the original inventory with the IP address from the ACT inventory
It will remove the serial_number field from all hosts
To create an ACT topology use the command actTopology
*/
func inventoryUpdate(originalInventory, actInventory, outputFile string) {
	// Read the original inventory file
	originalData, err := os.ReadFile(originalInventory)
	if err != nil {
		log.Fatalf("Error reading original inventory: %v", err)
	}

	// Read the act inventory file
	actData, err := os.ReadFile(actInventory)
	if err != nil {
		log.Fatalf("Error reading act inventory: %v", err)
	}

	var original, newInventory ansible.Inventory
	var act act.TopologyConfig

	// Unmarshal both YAML files
	err = yaml.Unmarshal(originalData, &original)
	if err != nil {
		log.Fatalf("Error unmarshaling original inventory: %v", err)
	}

	err = yaml.Unmarshal(actData, &act)
	if err != nil {
		log.Fatalf("Error unmarshaling act inventory: %v", err)
	}

	// Deep copy the original inventory to the new inventory
	newInventoryData, err := yaml.Marshal(&original)
	if err != nil {
		log.Fatalf("Error marshaling original inventory: %v", err)
	}
	err = yaml.Unmarshal(newInventoryData, &newInventory)
	if err != nil {
		log.Fatalf("Error unmarshaling to new inventory: %v", err)
	}

	actHosts := make(map[string]string)
	for _, node := range act.Nodes {
		actHosts[node.Name] = node.IPAddr
	}

	// Update the new inventory
	updateInventory(newInventory.All.Children, actHosts)
	var buf bytes.Buffer
	// create a writer with a buffer

	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	// Marshal the updated inventory back to YAML
	err = encoder.Encode(newInventory)
	if err != nil {
		log.Fatalf("Error marshaling updated inventory: %v", err)
	}

	// Write the updated inventory to a new file
	err = os.WriteFile(outputFile, buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Error writing updated inventory: %v", err)
	}

	fmt.Printf("Updated inventory written to %v.", outputFile)
}

func updateInventory(children map[string]interface{}, actHosts map[string]string) {
	for _, child := range children {
		switch v := child.(type) {
		case map[string]interface{}:
			if hosts, ok := v["hosts"].(map[string]interface{}); ok {
				for hostname, hostData := range hosts {
					if data, ok := hostData.(map[string]interface{}); ok {
						// Update ansible_host if it exists in actHosts
						if newIP, exists := actHosts[hostname]; exists {
							data["ansible_host"] = newIP
						}
						// Remove serial_number field
						delete(data, "serial_number")
					}
				}
			}
			// Recursively update nested children
			updateInventory(v, actHosts)
		}
	}
}
