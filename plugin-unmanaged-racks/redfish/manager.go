package redfish

type Manager struct {
	OdataContext    string         `json:"@odata.context"`
	Etag            string         `json:"@odata.etag,omitempty"`
	OdataID         string         `json:"@odata.id"`
	OdataType       string         `json:"@odata.type"`
	Name            string         `json:"Name"`
	ManagerType     string         `json:"ManagerType"`
	ID              string         `json:"Id"`
	UUID            string         `json:"UUID"`
	FirmwareVersion string         `json:"FirmwareVersion"`
	Status          *ManagerStatus `json:"Status,omitempty"`
}

type ManagerStatus struct {
	State string `json:"State"`
}
