package log

import (
	syslog "log"
	"log/slog"

	zlog "github.com/liut/staffio-backend/log"
)

type logger struct {
	slg *slog.Logger
}

// Default 默认实例
var Default Logger

func init() {
	syslog.SetFlags(syslog.Ltime | syslog.Lshortfile)
	Default = &logger{slg: slog.Default()}
}

func SetLogger(logger Logger) {
	if logger != nil {
		Default = logger
		zlog.SetLogger(logger)
	}
}

func GetLogger() Logger {
	return Default
}

func (z *logger) Debugw(msg string, keysAndValues ...interface{}) {
	z.slg.Debug(msg, keysAndValues...)
}

func (z *logger) Infow(msg string, keysAndValues ...interface{}) {
	z.slg.Info(msg, keysAndValues...)
}

func (z *logger) Warnw(msg string, keysAndValues ...interface{}) {
	z.slg.Warn(msg, keysAndValues...)
}

func (z *logger) Errorw(msg string, keysAndValues ...interface{}) {
	z.slg.Error(msg, keysAndValues...)
}

func (z *logger) Fatalw(msg string, keysAndValues ...interface{}) {
	syslog.Fatal(msg, keysAndValues)
}
