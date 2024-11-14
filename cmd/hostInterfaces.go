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
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/rogerscuall/crispy-enigma/hostfiles"
	"github.com/spf13/cobra"
)

// hostInterfacesCmd represents the hostInterfaces command
var hostInterfacesCmd = &cobra.Command{
	Use:   "hostInterfaces",
	Short: "Reads a CSV file and creates a YAML file with the host interfaces",
	Long: `It will read a CSV file with the following format:
interface,description,shutdown
Ethernet1,description,up
And will create a YAML file with the following format:
---
csc_ethernet_interfaces:
	- name: "Ethernet1"
	  description: "description"
	  state: "up"
Both files will have the same name, but different extensions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hostInterfaces called")
		folder := cmd.Flag("folder").Value.String()
		interfaceRange := cmd.Flag("range").Value.String()
		fmt.Println("Folder:", folder)
		fmt.Println("Range:", interfaceRange)
		base, lower, higher, err := hostfiles.ParseInterfaceString(interfaceRange)
		if err != nil {
			cobra.CheckErr(err)
		}
		fmt.Printf("Base: %v, Lower: %v, Higher: %v\n", base, lower, higher)
		defaultInterfaces := hostfiles.CreateDefaultInterfaces(base, lower, higher)
		fmt.Printf("Default interfaces: %v\n", defaultInterfaces)
		files, err := getCsvFiles(folder)
		if err != nil {
			fmt.Println(err)
		}

		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			// Read the file
			reader := csv.NewReader(f)
			records, err := reader.ReadAll()
			if err != nil {
				cobra.CheckErr(err)
			}

			var interfaces []hostfiles.Interface

			// Iterate through the records

			for i, record := range records {
				if i == 0 { // Skip header
					continue
				}
				shutdown := true
				if record[2] == "false" {
					shutdown = false
				}
				interfaces = append(interfaces, hostfiles.Interface{
					Name:        record[0],
					Description: record[1],
					Shutdown:    shutdown,
				})
			}
			log.Println("Interfaces:", interfaces)
			log.Println("File:", file)
			nonDefaultMap := make(map[string]hostfiles.Interface)
			for _, item := range interfaces {
				nonDefaultMap[item.Name] = item
			}
			for _, item := range defaultInterfaces {
				if _, exists := nonDefaultMap[item.Name]; !exists {
					nonDefaultMap[item.Name] = item
				}
			}
			var finalInterfaces []hostfiles.Interface
			for _, item := range nonDefaultMap {
				finalInterfaces = append(finalInterfaces, item)
			}
			log.Println("Interfaces:", interfaces)
			log.Println("File:", file)
			log.Println("Default interfaces:", defaultInterfaces)
			log.Println("Final interfaces:", finalInterfaces)
			hostfiles.WriteYamlFile(file, finalInterfaces)

		}
	},
}

func init() {
	rootCmd.AddCommand(hostInterfacesCmd)
	hostInterfacesCmd.Flags().StringP("folder", "f", "", "Path to the folder")
	hostInterfacesCmd.Flags().StringP("range", "r", "Ethernet10-20", "Interfaces range")
	err := hostInterfacesCmd.MarkFlagRequired("folder")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

}
