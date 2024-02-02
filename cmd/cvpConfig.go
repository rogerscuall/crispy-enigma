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

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/rogerscuall/crispy-enigma/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/aristanetworks/go-cvprac.v2/client"
)

var cvaasURL = "https://www.arista.io/cvpservice"

// cvpConfigCmd represents the cvpConfig command
var cvpConfigCmd = &cobra.Command{
	Use:   "cvpConfig",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cvpConfig called")
		folder := cmd.Flag("folder").Value.String()
		log.Print("Folder:", folder)
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			log.Print("Debug mode enabled")
			app.Debug = true
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
			cvpClient.Hosts = []string{cvpURL}
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
		files, err := getConfigFiles(folder)
		if err != nil {
			fmt.Println(err)
		}
		for _, file := range files {
			app.DebugLog("File Name: %v\n", file)
			deviceName := strings.TrimSuffix(path.Base(file), ".cfg")
			dev, err := cvpClient.API.GetDeviceByName(deviceName)
			if err != nil {
				log.Printf("Device %v not found in CVP", deviceName)
				log.Printf("ERROR: %s", err)
				continue
			}
			if dev == nil {
				log.Printf("Device %v not found in CVP", deviceName)
				continue
			}
			app.DebugLog("Device Hostname: %v\n", dev.Hostname)
			app.DebugLog("Device Serial Number: %v\n", dev.SerialNumber)
			app.DebugLog("Device System MAC Address: %v\n", dev.SystemMacAddress)
			config, err := cvpClient.API.GetConfigletsByDeviceID(dev.SystemMacAddress)
			if err != nil {
				log.Printf("Configlets for device %v not found in CVP", deviceName)
				log.Printf("ERROR: %s", err)
			}

			f, err := os.Open(file)
			if err != nil {
				log.Printf("Error: %v", err)
				continue
			}
			defer f.Close()
			newConfig, err := io.ReadAll(f)
			if err != nil {
				log.Printf("Error reading file: %v", err)
				continue
			}
			app.DebugLog("Number of Configlets: %v\n", len(config))
			for _, configlet := range config {
				app.DebugLog("Configlet Name: %v\n", configlet.Name)
				if err != nil {
					log.Printf("Error reading file: %v\n", err)
					continue
				}
				edits := myers.ComputeEdits(span.URIFromPath(file), configlet.Config, string(newConfig))
				diff := fmt.Sprint(gotextdiff.ToUnified("running-config", "intended-config", configlet.Config, edits))
				if diff != "" {
					fmt.Println("Device Config Diff:", deviceName)
					fmt.Print(diff)
				}
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(cvpConfigCmd)
	cvpURL = os.Getenv("CVP_URL")
	cvpUsername = os.Getenv("CVP_USERNAME")
	cvpPassword = os.Getenv("CVP_PASSWORD")
	cvpToken = os.Getenv("CVP_TOKEN")
	app = pkg.NewApplication()
	cvpConfigCmd.Flags().StringP("folder", "f", "", "Folder where the structured config YAML files are located")
	err := cvpConfigCmd.MarkFlagRequired("folder")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	cvpConfigCmd.Flags().BoolP("debug", "v", false, "Debug")
}
