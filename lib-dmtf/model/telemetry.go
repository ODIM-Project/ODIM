//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

package model

type MetricDefinitions struct {
	ODataID                 string             `json:"@odata.id"`
	ODataType               string             `json:"@odata.type"`
	ID                      string             `json:"Id"`
	Name                    string             `json:"Name"`
	Accuracy                int                `json:"Accuracy,omitempty"`
	Calculable              string             `json:"Calculable,omitempty"`
	CalculationAlgorithm    string             `json:"CalculationAlgorithm,omitempty"`
	CalculationParameters   []CalculationParam `json:"CalculationParameters,omitempty"`
	CalculationTimeInterval string             `json:"CalculationTimeInterval,omitempty"`
	Calibration             int                `json:"Calibration,omitempty"`
	DiscreteValues          []string           `json:"DiscreteValues,omitempty"`
	Implementation          string             `json:"Implementation,omitempty"`
	IsLinear                bool               `json:"IsLinear,omitempty"`
	MaxReadingRange         int                `json:"MaxReadingRange,omitempty"`
	MetricDataType          string             `json:"MetricDataType,omitempty"`
	MetricProperties        []string           `json:"MetricProperties,omitempty"`
	MetricType              string             `json:"MetricType,omitempty"`
	MinReadingRange         int                `json:"MinReadingRange,omitempty"`
	OEMCalculationAlgorithm string             `json:"OEMCalculationAlgorithm,omitempty"`
	PhysicalContext         string             `json:"PhysicalContext,omitempty"`
	Precision               int                `json:"Precision,omitempty"`
	SensingInterval         string             `json:"SensingInterval,omitempty"`
	TimestampAccuracy       string             `json:"TimestampAccuracy,omitempty"`
	Units                   string             `json:"Units,omitempty"`
	Wildcards               []WildCard         `json:"Wildcards,omitempty"`
}

type CalculationParam struct {
	ResultMetric string `json:"ResultMetric,omitempty"`
	SourceMetric string `json:"SourceMetric,omitempty"`
}

type MetricReportDefinitions struct {
	ODataID                       string     `json:"@odata.id"`
	ODataType                     string     `json:"@odata.type"`
	ID                            string     `json:"Id"`
	Name                          string     `json:"Name"`
	AppendLimit                   int        `json:"AppendLimit,omitempty"`
	Links                         MetricLink `json:"Links,omitempty"`
	MetricProperties              []string   `json:"MetricProperties,omitempty"`
	MetricReport                  Oid        `json:"MetricReport,omitempty"`
	MetricReportDefinitionEnabled bool       `json:"MetricReportDefinitionEnabled,omitempty"`
	MetricReportDefinitionType    string     `json:"MetricReportDefinitionType,omitempty"`
	MetricReportHeartbeatInterval string     `json:"MetricReportHeartbeatInterval,omitempty"`
	Metrics                       []Metric   `json:"Metrics,omitempty"`
	ReportActions                 []string   `json:"ReportActions,omitempty"`
	ReportTimespan                string     `json:"ReportTimespan,omitempty"`
	ReportUpdates                 string     `json:"ReportUpdates,omitempty"`
	Schedule                      Schedule   `json:"Schedule,omitempty"`
	Status                        Status     `json:"Status,omitempty"`
	SuppressRepeatedMetricValue   bool       `json:"SuppressRepeatedMetricValue,omitempty"`
	Wildcards                     []WildCard `json:"Wildcards,omitempty"`
}

type Schedule struct {
	EnabledDaysOfMonth  []string `json:"EnabledDaysOfMonth,omitempty"`
	EnabledDaysOfWeek   []string `json:"EnabledDaysOfWeek,omitempty"`
	EnabledIntervals    []string `json:"EnabledIntervals,omitempty"`
	EnabledMonthsOfYear []string `json:"EnabledMonthsOfYear,omitempty"`
	InitialStartTime    string   `json:"InitialStartTime,omitempty"`
	Lifetime            string   `json:"Lifetime,omitempty"`
	MaxOccurrences      int      `json:"MaxOccurrences,omitempty"`
	Name                string   `json:"Name,omitempty"`
	RecurrenceInterval  string   `json:"RecurrenceInterval,omitempty"`
}

