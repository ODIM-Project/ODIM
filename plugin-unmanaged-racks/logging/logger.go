package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var staticLogger *logger

type logger struct {
	logLevel *zap.AtomicLevel
	*zap.SugaredLogger
}

func (l logger) Print(args ...interface{}) {
	l.SugaredLogger.Info(args)
}

func (l logger) Println(args ...interface{}) {
	l.SugaredLogger.Info(args)
}

func (l logger) Error(args ...interface{}) {
	l.SugaredLogger.Error(args)
}

func (l logger) Warn(args ...interface{}) {
	l.SugaredLogger.Warn(args)
}

func (l logger) Info(args ...interface{}) {
	l.SugaredLogger.Info(args)
}

func (l logger) Debug(args ...interface{}) {
	l.SugaredLogger.Debug(args)
}

func (l logger) SetLevel(ll string) {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(ll))
	if err != nil {
		l.Info("Cannot change log level to %s", ll)
	}
	l.logLevel.SetLevel(level)
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
	ll := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:          "T",
			LevelKey:         "L",
			NameKey:          "N",
			CallerKey:        "C",
			FunctionKey:      "",
			MessageKey:       "M",
			StacktraceKey:    "",
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      zapcore.CapitalColorLevelEncoder,
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

	staticLogger = &logger{
		logLevel:      &ll,
		SugaredLogger: l.Sugar(),
	}
}
