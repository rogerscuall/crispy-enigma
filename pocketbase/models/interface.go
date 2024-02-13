package models

type Interface struct {
	ID             string `json:"id"`
	CollectionID   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	Updated        string `json:"updated"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	State          bool   `json:"state"`
	VLANID         string `json:"vlan_id"`
}

type InterfaceErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
