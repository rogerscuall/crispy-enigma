package pkg

import (
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	n "github.com/netbox-community/go-netbox/v3"
	"github.com/rogerscuall/crispy-enigma/models"
)

type Application struct {
	Devices        []*models.Config
	NetBoxclient   *n.APIClient
	InfobloxClient *ibclient.Connector
}

// func NewApplication() *Application {
// 	//WARNING: this is an extremely important command, do not remove.
// 	// Configure the schema to be https.
	

// 	return &Application{
// 		Devices: make([]*models.Config, 0),
// 	}
// }

func (a *Application) AddDevice(device *models.Config) {
	if a.Devices == nil {
		a.Devices = make([]*models.Config, 0)
	}
	a.Devices = append(a.Devices, device)
}
