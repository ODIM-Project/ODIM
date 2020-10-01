package redfish

type Chassis struct {
	Ocontext           string           `json:"@odata.context"`
	Oid                string           `json:"@odata.id"`
	Otype              string           `json:"@odata.type"`
	Oetag              string           `json:"@odata.etag,omitempty"`
	ID                 string           `json:"Id"`
	Description        string           `json:"Description"`
	Name               string           `json:"Name"`
	AssetTag           string           `json:"AssetTag"`
	ChassisType        string           `json:"ChassisType"`
	DepthMm            int              `json:"DepthMm"`
	EnvironmentalClass string           `json:"EnvironmentalClass"`
	HeightMm           int              `json:"HeightMm"`
	IndicatorLED       int              `json:"IndicatorLED"`
	Manufacturer       string           `json:"Manufacturer"`
	Model              string           `json:"Model"`
	PartNumber         string           `json:"PartNumber"`
	PowerState         string           `json:"PowerState"`
	SerialNumber       string           `json:"SerialNumber"`
	SKU                string           `json:"SKU"`
	UUID               string           `json:"UUID"`
	WeightKg           int              `json:"WeightKg"`
	WidthMm            int              `json:"WidthMm"`
	Links              Links            `json:"Links"`
	Location           Location         `json:"Location"`
	LogServices        LogServices      `json:"LogServices"`
	Assembly           Assembly         `json:"Assembly"`
	NetworkAdapters    NetworkAdapters  `json:"NetworkAdapters"`
	PCIeSlots          PCIeSlots        `json:"PCIeSlots"`
	PhysicalSecurity   PhysicalSecurity `json:"PhysicalSecurity"`
	Power              Power            `json:"Power,omitempty"`
	Sensors            Sensors          `json:"Sensors,omitempty"`
	Status             Status           `json:"Status"`
	Thermal            Thermal          `json:"Thermal,omitempty"`
}

type Location struct {
	Oid string `json:"@odata.id"`
}

type LogServices struct {
	Oid                 string  `json:"@odata.id"`
	Ocontext            string  `json:"@odata.context,omitempty"`
	Otype               string  `json:"@odata.type,omitempty"`
	Oetag               string  `json:"@odata.etag,omitempty"`
	ID                  string  `json:"Id,omitempty"`
	Description         string  `json:"Description,omitempty"`
	Name                string  `json:"Name,omitempty"`
	DateTime            string  `json:"DateTime,omitempty"`
	DateTimeLocalOffset string  `json:"DateTimeLocalOffset,omitempty"`
	Entries             Entries `json:"Entries,omitempty"`
	LogEntryType        string  `json:"LogEntryType,omitempty"`
	MaxNumberOfRecords  int     `json:"MaxNumberOfRecords,omitempty"`
	OverWritePolicy     string  `json:"OverWritePolicy,omitempty"`
	ServiceEnabled      bool    `json:"ServiceEnabled,omitempty"`
	Status              Status  `json:"Status,omitempty"`
}

type Entries struct {
	Oid string `json:"@odata.id"`
}

type Assembly struct {
	Oid string `json:"@odata.id"`
}

type NetworkAdapters struct {
	Oid string `json:"@odata.id"`
}

type PCIeSlots struct {
	Oid string `json:"@odata.id"`
}

type PhysicalSecurity struct {
	IntrusionSensor       string
	IntrusionSensorNumber int
	IntrusionSensorReArm  string
}

type Power struct {
	Oid string `json:"@odata.id"`
}

type Sensors struct {
	Oid string `json:"@odata.id"`
}

type Status struct {
	Oid          string `json:"@odata.id"`
	Ocontext     string `json:"@odata.context,omitempty"`
	Oetag        string `json:"@odata.etag,omitempty"`
	Otype        string `json:"@odata.type,omitempty"`
	Description  string `json:"description,omitempty"`
	ID           string `json:"Id,omitempty"`
	Name         string `json:"Name,omitempty"`
	Health       string `json:"Health,omitempty"`
	HealthRollup string `json:"HealthRollup,omitempty"`
	State        string `json:"State,omitempty"`
	Oem          Oem    `json:"Oem,omitempty"`
}

type Thermal struct {
	Oid string `json:"@odata.id"`
}

type Oem struct {
}

type Links struct {
	Contains                 []Link `json:",omitempty"`
	Chassis                  []Link `json:",omitempty"`
	ComputerSystems          []Link `json:",omitempty"`
	ConsumingComputerSystems []Link `json:",omitempty"`
	ContainedBy              []Link `json:",omitempty"`
	CooledBy                 []Link `json:",omitempty"`
	Endpoints                []Link `json:",omitempty"`
	Drives                   []Link `json:",omitempty"`
	ManagedBy                []Link `json:",omitempty"`
	Oem                      *Oem   `json:",omitempty"`
	ManagersInChassis        []Link `json:",omitempty"`
	PCIeDevices              []Link `json:",omitempty"`
	PCIeFunctions            []Link `json:",omitempty"`
	PoweredBy                []Link `json:",omitempty"`
	Processors               []Link `json:",omitempty"`
	ResourceBlocks           []Link `json:",omitempty"`
	Storage                  []Link `json:",omitempty"`
	SupplyingComputerSystems []Link `json:",omitempty"`
	Switches                 []Link `json:",omitempty"`
}
