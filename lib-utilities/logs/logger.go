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
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var priorityLogFields = []string{
	"host",
	"thread_name",
	"process_id",
	"message_id",
}

var syslogPriorityNumerics = map[string]int8{
	"panic":   8,
	"fatal":   9,
	"error":   11,
	"warn":    12,
	"warning": 12,
	"info":    14,
	"debug":   15,
	"trace":   15,
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

// SysLogFormatter implements logrus Format interface. It provides a formatter for odim in syslog format
type SysLogFormatter struct{}

// Log is an instance of logrus.Entry.
// This instance should be initialized at the start of the each service
// and will be used for logging
var Log *logrus.Entry

// LogFormat is custom type created for the log formats supported by ODIM
type LogFormat uint32

// log formats supported by ODIM
const (
	SyslogFormat LogFormat = iota
	JSONFormat
)

func getProcessLogDetails(ctx context.Context) logrus.Fields {
	var fields = make(map[string]interface{})
	TransactionID := strings.Replace(fmt.Sprintf("%v", ctx.Value("transactionid")), "\n", "", -1)
	ProcessName := strings.Replace(fmt.Sprintf("%v", ctx.Value("processname")), "\n", "", -1)
	ThreadID := strings.Replace(fmt.Sprintf("%v", ctx.Value("threadid")), "\n", "", -1)
	ActionName := strings.Replace(fmt.Sprintf("%v", ctx.Value("actionname")), "\n", "", -1)
	ThreadName := strings.Replace(fmt.Sprintf("%v", ctx.Value("threadname")), "\n", "", -1)
	ActionID := strings.Replace(fmt.Sprintf("%v", ctx.Value("actionid")), "\n", "", -1)
	fields["transaction_id"] = TransactionID
	fields["process_name"] = ProcessName
	fields["thread_id"] = ThreadID
	fields["action_name"] = ActionName
	fields["message_id"] = ActionName
	fields["thread_name"] = ThreadName
	fields["action_id"] = ActionID

	return fields
}

func init() {
	Log = logrus.NewEntry(logrus.New())
}

// Adorn adds the fields to Log variable
func Adorn(m logrus.Fields) {
	Log = Log.WithFields(m)
}

// LogWithFields add fields to log
func LogWithFields(ctx context.Context) *logrus.Entry {
	fields := getProcessLogDetails(ctx)
	return Log.WithFields(fields)
}

// Format renders a log in syslog format
func (f *SysLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	level := entry.Level.String()
	priorityNumber := findSysLogPriorityNumeric(entry, level)
	sysLogMsg := fmt.Sprintf("<%d>%s %s", priorityNumber, "1", entry.Time.UTC().Format(time.RFC3339))
	sysLogMsg = formatPriorityFields(entry, sysLogMsg)
	sysLogMsg = formatStructuredFields(entry, sysLogMsg)
	for k, v := range logFields {
		if accountLog, present := formatSyslog(k, v, entry); present {
			sysLogMsg = fmt.Sprintf("%s %s", sysLogMsg, accountLog)
		}
	}
	if _, ok := entry.Data["auth"]; ok {
		sysLogMsg = formatAuthStructFields(entry, sysLogMsg, priorityNumber)
	}
	if _, ok := entry.Data["audit"]; ok {
		sysLogMsg = formatAuditStructFields(entry, sysLogMsg, priorityNumber)
		return append([]byte(sysLogMsg), '\n'), nil
	}

	sysLogMsg = fmt.Sprintf("%s %s", sysLogMsg, entry.Message)
	return append([]byte(sysLogMsg), '\n'), nil
}

// findSysLogPriorityNumeric is used to find the log priority number
func findSysLogPriorityNumeric(entry *logrus.Entry, level string) int8 {
	if _, ok := entry.Data["auth"].(bool); ok {
		sCode := entry.Data["statuscode"].(int32)
		if getResponseStatus(sCode) {
			return 86
		}
		return 84
	}
	if _, ok := entry.Data["audit"].(bool); ok {
		sCode := entry.Data["statuscode"].(int32)
		if getResponseStatus(sCode) {
			return 110
		}
		return 107
	}
	return syslogPriorityNumerics[level]
}

func formatPriorityFields(entry *logrus.Entry, msg string) string {
	for _, v := range priorityLogFields {
		if val, ok := entry.Data[v]; ok {
			msg = fmt.Sprintf("%s %v", msg, val)
		}
	}

	return msg
}

// formatAuthStructFields used to format the syslog message
func formatAuthStructFields(entry *logrus.Entry, msg string, priorityNo int8) string {
	var sessionUserName, sessionRoleID string
	respStatusCode := int32(http.StatusUnauthorized)
	tokenMsg := ""

	if entry.Data["sessionuserid"] != nil {
		sessionUserName = entry.Data["sessionuserid"].(string)
	}
	if entry.Data["sessionroleid"] != nil {
		sessionRoleID = entry.Data["sessionroleid"].(string)
	}
	if entry.Data["statuscode"] != nil {
		respStatusCode = entry.Data["statuscode"].(int32)
	}
	msg = fmt.Sprintf("%s [account@1 user=\"%s\" roleID=\"%s\"]", msg, sessionUserName, sessionRoleID)
	if priorityNo == 86 {
		msg = fmt.Sprintf("%s %s", msg, "Authentication/Authorization successful")
	} else {
		errMsg := "Authentication/Authorization failed"
		if respStatusCode == http.StatusForbidden {
			errMsg = "Authorization failed"
		} else if respStatusCode == http.StatusUnauthorized {
			errMsg = "Authentication failed"
		}
		msg = fmt.Sprintf("%s %s %s", msg, errMsg, tokenMsg)
	}

	return msg
}

// formatStructuredFields is used to create structured fields for log
func formatStructuredFields(entry *logrus.Entry, msg string) string {
	var transID, processName, actionID, actionName, threadID, threadName string
	fields := []string{"process_name", "transaction_id", "action_id", "action_name", "thread_id", "thread_name"}

	for _, value := range fields {
		if val, ok := entry.Data[value]; ok {
			if val != nil {
				switch value {
				case "process_name":
					processName = val.(string)
				case "transaction_id":
					transID = val.(string)
				case "action_id":
					actionID = val.(string)
				case "action_name":
					actionName = val.(string)
				case "thread_id":
					threadID = val.(string)
				case "thread_name":
					threadName = val.(string)
				}
			}
		}
	}
	if transID != "" {
		msg = fmt.Sprintf("%s [process@1 processName=\"%s\" transactionID=\"%s\" actionID=\"%s\" actionName=\"%s\" threadID=\"%s\" threadName=\"%s\"]", msg, processName, transID, actionID, actionName, threadID, threadName)
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

// Convert the log format to a string.
func (format LogFormat) String() string {
	if b, err := format.MarshalText(); err == nil {
		return string(b)
	}
	return "unknown_log_format"
}

// ParseLogFormat takes a string level and returns the log format.
func ParseLogFormat(format string) (LogFormat, error) {
	switch strings.ToLower(format) {
	case "syslog":
		return SyslogFormat, nil
	case "json":
		return JSONFormat, nil
	}

	var lf LogFormat
	return lf, fmt.Errorf("invalid log format : %s", format)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (format *LogFormat) UnmarshalText(text []byte) error {
	l, err := ParseLogFormat(string(text))
	if err != nil {
		return err
	}

	*format = l
	return nil
}

// MarshalText will validate the log format and return the corresponding string
func (format LogFormat) MarshalText() ([]byte, error) {
	switch format {
	case SyslogFormat:
		return []byte("syslog"), nil
	case JSONFormat:
		return []byte("json"), nil
	}

	return nil, fmt.Errorf("invalid log format %d", format)
}

// SetFormatter set the format for logging
func SetFormatter(format LogFormat) {
	switch format {
	case SyslogFormat:
		Log.Logger.SetFormatter(&SysLogFormatter{})
	case JSONFormat:
		Log.Logger.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "log_timestamp",
			logrus.FieldKeyLevel: "log_level",
			logrus.FieldKeyMsg:   "log_message",
		}})
	default:
		Log.Logger.SetFormatter(&SysLogFormatter{})
		Log.Info("Something went wrong! Setting the default format Syslog for logging")
	}
}
