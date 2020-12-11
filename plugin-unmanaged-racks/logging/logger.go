/*
 * Copyright (c) 2020 Intel Corporation
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

package logging

import (
	"github.com/antonfisher/nested-logrus-formatter"
	"github.com/kataras/golog"
	"github.com/sirupsen/logrus"
)

var staticLogger *logger

type logger struct {
	logLevel *logrus.Level
	*logrus.Logger
}

func (l logger) Print(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l logger) Println(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l logger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l logger) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}

func (l logger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l logger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

func (l logger) SetLevel(lls string) {
	ll, err := logrus.ParseLevel(lls)
	if err != nil {
		l.Errorf("Cannot change log level to %s, defaulting to INFO", lls)
		l.Logger.SetLevel(logrus.InfoLevel)

	}
	l.Logger.SetLevel(ll)
}

// SetLogLevel sets log level on URP's logger
func SetLogLevel(logLevel string) {
	staticLogger.SetLevel(logLevel)
}

// GetLogger returns instance of URP logger
func GetLogger() golog.ExternalLogger {
	return *staticLogger
}

// Error logs error message using URP logger
func Error(i ...interface{}) {
	staticLogger.Error(i...)
}

// Errorf logs error message using URP logger
func Errorf(t string, i ...interface{}) {
	staticLogger.Errorf(t, i...)
}

// Warn logs warning message using URP logger
func Warn(i ...interface{}) {
	staticLogger.Warn(i...)
}

// Warnf logs warning message using URP logger
func Warnf(t string, i ...interface{}) {
	staticLogger.Warnf(t, i...)
}

// Info logs info message using URP logger
func Info(i ...interface{}) {
	staticLogger.Info(i...)
}

// Infof logs info message using URP logger
func Infof(t string, i ...interface{}) {
	staticLogger.Infof(t, i...)
}

// Debug logs debug message using URP logger
func Debug(i ...interface{}) {
	staticLogger.Debug(i...)
}

// Debugf logs debug message using URP logger
func Debugf(t string, i ...interface{}) {
	staticLogger.Debugf(t, i...)
}

// Fatal logs fatal message using URP logger
func Fatal(i ...interface{}) {
	staticLogger.Fatal(i...)
}

func init() {
	ll := logrus.DebugLevel
	l := logrus.New()
	l.SetFormatter(&formatter.Formatter{})

	staticLogger = &logger{
		logLevel: &ll,
		Logger:   l,
	}
}
