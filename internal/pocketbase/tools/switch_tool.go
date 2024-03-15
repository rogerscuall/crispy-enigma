package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rogerscuall/crispy-enigma/internal/pocketbase"
	"github.com/rogerscuall/crispy-enigma/internal/pocketbase/models"
	"github.com/tmc/langchaingo/tools"
)

// SwitchTool needs to implement the Tool interface
var _ tools.Tool = &SwitchTool{}

type SwitchTool struct {
	Client     *pocketbase.Client
	Interfaces []models.Interface
	Vlans      []models.VLAN
}

func (st *SwitchTool) Name() string {
	return "switch"
}
func (st *SwitchTool) Description() string {
	return "Useful to return all the interfaces of a switch using its name as input."
}

func (st *SwitchTool) Call(ctx context.Context, input string) (string, error) {
	interfaces, err := st.Client.GetInterfacesFromSwitchName(input)
	if err != nil {
		return "", err
	}
	st.Interfaces = interfaces
	response := "Interfaces: \n"
	for _, intf := range interfaces {
		resp, err := st.Client.View("vlan", intf.VLAN)
		if err != nil {
			return "", err
		}
		vl := models.VLAN{}
		err = json.Unmarshal(resp, &vl)
		if err != nil {
			return "", err
		}
		response += fmt.Sprintf("Name: %s, ID %s, Description: %s, State: %t, VLAN Name: %s, Switch: %s\n", intf.Name, intf.ID, intf.Description, intf.State, vl, intf.Switch)
	}
	return response, nil
}


