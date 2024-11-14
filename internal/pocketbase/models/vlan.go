package models

type VLAN struct {
	Collection
	Name        string `json:"name"`
	Description string `json:"description"`
}

type VLANList struct {
	CollectionList
	Items []VLAN `json:"items"`
}
