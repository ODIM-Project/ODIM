package redfish

type MessageData struct {
	OdataType string  `json:"@odata.type"`
	Name      string  `json:"Name"`
	Context   string  `json:"@odata.context"`
	Events    []Event `json:"Events"`
}

// Event contains the details of the event subscribed from PMB
type Event struct {
	MemberID          string      `json:"MemberId,omitempty"`
	EventType         string      `json:"EventType"`
	EventGroupID      int         `json:"EventGroupId,omitempty"`
	EventID           string      `json:"EventId"`
	Severity          string      `json:"Severity"`
	EventTimestamp    string      `json:"EventTimestamp"`
	Message           string      `json:"Message"`
	MessageArgs       []string    `json:"MessageArgs,omitempty"`
	MessageID         string      `json:"MessageId"`
	Oem               interface{} `json:"Oem,omitempty"`
	OriginOfCondition *Link       `json:"OriginOfCondition,omitempty"`
}
