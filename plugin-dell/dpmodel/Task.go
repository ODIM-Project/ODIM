package dpmodel

type Task struct {
	ID          string `json:"Id"`
	Description string `json:"Description"`
	Name        string `json:"Name"`
	StartTime   string `json:"StartTime"`
	TaskState   string `json:"TaskState"`
	TaskStatus  string `json:"TaskStatus"`
	EndTime     string
	Messages    []*Message `json:"Messages"`
	Oem         OemTask    `json:"Oem"`
}

// Message Model
type Message struct {
	Message           string   `json:"Message"`
	MessageID         string   `json:"MessageId"`
	MessageArgs       []string `json:"MessageArgs"`
	Oem               OemTask  `json:"Oem"`
	RelatedProperties []string `json:"RelatedProperties"`
	Resolution        string   `json:"Resolution"`
	Severity          string   `json:"Severity"`
}

// Oem Model
type OemTask struct {
	Dell DellTask `json:"Dell"`
}

type DellTask struct {
	CompletionTime    string   `json:"CompletionTime"`
	Description       string   `json:"Description"`
	EndTime           string   `json:"EndTime"`
	Id                string   `json:"Id"`
	JobState          string   `json:"JobState"`
	JobType           string   `json:"JobType"`
	Message           string   `json:"Message"`
	MessageArgs       []string `json:"MessageArgs"`
	MessageId         string   `json:"MessageId"`
	Name              string   `json:"Name"`
	PercentComplete   int      `json:"PercentComplete"`
	StartTime         string   `json:"StartTime"`
	TargetSettingsURI string   `json:"TargetSettingsURI"`
}
