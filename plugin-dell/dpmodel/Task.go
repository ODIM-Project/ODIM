/*
 * Copyright (c) 2021 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
