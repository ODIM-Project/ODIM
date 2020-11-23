package logging

import (
	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var staticLogger *logger

type logger struct {
	logLevel *logrus.Level
	*logrus.Logger
}

func (l logger) Print(args ...interface{}) {
	l.Logger.Info(args)
}

func (l logger) Println(args ...interface{}) {
	l.Logger.Info(args)
}

func (l logger) Error(args ...interface{}) {
	l.Logger.Error(args)
}

func (l logger) Warn(args ...interface{}) {
	l.Logger.Warn(args)
}

func (l logger) Info(args ...interface{}) {
	l.Logger.Info(args)
}

func (l logger) Debug(args ...interface{}) {
	l.Logger.Debug(args)
}

func (l logger) SetLevel(lls string) {
	ll, err := logrus.ParseLevel(lls)
	if err != nil {
		l.Errorf("Cannot change log level to %s, defaulting to INFO", lls)
		l.Logger.SetLevel(logrus.InfoLevel)

	}
	l.Logger.SetLevel(ll)
}

func Logger() logger {
	return *staticLogger
}

func Error(i ...interface{}) {
	staticLogger.Error(i)
}

func Errorf(t string, i ...interface{}) {
	staticLogger.Errorf(t, i...)
}

func Warn(i ...interface{}) {
	staticLogger.Warn(i)
}

func Warnf(t string, i ...interface{}) {
	staticLogger.Warnf(t, i...)
}

func Info(i ...interface{}) {
	staticLogger.Info(i)
}

func Infof(t string, i ...interface{}) {
	staticLogger.Infof(t, i...)
}

func Debug(i ...interface{}) {
	staticLogger.Debug(i)
}

func Debugf(t string, i ...interface{}) {
	staticLogger.Debugf(t, i...)
}

func Fatal(i ...interface{}) {
	staticLogger.Fatal(i)
}

func init() {
	ll := logrus.DebugLevel
	l := logrus.New()
	l.SetFormatter(&formatter.Formatter{
		//FieldsOrder:           nil,
		//TimestampFormat:       "",
		//HideKeys:              false,
		//NoColors:              false,
		//NoFieldsColors:        false,
		//NoFieldsSpace:         false,
		//ShowFullLevel:         false,
		//NoUppercaseLevel:      false,
		//TrimMessages:          false,
		//CallerFirst:           false,
		//CustomCallerFormatter: nil,
	})

	staticLogger = &logger{
		logLevel: &ll,
		Logger:   l,
	}
}
