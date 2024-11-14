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
	"net/url"
	"os"

	"github.com/go-resty/resty/v2"
	n "github.com/netbox-community/go-netbox/v3"
	mo "github.com/rogerscuall/crispy-enigma/models"
	"github.com/rogerscuall/crispy-enigma/netbox"
	"github.com/rogerscuall/crispy-enigma/pkg"
	"github.com/spf13/cobra"
)


// netboxUpdateCmd represents the netboxUpdate command
var netboxUpdateCmd = &cobra.Command{
	Use:   "netboxUpdate",
	Short: "Update Netbox with the devices in the folder",
	Long: `Using the AVD structured config and some additional information update Netbox to reflect the devices and their configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("netboxUpdate called")
		folder := cmd.Flag("folder").Value.String()
		log.Print("Folder:", folder)

		url, err := url.Parse(netboxURL)
		cobra.CheckErr(err)
		token, err := createToken(netboxUsername, netboxPassword, url)
		cobra.CheckErr(err)
		app.NetBoxclient = n.NewAPIClientFor(url.Host, token)
		c := app.NetBoxclient.GetConfig()
		c.Scheme = "https"
		// Fetches all the .yml files in the given path
		files, err := getYmlFiles(folder)
		if err != nil {
			fmt.Println(err)
		}
		for _, file := range files {
			fmt.Println("Working on file:", file)
			c, err := mo.NewConfigFromYaml(file)
			if err != nil {
				cobra.CheckErr(err)
			}
			app.AddDevice(c)
			netbox.Work(app)
		}
	},
}

func init() {
	rootCmd.AddCommand(netboxUpdateCmd)
	netboxURL = os.Getenv("NETBOX_URL")
	netboxToken = os.Getenv("NETBOX_TOKEN")
	netboxUsername = os.Getenv("NETBOX_USERNAME")
	netboxPassword = os.Getenv("NETBOX_PASSWORD")
	app = pkg.NewApplication()

	netboxUpdateCmd.Flags().StringP("folder", "f", "", "Path to the folder")
	err := infoUpdateCmd.MarkFlagRequired("folder")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// netboxUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// netboxUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createToken(usr, pwd string, url *url.URL) (string, error) {
	client := resty.New()
	client.SetBaseURL("https://" + url.Host)

	body := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, usr, pwd)

	result := make(map[string]interface{})
	_, err := client.R().
		SetResult(&result).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("/api/users/tokens/provision/")

	if err != nil {
		return "", fmt.Errorf("error requesting a token: %w", err)
	}

	if val, ok := result["key"]; ok {
		return val.(string), nil
	}

	return "", fmt.Errorf("empty token")
}
