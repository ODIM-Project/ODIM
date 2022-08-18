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
	"github.com/sirupsen/logrus"
)

// Logger is used when you import the log package in your service
// and logging using the package level methods.
// This can be customized using the functions InitSysLogger or InitJSONLogger
var Logger *logrus.Entry

// LogFormat type
type LogFormat uint32

const (
	// SysLogFormat will choose constomized ODIMSysLogFormatter for logging
	SysLogFormat LogFormat = iota
	// JSONFormat will choose built in JSON format for logging
	JSONFormat
)

// Config is used for user configuration
type Config struct {
	LogFormat LogFormat
}

// InitLogger sets up the Logger and sets up the format and level
func InitLogger(c *Config) {
	Logger = logrus.NewEntry(logrus.New())

	// setting logger format
	switch c.LogFormat {
	case SysLogFormat:
		Logger.Logger.SetFormatter(&ODIMSysLogFormatter{})
	case JSONFormat:
		Logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		Logger.Logger.SetFormatter(&ODIMSysLogFormatter{})
	}

	// set the minimum level for logging
	Logger.Logger.SetLevel(logrus.DebugLevel)
}

// Trace calls Trace method on the package level
func Trace(args ...interface{}) {
	Logger.Trace(args...)
}

// Debug calls Debug method on the package level
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Print calls Print method on the package level
func Print(args ...interface{}) {
	Logger.Print(args...)
}

// Info calls Info method on the package level
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Warn calls Warn method on the package level
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Warning calls Warning method on the package level
func Warning(args ...interface{}) {
	Logger.Warning(args...)
}

// Error calls Error method on the package level
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Fatal calls Fatal method on the package level
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Panic calls Panic method on the package level
func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

// Tracef calls Tracef method on the package level
func Tracef(format string, args ...interface{}) {
	Logger.Tracef(format, args...)
}

// Debugf calls Debugf method on the package level
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Infof calls Infof method on the package level
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Printf calls Printf method on the package level
func Printf(format string, args ...interface{}) {
	Logger.Printf(format, args...)
}

// Warnf calls Warnf method on the package level
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Warningf calls Warningf method on the package level
func Warningf(format string, args ...interface{}) {
	Logger.Warningf(format, args...)
}

// Errorf calls Errorf method on the package level
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatalf calls Fatalf method on the package level
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// Panicf calls Panicf method on the package level
func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}

// Traceln calls Traceln method on the package level
func Traceln(args ...interface{}) {
	Logger.Traceln(args...)
}

// Debugln calls Debugln method on the package level
func Debugln(args ...interface{}) {
	Logger.Debugln(args...)
}

// Infoln calls Infoln method on the package level
func Infoln(args ...interface{}) {
	Logger.Infoln(args...)
}

// Println calls Println method on the package level
func Println(args ...interface{}) {
	Logger.Println(args...)
}

// Warnln calls Warnln method on the package level
func Warnln(args ...interface{}) {
	Logger.Warnln(args...)
}

// Warningln calls Warningln method on the package level
func Warningln(args ...interface{}) {
	Logger.Warningln(args...)
}

// Errorln calls Errorln method on the package level
func Errorln(args ...interface{}) {
	Logger.Errorln(args...)
}

// Fatalln calls Fatalln method on the package level
func Fatalln(args ...interface{}) {
	Logger.Fatalln(args...)
}

// Panicln calls Panicln method on the package level
func Panicln(args ...interface{}) {
	Logger.Panicln(args...)
}

// TraceWithFileds calls Trace method on package level after appending the fields passed
func TraceWithFileds(fields map[string]interface{}, args ...interface{}) {
	data := getFields(fields)
	Logger.WithFields(data).Trace(args...)
}

// DebugWithFileds calls Debug method on package level after appending the fields passed
func DebugWithFileds(fields map[string]interface{}, args ...interface{}) {
	data := getFields(fields)
	Logger.WithFields(data).Debug(args...)
}

// InfoWithFileds calls Info method on package level after appending the fields passed
func InfoWithFileds(fields map[string]interface{}, args ...interface{}) {
	data := getFields(fields)
	Logger.WithFields(data).Info(args...)
}

// PrintWithFileds calls Info method on package level after appending the fields passed
func PrintWithFileds(fields map[string]interface{}, args ...interface{}) {
	data := getFields(fields)
	Logger.WithFields(data).Info(args...)
}

// WarnWithFileds calls Warn method on package level after appending the fields passed
func WarnWithFileds(fields map[string]interface{}, args ...interface{}) {
	data := getFields(fields)
	Logger.WithFields(data).Warn(args...)
}

// ErrorWithFileds calls Error method on package level after appending the fields passed
func ErrorWithFileds(fields map[string]interface{}, args ...interface{}) {
	data := getFields(fields)
	Logger.WithFields(data).Error(args...)
}

// PanicWithFileds calls Panic method on package level after appending the fields passed
func PanicWithFileds(fields map[string]interface{}, args ...interface{}) {
	data := getFields(fields)
	Logger.WithFields(data).Panic(args...)
}

// getFields converts map[string]interface{} to logrus.Fields
func getFields(fields map[string]interface{}) logrus.Fields {
	data := make(logrus.Fields)
	for k, v := range fields {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	return data
}
