package tcommon

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
)

var (
	// ConfigFilePath holds the value of odim config file path
	ConfigFilePath string
)

const (
	// IterationCount is a value that needs to be added in context
	// to track the number of threads created
	IterationCount = "IterationCount"
)

// TrackConfigFileChanges monitors the config changes using fsnotfiy
func TrackConfigFileChanges(errChan chan error) {
	eventChan := make(chan interface{})
	format := config.Data.LogFormat
	go common.TrackConfigFileChanges(ConfigFilePath, eventChan, errChan)
	for {
		select {
		case info := <-eventChan:
			l.Log.Info(info) // new data arrives through eventChan channel
			if l.Log.Level != config.Data.LogLevel {
				l.Log.Info("Log level is updated, new log level is ", config.Data.LogLevel)
				l.Log.Logger.SetLevel(config.Data.LogLevel)
			}
			if format != config.Data.LogFormat {
				l.SetFormatter(config.Data.LogFormat)
				format = config.Data.LogFormat
				l.Log.Info("Log format is updated, new log format is ", config.Data.LogFormat)
			}
		case err := <-errChan:
			l.Log.Error(err)
		}
	}
}
