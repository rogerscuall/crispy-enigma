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
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/rogerscuall/crispy-enigma/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/aristanetworks/go-cvprac.v2/client"
)

/*LISTEN: This is a possibility but hard to implement.
This command should be similar to cvpConfig but instead of checking for AVD configlets,
 it checks again the last known running-config in CVP with: cvpClient.API.GetInventoryConfiguration(dev.SystemMacAddress).
The idea is to compare the AVD config with the last known running config from device that exist in CVP.
the are 2 **problems**:
 * the AVD config and the device config are sometimes the same but in different order.
 * the device config does not show the defaults and the AVD config does.
	* device config: server x.x.x.x
	* AVD config: server x.x.x.x vrf default
TODO: Order both configurations before comparing them.
Once order is applied the position is lost and we will only provide a yes or no answer.
*/
// compareDeviceCmd represents the compareDevice command
var compareDeviceCmd = &cobra.Command{
	Use:   "compareDevice",
	Short: "Download running-config from CVP and compares with intended config",
	Long: `For every device in the folder will check for the last known running-config in CVP
and compare with the AVD config. The last known running-config and the AVD config are out of order even if they are the same.
Due to this, the output will only provide a yes or no answer. 
Useful to check at the pipeline level there is a config conflict in CVP. 
If token is set, it will to check at the pipeline level if a the build will update CVP. If token is set, it will
take precedence over username and password. If CVP_URL is not set, it will use CVAAS.
The following variables can be set as environment variables or in a .env file.
- CVP_USERNAME: CVP Username
- CVP_PASSWORD: CVP Password
- CVP_TOKEN: CVP Token takes precedence over username and password
- CVP_URL: CVP URL if not set it will use CVAAS
The AVD configlet is named <FABRIC_NAME>_<DEVICE_NAME>.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("compareDevice called")
		folder := cmd.Flag("folder").Value.String()
		log.Print("Folder:", folder)
		debug, _ := cmd.Flags().GetBool("verbose")
		if debug {
			log.Print("Debug mode enabled")
			app.Debug = true
		}
		checkConfiglets, _ := cmd.Flags().GetBool("check-all-configlets")
		if checkConfiglets {
			log.Print("Checking all configlets")
		}
		files, err := getConfigFiles(folder)
		if err != nil {
			log.Fatalf("Error reading folder: %v", err)
		}
		cvpClient, _ := client.NewCvpClient(
			client.Protocol("https"),
			client.Port(443),
			client.Debug(debug))
		// use CVAAS if CVP_URL is not set
		if cvpURL == "" {
			app.DebugLog("CVP_URL not set, using %v", cvaasURL)
			cvpClient.Client.HostURL = cvaasURL
		} else {
			app.DebugLog("Using CVP_URL: %v", cvpURL)
			hosts := []string{cvpURL}
			cvpClient.SetHosts(hosts...)
		}
		// if token is set, use token authentication and ignore username and password
		if cvpToken != "" {
			app.DebugLog("Using Token authentication")
			cvpClient.Client.SetAuthToken(cvpToken)
		} else {
			app.DebugLog("Using Username and Password authentication")
			if err := cvpClient.Connect(cvpUsername, cvpPassword); err != nil {
				log.Fatalf("ERROR: %s", err)
			}
		}
		app.CVPClient = cvpClient
		// verify we have at least one device in inventory
		data, err := cvpClient.API.GetCvpInfo()
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
		app.DebugLog("Data: %v\n", data)
		// testing authentication by getting cvp info
		info, err := app.CVPClient.API.GetCvpInfo()
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
		app.DebugLog("CVP Info: %v\n", info)
		err = os.Mkdir("running-config", 0755)
		if err != nil {
			log.Printf("Error creating running-config directory: %v", err)
		}

		var inSync = true
		totalDiff := make(map[string]string)
		for _, file := range files {
			app.DebugLog("File Name: %v\n", file)
			deviceName := strings.TrimSuffix(path.Base(file), ".cfg")
			if fqdnSuffix != "" {
				deviceName = deviceName + "." + fqdnSuffix
			}
			log.Printf("Working on device: %v\n", deviceName)
			dev, err := cvpClient.API.GetDeviceByName(deviceName)
			if err != nil {
				log.Printf("Device %v not found in CVP", deviceName)
				log.Printf("ERROR: %s", err)
				inSync = false
				continue
			}
			if dev == nil {
				log.Printf("Device %v not found in CVP", deviceName)
				inSync = false
				continue
			}
			app.DebugLog("Device Hostname: %v\n", dev.Hostname)
			app.DebugLog("Device Serial Number: %v\n", dev.SerialNumber)
			app.DebugLog("Device System MAC Address: %v\n", dev.SystemMacAddress)

			deviceConfig, err := cvpClient.API.GetInventoryConfiguration(dev.SystemMacAddress)
			if err != nil {
				log.Printf("Error getting device configuration: %v", err)
				inSync = false
				continue
			}

			if deviceName == "dc1-leaf2bx" {
				fmt.Println("Device Configuration conf:", deviceConfig.Output)
			}

			// config := deviceConfig.Output

			f, err := os.Open(file)
			if err != nil {
				log.Printf("Error: %v", err)
				inSync = false
				continue
			}
			defer f.Close()
			newConfig, err := io.ReadAll(f)
			if err != nil {
				log.Printf("Error reading file: %v", err)
				inSync = false
				continue
			}

			// edits := myers.ComputeEdits(span.URIFromPath(file), config, string(newConfig))
			// diff := fmt.Sprint(gotextdiff.ToUnified("running-config", "designed-config", config, edits))
			// if diff != "" {
			// 	totalDiff[deviceName] = diff
			// } else {
			// 	log.Printf("Device %v config is in sync\n", deviceName)
			// }
			// create a file with the running config
			fileName := fmt.Sprintf("running-config/%v.cfg", deviceName)
			err = os.WriteFile(fileName, newConfig, 0644)
			if err != nil {
				log.Printf("Error writing file: %v\n", err)
			}

		}
		if len(totalDiff) == 0 && inSync {
			fmt.Println("All devices are in sync")
		}
		for name, diff := range totalDiff {
			fmt.Println("Device Config Diff:", name)
			fmt.Print(diff)
		}
	},
}

func init() {
	configCmd.AddCommand(compareDeviceCmd)
	cvpURL = os.Getenv("CVP_URL")
	cvpUsername = os.Getenv("CVP_USERNAME")
	cvpPassword = os.Getenv("CVP_PASSWORD")
	cvpToken = os.Getenv("CVP_TOKEN")
	app = pkg.NewApplication()
	compareDeviceCmd.Flags().StringP("folder", "f", "", "Folder where the structured config YAML files are located")
	compareDeviceCmd.Flags().BoolP("check-all-configlets", "c", false, "Check all configlets instead of the AVD configlet")
	err := compareDeviceCmd.MarkFlagRequired("folder")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

}
