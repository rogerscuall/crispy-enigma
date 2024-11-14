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

	"github.com/rogerscuall/crispy-enigma/internal/act"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var hostKeys act.HostKeyList

// actInventoryCmd represents the actInventory command
var actInventoryCmd = &cobra.Command{
	Use:   "actInventory",
	Short: "Takes an Ansible AVD inventory and updates it with an ACT topology",
	Long: `
This script will update the Ansible AVD inventory with the ACT topology, and remove other fields like:
- serial_number
- is_deployed
It also can take a host_key defined in the .crispy-enigma.yaml file and update the inventory with the new value.
This file should be in the same directory from where the script is executed.
Example of .crispy-enigma.yaml:
host_keys:
  - host: cvp
    key: ansible_host
    newvalue: "10-255-31-17.some.act.arista.com"
  - host: cvp
    key: ansible_user
    newvalue: admin
  ...
To create an ACT topology use the command actTopology`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("actInventory called")
		inventory := cmd.Flag("inventory").Value.String()
		outputFile := cmd.Flag("update").Value.String()
		if err := viper.UnmarshalKey("host_keys", &hostKeys); err != nil {
			fmt.Println("Error unmarshaling provider: " + err.Error())
		}

		data, err := os.ReadFile(inventory)
		if err != nil {
			cobra.CheckErr(err)
		}

		var inventoryData interface{}
		err = yaml.Unmarshal(data, &inventoryData)
		if err != nil {
			cobra.CheckErr(err)
		}
		keysToRemove := []string{"serial_number", "is_deployed"}

		act.RemoveKeys(inventoryData, keysToRemove)
		act.RemoveNulls(inventoryData)
		for _, hostKey := range hostKeys {
			act.UpdateHostKey(inventoryData, hostKey.Host, hostKey.Key, hostKey.NewValue)
		}
		updatedData, err := yaml.Marshal(&inventoryData)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		err = os.WriteFile(outputFile, updatedData, 0644)
		if err != nil {
			cobra.CheckErr(err)
		}
		fmt.Println("Inventory updated and saved to", outputFile)

	},
}

func init() {
	rootCmd.AddCommand(actInventoryCmd)
	actInventoryCmd.Flags().StringP("inventory", "i", "inventory.yml", "Inventory file")
	actInventoryCmd.Flags().StringP("update", "u", "act-inventory.yml", "Output file")
}
