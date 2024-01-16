package pkg

import (
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	"github.com/netbox-community/go-netbox/v3/netbox/client"
	"github.com/rogerscuall/crispy-enigma/models"
)

type Application struct {
	Devices        []*models.Config
	NetBoxclient   *client.NetBoxAPI
	InfobloxClient *ibclient.Connector
}

func NewApplication() *Application {
	//WARNING: this is an extremely important command, do not remove.
	client.DefaultSchemes = []string{"https"}
	return &Application{
		Devices: make([]*models.Config, 0),
	}
}

func (a *Application) AddDevice(device *models.Config) {
	if a.Devices == nil {
		a.Devices = make([]*models.Config, 0)
	}
	a.Devices = append(a.Devices, device)
}
