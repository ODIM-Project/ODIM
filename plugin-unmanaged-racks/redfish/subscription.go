package redfish

type EvtSubPost struct {
	Name                 string   `json:"Name"`
	Destination          string   `json:"Destination"`
	EventTypes           []string `json:"EventTypes,omitempty"`
	MessageIds           []string `json:"MessageIds,omitempty"`
	ResourceTypes        []string `json:"ResourceTypes,omitempty"`
	Context              string   `json:"Context"`
	Protocol             string   `json:"Protocol"`
	SubscriptionType     string   `json:"SubscriptionType"`
	EventFormatType      string   `json:"EventFormatType"`
	SubordinateResources bool     `json:"SubordinateResources"`
	OriginResources      []Link   `json:"OriginResources"`
}
