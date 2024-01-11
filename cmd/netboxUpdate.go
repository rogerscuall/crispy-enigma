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
	"os"

	n "github.com/netbox-community/go-netbox/v3/netbox"
	"github.com/rogerscuall/crispy-enigma/netbox"
	"github.com/spf13/cobra"
)

var (
	//netboxURL      string
	netboxUsername string
	netboxPassword string
	netboxToken    string
	netboxURL      string
)

// netboxUpdateCmd represents the netboxUpdate command
var netboxUpdateCmd = &cobra.Command{
	Use:   "netboxUpdate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("netboxUpdate called")
		nb := n.NewNetboxWithAPIKey(netboxURL, netboxToken)
		netbox.Work(nb)
		// role := int64(1)
		// name := "test-device-rgo"
		// deviceType := int64(8)
		// site := int64(1)
		// device := &models.WritableDeviceWithConfigContext{
		// 	Name:       &name,
		// 	Role:       &role,
		// 	DeviceType: &deviceType,
		// 	DeviceRole: &role,
		// 	Site:       &site,
		// 	Tags:       []*models.NestedTag{},
		// }
		// create, err := nb.Dcim.DcimDevicesCreate(&dcim.DcimDevicesCreateParams{
		// 	Context: context.Background(),
		// 	Data:    device,
		// }, nil)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Printf("%#v\n", create)
	},
}

func init() {
	rootCmd.AddCommand(netboxUpdateCmd)

	netboxURL = os.Getenv("NETBOX_URL")
	netboxToken = os.Getenv("NETBOX_TOKEN")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// netboxUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// netboxUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// func createToken(usr, pwd string, url *url.URL) (string, error) {
// 	client := resty.New()
// 	client.SetBaseURL("https://" + url.Host)

// 	body := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, usr, pwd)

// 	result := make(map[string]interface{})
// 	_, err := client.R().
// 		SetResult(&result).
// 		SetHeader("Content-Type", "application/json").
// 		SetBody(body).
// 		Post("/api/users/tokens/provision/")

// 	if err != nil {
// 		return "", fmt.Errorf("error requesting a token: %w", err)
// 	}

// 	if val, ok := result["key"]; ok {
// 		return val.(string), nil
// 	}

// 	return "", fmt.Errorf("empty token")
// }
