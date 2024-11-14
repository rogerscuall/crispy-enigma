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
	"path"
	"strings"

	"github.com/rogerscuall/crispy-enigma/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/aristanetworks/go-cvprac.v2/client"
)

// cvpPendingTaskCmd represents the cvpPendingTask command
var cvpPendingTaskCmd = &cobra.Command{
	Use:   "cvpPendingTask",
	Short: "Check to see if any device on CVP has a pending task",
	Long: `For every device in the current folder, check if there is a pending task in CVP.
We will print a information of the device only if the WorkOrderUserDefinedStatus is "Pending"
and the WorkOrderState is "ACTIVE"
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cvpPendingTask called")
		folder := cmd.Flag("folder").Value.String()
		log.Print("Folder:", folder)
		debug, _ := cmd.Flags().GetBool("verbose")
		if debug {
			log.Print("Debug mode enabled")
			app.Debug = true
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
		tasks, err := cvpClient.API.GetAllTasks()
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
		for _, task := range tasks {
			app.DebugLog("Task ID: %v\n", task.WorkOrderID)
			app.DebugLog("Task Device Name: %v\n", task.WorkOrderDetails.NetElementHostName)
			app.DebugLog("Task Status: %v\n", task.WorkOrderState)
			app.DebugLog("Task  User Status: %v\n", task.WorkOrderUserDefinedStatus)
			app.DebugLog("Task Description: %v\n", task.Description)
		}
		var finalOutputClean = true
		for _, file := range files {
			app.DebugLog("File Name: %v\n", file)
			deviceName := strings.TrimSuffix(path.Base(file), ".cfg")
			if fqdnSuffix != "" {
				deviceName = deviceName + "." + fqdnSuffix
			}
			clean := true
			for _, task := range tasks {
				if task.WorkOrderDetails.NetElementHostName == deviceName {
					if task.WorkOrderState == "ACTIVE" && task.WorkOrderUserDefinedStatus == "Pending" {
						clean = false
						log.Printf("The device %v has pending task %v", deviceName, task.WorkOrderID)
					}
				}
			}
			if clean {
				log.Printf("The device %v is clean", deviceName)
			} else {
				finalOutputClean = false
			}
			clean = true
		}
		if finalOutputClean {
			log.Printf("No pending tasks for this network")
		} else {
			log.Printf("There are pending tasks for this network")
			
		}
	},
}

func init() {
	rootCmd.AddCommand(cvpPendingTaskCmd)
	cvpURL = os.Getenv("CVP_URL")
	cvpUsername = os.Getenv("CVP_USERNAME")
	cvpPassword = os.Getenv("CVP_PASSWORD")
	cvpToken = os.Getenv("CVP_TOKEN")
	fqdnSuffix = os.Getenv("FQDN_SUFFIX")
	app = pkg.NewApplication()
	cvpPendingTaskCmd.Flags().StringP("folder", "f", "", "Folder to store running-config files")
	err := cvpPendingTaskCmd.MarkFlagRequired("folder")
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

}
