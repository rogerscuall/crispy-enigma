package models

type Switch struct {
	Collection
	Name       string   `json:"name"`
	Interfaces []string `json:"interfaces"`
	Vlans      []string `json:"vlans"`
}

type SwitchList struct {
	CollectionList
	Items []Switch `json:"items"`
}
