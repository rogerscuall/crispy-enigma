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

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	"github.com/rogerscuall/crispy-enigma/infoblox"
	"github.com/rogerscuall/crispy-enigma/models"
	"github.com/rogerscuall/crispy-enigma/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// infobloxNextNetworkCmd represents the infobloxNextNetwork command
var infobloxNextNetworkCmd = &cobra.Command{
	Use:   "infobloxNextNetwork",
	Short: "Gather the next available network from Infoblox",
	Long: `Reserve the next available network from Infoblox and update the NETWORK_SERVICES.yml 
file with the new network on AVD to update the network services`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("infobloxNextNetwork called")
		filePath := cmd.Flag("file").Value.String()
		log.Print("File: ", filePath)
		containerName := cmd.Flag("container").Value.String()
		log.Print("Container Name: ", containerName)
		maskPrefix, err := cmd.Flags().GetUint("mask")
		if err != nil {
			cobra.CheckErr(err)
		}
		log.Print("Mask: ", maskPrefix)
		numberNetwork, err := cmd.Flags().GetUint("number")
		if err != nil {
			cobra.CheckErr(err)
		}
		log.Print("Number of subnets: ", numberNetwork)
		file, err := os.Open(filePath)
		if err != nil {
			cobra.CheckErr(err)
		}
		defer file.Close()
		networkService := &models.NetworkService{}
		decoder := yaml.NewDecoder(file)
		err = decoder.Decode(networkService)
		if err != nil {
			cobra.CheckErr(err)
		}
		conn, err := infoblox.CreateConnector(infoURL, infoWapiVersion, infoUsername, infoPassword)
		if err != nil {
			cobra.CheckErr(err)
		}
		defer conn.Logout()
		objMgr := ibclient.NewObjectManager(conn, "myclient", "")

		for i := 0; i < int(numberNetwork); i++ {
			network, err := infoblox.GetNextAvailableNetwork(objMgr, "default", containerName, uint(maskPrefix))
			if err != nil {
				panic(err)
			}
			log.Println("Network: ", network)
			svi := models.NewSVI(network)
			networkService.Tenants[0].VRFs[0].SVIs = append(networkService.Tenants[0].VRFs[0].SVIs, *svi)
		}

		updatedYml, err := yaml.Marshal(networkService)
		if err != nil {
			cobra.CheckErr(err)
		}
		tmpFile, err := os.Create("new_network_services.yml")
		if err != nil {
			cobra.CheckErr(err)
		}
		defer tmpFile.Close()
		if _, err := tmpFile.Write(updatedYml); err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(infobloxNextNetworkCmd)

	// Initialize global variables
	infoURL = os.Getenv("INFOBLOX_URL")
	infoUsername = os.Getenv("INFOBLOX_USERNAME")
	infoPassword = os.Getenv("INFOBLOX_PASSWORD")
	infoWapiVersion = os.Getenv("INFOBLOX_WAPI_VERSION")

	app = pkg.NewApplication()
	var err error
	app.InfobloxClient, err = infoblox.CreateConnector(infoURL, infoWapiVersion, infoUsername, infoPassword)
	if err != nil {
		panic(err)
	}
	infobloxNextNetworkCmd.Flags().StringP("container", "c", "10.0.0.0/16", "Name of the container to create the network")
	infobloxNextNetworkCmd.Flags().StringP("file", "f", "", "Path to the file NETWORK_SERVICES.yml")
	err = infobloxNextNetworkCmd.MarkFlagRequired("file")
	if err != nil {
		cobra.CheckErr(err)
	}
	infobloxNextNetworkCmd.Flags().UintP("mask", "m", 24, "Mask for the subnet for the new network")
	infobloxNextNetworkCmd.Flags().UintP("number", "n", 1, "Number of subnets to add")
}
