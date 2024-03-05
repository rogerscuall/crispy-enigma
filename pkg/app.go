package pkg

import (
	"fmt"

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	n "github.com/netbox-community/go-netbox/v3"
	"github.com/rogerscuall/crispy-enigma/internal/pocketbase"
	"github.com/rogerscuall/crispy-enigma/models"
	"gopkg.in/aristanetworks/go-cvprac.v2/client"
)

type Application struct {
	Devices        []*models.Config
	NetBoxclient   *n.APIClient
	InfobloxClient *ibclient.Connector
	CVPClient      *client.CvpClient
	Debug          bool
	PbClient       *pocketbase.Client
}

func NewApplication() *Application {
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

func (a *Application) DebugLog(format string, v ...interface{}) {
	if a.Debug {
		fmt.Printf(format, v...)
		fmt.Println()
	}
}
