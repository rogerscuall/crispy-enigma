package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rogerscuall/crispy-enigma/internal/pocketbase"
	"github.com/rogerscuall/crispy-enigma/internal/pocketbase/models"
	"github.com/tmc/langchaingo/tools"
)

var _ tools.Tool = &InterfaceToolVLAN{}

type InterfaceToolVLAN struct {
	Client           *pocketbase.Client
	Interfaces       []models.Interface
	WorkingInterface *WorkingInterfaceTool
}

func (it *InterfaceToolVLAN) Name() string {
	return "interface"
}

func (it *InterfaceToolVLAN) Description() string {
	return "Useful update the vlan of an interface, by providing vlan name as input. Use this tool after using the WorkingInterfaceTool"
}

func (it *InterfaceToolVLAN) Call(ctx context.Context, input string) (string, error) {
	fmt.Println(it.WorkingInterface.InterfaceID)
	fmt.Println(input)
	resp, err := it.Client.List("vlan", pocketbase.Params{Filters: fmt.Sprintf("name=\"%s\"", input)})
	if err != nil {
		return "", err
	}
	vl := models.VLANList{}
	err = json.Unmarshal(resp, &vl)
	if err != nil {
		return "", err
	}
	if len(vl.Items) == 0 {
		return "", fmt.Errorf("vlan %s not found", input)
	} else if len(vl.Items) > 1 {
		return "", fmt.Errorf("multiple vlans found with name %s", input)
	}
	vlanID := vl.Items[0].ID
	resp, err = it.Client.View("interface", it.WorkingInterface.InterfaceID)
	if err != nil {
		return "", err
	}
	interf := models.Interface{}
	err = json.Unmarshal(resp, &interf)
	if err != nil {
		return "", err
	}
	body := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		State       bool   `json:"state"`
		VLAN        string `json:"vlan"`
	}{
		Name: interf.Name,
		VLAN: vlanID,
	}
	err = it.Client.Update("interface", it.WorkingInterface.InterfaceID, body)
	if err != nil {
		return "", err
	}
	return "VLAN updated successfully", nil
}

var _ tools.Tool = &WorkingInterfaceTool{}

type WorkingInterfaceTool struct {
	InterfaceID string
}

func (wit *WorkingInterfaceTool) Name() string {
	return "workinginterface"
}

func (wit *WorkingInterfaceTool) Description() string {
	return "Useful to store the interface ID the user wants to work with. Use this tool before using the InterfaceToolVLAN"
}

func (wit *WorkingInterfaceTool) Call(ctx context.Context, input string) (string, error) {
	fmt.Println("Selecting the working interface as " + input)
	wit.InterfaceID = input
	return "Interface stored successfully", nil
}
