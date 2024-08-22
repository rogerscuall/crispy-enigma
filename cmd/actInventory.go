/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/rogerscuall/crispy-enigma/internal/ansible"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// actInventoryCmd represents the actInventory command
var actInventoryCmd = &cobra.Command{
	Use:   "actInventory",
	Short: "Takes an Ansible AVD inventory and updates it with an ACT inventory",
	Long: `ACT has a specific management interface, we need our ACT devices to use that one.
Most of the time it does not match the management interface in the Ansible AVD inventory.
This script will update the Ansible AVD inventory with the ACT inventory, and remove other fields like:
- serial_number
The act-inventory.yml should have a VEOS group with all the hosts listed after them like this example:
...
VEOS:
      hosts:
        ATL-ADM-BL201:
          ansible_host: 10.255.83.101
          ansible_user: cvpadmin
          ansible_ssh_pass: cvp123!
        ATL-ADM-BL202:
          ansible_host: 10.255.3.117
          ansible_user: cvpadmin
          ansible_ssh_pass: cvp123!
        ATL-ADM-LF203:
          ansible_host: 10.255.29.188
          ansible_user: cvpadmin
          ansible_ssh_pass: cvp123!
...`,
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
	actInventoryCmd.Flags().StringP("act", "a", "act-inventory.yml", "ACT inventory file")
	actInventoryCmd.Flags().StringP("original", "o", "inventory.yml", "Original inventory file")
	actInventoryCmd.Flags().StringP("output", "O", "updated-inventory.yml", "Output file")
}

/*
inventoryUpdate updates the original inventory with the ACT inventory
It will update the IP address of all the hosts in the original inventory with the IP address from the ACT inventory
It will remove the serial_number field from all hosts
The act-inventory.yml should have a VEOS group with all the hosts listed after them like this example:
```yaml
...
VEOS:

	hosts:
	  ATL-ADM-BL201:
	    ansible_host: 10.255.83.101
	    ansible_user: cvpadmin
	    ansible_ssh_pass: cvp123!
	  ATL-ADM-BL202:
	    ansible_host: 10.255.3.117
	    ansible_user: cvpadmin
	    ansible_ssh_pass: cvp123!
	  ATL-ADM-LF203:
	    ansible_host: 10.255.29.188
	    ansible_user: cvpadmin
	    ansible_ssh_pass: cvp123!

...
```
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

	var original, act, newInventory ansible.Inventory

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

	// Create a map of hostnames to IP addresses from act inventory
	actHosts := make(map[string]string)
	if veos, ok := act.All.Children["VEOS"].(map[string]interface{}); ok {
		if hosts, ok := veos["hosts"].(map[string]interface{}); ok {
			for hostname, data := range hosts {
				if hostData, ok := data.(map[string]interface{}); ok {
					if ip, ok := hostData["ansible_host"].(string); ok {
						actHosts[hostname] = ip
					}
				}
			}
		}
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
