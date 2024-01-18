package infoblox

import (
	"fmt"
	"net/netip"

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
)



func CreateConnector(url, wapi, username, password string) (*ibclient.Connector, error) {
	transportConfig := ibclient.NewTransportConfig("false", 20, 10)
	transportConfig.SslVerify = false
	requestBuilder := &ibclient.WapiRequestBuilder{}
	requestor := &ibclient.WapiHttpRequestor{}
	hostConfig := ibclient.HostConfig{
		Scheme:  "https",
		Host:    url,
		Version: wapi,
		Port:    "443",
	}
	authConfig := ibclient.AuthConfig{
		Username: username,
		Password: password,
	}
	conn, err := ibclient.NewConnector(hostConfig, authConfig, transportConfig, requestBuilder, requestor)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func GetNextAvailableNetwork(objMgr ibclient.IBObjectManager, netview, containerName string, mask uint) (string, error) {
	container, err := netip.ParsePrefix(containerName)
	if err != nil {
		return "", err
	}
	if mask > 32 {
		return "", fmt.Errorf("invalid mask: %d", mask)
	}
	_, err = objMgr.GetNetworkContainer("default", container.String(), false, nil)
	if err != nil {
		return "", err
	}
	extraAttrs := ibclient.EA{
		"Simple":  "Network",
		"Comment": "This is a comment",
	}
	var network *ibclient.Network

	network, err = objMgr.AllocateNetwork("default", containerName, false, mask, "used", extraAttrs)
	if err != nil {
		return "", err
	}
	return network.Cidr, nil
}
