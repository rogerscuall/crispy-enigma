package pocketbase

import (
	"encoding/json"
	"fmt"

	"github.com/rogerscuall/crispy-enigma/internal/pocketbase/models"
)

func (c *Client) GetInterfacesFromSwitchName(name string) ([]models.Interface, error) {
	par := Params{
		Page:    1,
		Size:    100,
		Filters: fmt.Sprintf("name=\"%s\"", name),
	}
	//Get the ID of the switch
	resp, err := c.List("switch", par)
	if err != nil {
		return []models.Interface{}, err
	}
	sw := models.SwitchList{}
	err = json.Unmarshal(resp, &sw)
	if err != nil {
		return []models.Interface{}, err
	}
	if len(sw.Items) == 0 {
		return []models.Interface{}, fmt.Errorf("switch %s not found", name)
	} else if len(sw.Items) > 1 {
		return []models.Interface{}, fmt.Errorf("multiple switches found with name %s", name)
	}
	//Get the interfaces of the switch
	par = Params{
		Page:    1,
		Size:    100,
		Filters: fmt.Sprintf("switch=\"%s\"", sw.Items[0].ID),
	}
	resp, err = c.List("interface", par)
	if err != nil {
		return []models.Interface{}, err
	}
	interfaces := models.InterfaceList{}
	err = json.Unmarshal(resp, &interfaces)
	if err != nil {
		return []models.Interface{}, err
	}
	if len(interfaces.Items) == 0 {
		return []models.Interface{}, fmt.Errorf("interfaces not found for switch %s", name)
	}
	return interfaces.Items, nil
}