type Metric struct {
	CollectionDuration  string   `json:"CollectionDuration,omitempty"`
	CollectionFunction  string   `json:"CollectionFunction,omitempty"`
	CollectionTimeScope string   `json:"CollectionTimeScope,omitempty"`
	MetricId            string   `json:"MetricId,omitempty"`
	MetricProperties    []string `json:"MetricProperties,omitempty"`
	Oem                 *Oem     `json:"Oem,omitempty"`
}

type MetricLink struct {
	Oem      *Oem  `json:"Oem,omitempty"`
	Triggers []Oid `json:"Triggers,omitempty"`
}

type Oid struct {
	ODataID string `json:"@odata.id"`
}

type MetricReports struct {
	Context                string        `json:"Context,omitempty"`
	MetricReportDefinition Oem           `json:"MetricReportDefinition,omitempty"`
	MetricValues           []MetricValue `json:"MetricValues,omitempty"`
	ReportSequence         string        `json:"ReportSequence,omitempty"`
	Timestamp              string        `json:"Timestamp,omitempty"`
}

type MetricValue struct {
	MetricDefinition Oem    `json:"MetricDefinition,omitempty"`
	MetricId         string `json:"MetricId,omitempty"`
	MetricProperty   string `json:"MetricProperty,omitempty"`
	MetricValue      string `json:"MetricValue,omitempty"`
	Oem              *Oem   `json:"Oem,omitempty"`
	Timestamp        string `json:"Timestamp,omitempty"`
}

type Triggers struct {
	ODataID                  string            `json:"@odata.id"`
	ODataType                string            `json:"@odata.type"`
	ID                       string            `json:"Id"`
	Name                     string            `json:"Name"`
	DiscreteTriggerCondition string            `json:"DiscreteTriggerCondition,omitempty"`
	DiscreteTriggers         []DiscreteTrigger `json:"DiscreteTriggers,omitempty"`
	EventTriggers            []string          `json:"EventTriggers,omitempty"`
	Links                    TriggerLinks      `json:"Links,omitempty"`
	MetricProperties         []string          `json:"MetricProperties,omitempty"`
	MetricType               string            `json:"MetricType,omitempty"`
	NumericThresholds        NumericThresholds `json:"NumericThresholds,omitempty"`
	Status                   Status            `json:"Status,omitempty"`
	TriggerActions           []string          `json:"TriggerActions,omitempty"`
	Wildcards                []WildCard        `json:"Wildcards,omitempty"`
}

type NumericThresholds struct {
	LowerCritical Threshold `json:"LowerCritical,omitempty"`
	LowerWarning  Threshold `json:"LowerWarning,omitempty"`
	UpperCritical Threshold `json:"UpperCritical,omitempty"`
	UpperWarning  Threshold `json:"UpperWarning,omitempty"`
}
type Threshold struct {
	Activation string `json:"Activation,omitempty"`
	DwellTime  string `json:"DwellTime,omitempty"`
	Reading    int    `json:"Reading,omitempty"`
}

type TriggerLinks struct {
	MetricReportDefinitions []Oid `json:"MetricReportDefinitions,omitempty"`
	Oem                     *Oem  `json:"Oem,omitempty"`
}

type DiscreteTrigger struct {
	DwellTime string `json:"DwellTime,omitempty"`
	Name      string `json:"Name,omitempty"`
	Severity  string `json:"Severity,omitempty"`
	Value     string `json:"Value,omitempty"`
}

type WildCard struct {
	Name   string   `json:"Name,omitempty"`
	Values []string `json:"Values,omitempty"`
}
