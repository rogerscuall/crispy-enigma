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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rogerscuall/crispy-enigma/hostfiles"
	"github.com/rogerscuall/crispy-enigma/internal/pocketbase"
	pbmodels "github.com/rogerscuall/crispy-enigma/internal/pocketbase/models"
	mo "github.com/rogerscuall/crispy-enigma/models"
	"github.com/spf13/cobra"
)

// pbAccessInterfaceCmd represents the pbAccessInterface command
var pbAccessInterfaceCmd = &cobra.Command{
	Use:   "pbAccessInterface",
	Short: "Gather access port information from PocketBase",
	Long:  `This command connects to PocketBase and retrieves information about access ports on a switch.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pbAccessInterface called")
		folder := cmd.Flag("folder").Value.String()
		log.Print("Folder:", folder)
		if verbose {
			log.Print("Debug mode enabled")
			app.Debug = true
		}
		if pbURL == "" {
			fmt.Println("PB_URL is not set")
			os.Exit(1)
		}
		if pbUsername == "" {
			fmt.Println("PB_USERNAME is not set")
			os.Exit(1)
		}
		if pbPassword == "" {
			fmt.Println("PB_PASSWORD is not set")
			os.Exit(1)
		}
		app.PbClient = pocketbase.NewClient(pbURL, pbUsername, pbPassword)

		files, err := getYmlFiles(folder)
		if err != nil {
			fmt.Println(err)
		}
		// create a host_vars folder if it doesn't exist
		err = os.MkdirAll("host_vars", os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		for _, file := range files {
			fmt.Println("Working on file:", file)
			var hostfile []hostfiles.Interface
			config, err := mo.NewConfigFromYaml(file)
			if err != nil {
				cobra.CheckErr(err)
			}
			// Get all the interfaces from a switch
			interfaces, err := app.PbClient.GetInterfacesFromSwitchName(config.Hostname)
			if err != nil {
				app.DebugLog("Device was not found %s", config.Hostname)
				app.DebugLog("Error: %s", err)
				continue
			}
			for _, intf := range interfaces {
				//Get VLAN Name
				var vlan pbmodels.VLAN
				resp, err := app.PbClient.View("vlan", intf.VLAN)
				if err != nil {
					app.DebugLog("VLAN was not found %s", intf.VLAN)
				} else {

					err = json.Unmarshal(resp, &vlan)
					if err != nil {
						app.DebugLog("Error unmarshalling VLAN %s", err)
					}
					fmt.Println("VLAN:", vlan.Name)
				}
				vlanNumber, err := strconv.Atoi(vlan.Name)
				if err != nil {
					app.DebugLog("Error converting VLAN to int %s", err)
				}
				hostfile = append(hostfile, hostfiles.Interface{
					Name:        intf.Name,
					Description: intf.Description,
					Shutdown:    intf.State,
					VLAN:        vlanNumber,
				})
			}
			hostfiles.WriteHostFile(fmt.Sprintf("host_vars/%s.yml", config.Hostname), hostfile)
		}
	},
}

func init() {
	rootCmd.AddCommand(pbAccessInterfaceCmd)
	pbURL = os.Getenv("PB_URL")
	pbUsername = os.Getenv("PB_USERNAME")
	pbPassword = os.Getenv("PB_PASSWORD")
	pbAccessInterfaceCmd.Flags().StringP("folder", "f", "", "Path to the folder")
	err := pbAccessInterfaceCmd.MarkFlagRequired("folder")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	//cvpConfigCmd.Flags().BoolP("debug", "v", false, "Debug")

}
