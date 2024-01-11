package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	"github.com/rogerscuall/crispy-enigma/models"
)

func isDataConflictError(err error) bool {
	return strings.Contains(err.Error(), "IB.Data.Conflict")
}

var (
	infoURL         string
	infoUsername    string
	infoPassword    string
	infoWapiVersion string
)

func init() {
	// Initialize global variables
	infoURL = os.Getenv("INFOBLOX_URL")
	infoUsername = os.Getenv("INFOBLOX_USERNAME")
	infoPassword = os.Getenv("INFOBLOX_PASSWORD")
	infoWapiVersion = os.Getenv("INFOBLOX_WAPI_VERSION")
}

func createConnector() (*ibclient.Connector, error) {
	transportConfig := ibclient.NewTransportConfig("false", 20, 10)
	transportConfig.SslVerify = false
	requestBuilder := &ibclient.WapiRequestBuilder{}
	requestor := &ibclient.WapiHttpRequestor{}
	hostConfig := ibclient.HostConfig{
		Scheme:  "https",
		Host:    infoURL,
		Version: infoWapiVersion,
		Port:    "443",
	}
	authConfig := ibclient.AuthConfig{
		Username: infoUsername,
		Password: infoPassword,
	}
	conn, err := ibclient.NewConnector(hostConfig, authConfig, transportConfig, requestBuilder, requestor)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {
	// Create a new connector
	conn, err := createConnector()
	if err != nil {
		panic(err)
	}
	defer conn.Logout()
	objMgr := ibclient.NewObjectManager(conn, "myclient", "")

	// Fetches all the .yml files in the given path
	files, err := getYmlFiles("../old")
	if err != nil {
		fmt.Println(err)
	}
	// update the information in infloblox
	for _, file := range files {
		fmt.Println("Working on file:", file)
		config, err := models.NewConfigFromYaml(file)
		if err != nil {
			log.Println("Skipping file:", file)
			log.Printf("Error: %v", err)
			continue
		}
		//update to the vlan interfaces
		for _, svi := range config.VlanInterfaces {
			log.Println("Working on SVI:", svi.Name)
			networkCIDR := svi.IpAddress
			if networkCIDR == "" {
				networkCIDR = svi.IpAddressVirtual
			}
			_, ipNet, err := net.ParseCIDR(networkCIDR)
			if err != nil {
				fmt.Println(err)
				return
			}
			log.Println(ipNet.String())
			log.Println("Network to create:", ipNet.String())
			extraAttrs := ibclient.EA{
				"Simple":  "Network",
				"Comment": "This is a comment",
			}
			infobloxNetwork, err := objMgr.CreateNetwork("default", ipNet.String(), false, svi.Description, extraAttrs)
			if err != nil {
				if isDataConflictError(err) {
					fmt.Println("Network already exists")
				} else {
					fmt.Println("Error creating network:", err)
				}
				continue
			}
			fmt.Println("Network Created: ", infobloxNetwork.Cidr)
		}
	}

}

// getYmlFiles returns a slice of all the .yml files in the given path
func getYmlFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".yml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
