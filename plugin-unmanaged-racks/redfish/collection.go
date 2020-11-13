package redfish

type Collection struct {
	OdataContext string `json:"@odata.context"`
	Etag         string `json:"@odata.etag,omitempty"`
	OdataID      string `json:"@odata.id"`
	OdataType    string `json:"@odata.type"`
	Description  string `json:"Description"`
	Name         string `json:"Name"`
	Members      []Link `json:"Members"`
	MembersCount int    `json:"Members@odata.count"`
}

func NewCollection(odataId, odataType string, members ...Link) Collection {
	return Collection{
		OdataContext: "/ODIM/v1/$metadata" + odataType,
		OdataID:      odataId,
		OdataType:    odataType,
		Members:      append([]Link{}, members...),
		MembersCount: len(members),
	}
}
