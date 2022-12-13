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

func TrackConfigFileChanges(errChan chan error) {
	eventChan := make(chan interface{})
	go common.TrackConfigFileChanges(ConfigFilePath, eventChan, errChan)
	for {
		l.Log.Info(<-eventChan) // new data arrives through eventChan channel
		if l.Log.Level != config.Data.LogLevel {
			l.Log.Info("Log level is updated, new log level is ", config.Data.LogLevel)
			l.Log.Logger.SetLevel(config.Data.LogLevel)
		}

	}
}

// TrackConfigErrors monitors the errors in goroutines of odim libraries and log the errors
func TrackConfigErrors(errChan chan error) {
	for {
		err := <-errChan
		l.Log.Error(err)
	}
}
