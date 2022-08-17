//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

package logs

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var priorityLogFields = []string{
	"host",
}

var logFields = map[string][]string{
	"account": {
		"user",
		"roleID",
	},
	"request": {
		"method",
		"resource",
		"requestBody",
	},
	"response": {
		"responseCode",
	},
}

// ODIMSysLogFormatter implements logrus Format interface. It provides a formatter for odim in syslog format
type ODIMSysLogFormatter struct{}

// Format renders a log in syslog format
func (f *ODIMSysLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	sysLogMsg := fmt.Sprintf("%s %s", entry.Time.UTC().Format(time.RFC3339), entry.Level)
	sysLogMsg = formatPriorityFields(entry, sysLogMsg)
	for k, v := range logFields {
		if accountLog, present := formatSyslog(k, v, entry); present {
			sysLogMsg = fmt.Sprintf("%s %s", sysLogMsg, accountLog)
		}
	}

	sysLogMsg = fmt.Sprintf("%s %s", sysLogMsg, entry.Message)
	return append([]byte(sysLogMsg), '\n'), nil
}

func formatPriorityFields(entry *logrus.Entry, msg string) string {
	present := true
	for _, v := range priorityLogFields {
		if val, ok := entry.Data[v]; ok {
			present = false
			msg = fmt.Sprintf("%s %v ", msg, val)
		}
	}
	if !present {
		msg = msg[:len(msg)-1]
	}
	return msg
}

func formatSyslog(logType string, logFields []string, entry *logrus.Entry) (string, bool) {
	isPresent := false
	msg := fmt.Sprintf("[%s@1 ", logType)
	for _, v := range logFields {
		if val, ok := entry.Data[v]; ok {
			isPresent = true
			msg = fmt.Sprintf("%s %s=\"%v\" ", msg, v, val)
		}
	}
	msg = msg[:len(msg)-1]
	return fmt.Sprintf("%s]", msg), isPresent
}
