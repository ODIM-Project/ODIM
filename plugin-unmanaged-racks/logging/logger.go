package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var staticLogger *zap.SugaredLogger

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

type bridge struct {
	*zap.SugaredLogger
}

func (b bridge) Print(i ...interface{}) {
	b.SugaredLogger.Info(i)
}

func (b bridge) Println(i ...interface{}) {
	b.SugaredLogger.Info(i)
}

func (b bridge) Error(i ...interface{}) {
	b.SugaredLogger.Error(i)
}

func (b bridge) Warn(i ...interface{}) {
	b.SugaredLogger.Warn(i)
}

func (b bridge) Info(i ...interface{}) {
	b.SugaredLogger.Info(i)
}

func (b bridge) Debug(i ...interface{}) {
	b.SugaredLogger.Debug(i)
}

func init() {
	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "T",
			LevelKey:      "L",
			NameKey:       "N",
			CallerKey:     "C",
			FunctionKey:   "",
			MessageKey:    "M",
			StacktraceKey: "",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(l.CapitalString()[:3])
			},
			EncodeTime:       zapcore.ISO8601TimeEncoder,
			EncodeDuration:   zapcore.StringDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: " ",
		},
	}

	l, e := cfg.Build(zap.AddCallerSkip(1))
	if e != nil {
		panic(e)
	}

	staticLogger = l.Sugar()
}
