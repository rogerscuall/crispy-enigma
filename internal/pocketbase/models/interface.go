package models

type Interface struct {
	Collection
	Name        string `json:"name"`
	Description string `json:"description"`
	State       bool   `json:"state"`
	VLAN        string `json:"vlan"`
	Switch      string `json:"switch"`
}

type InterfaceList struct {
	CollectionList
	Items []Interface `json:"items"`
}
