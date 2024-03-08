package i18n

import (
	zlog "github.com/liut/staffio/pkg/log"
)

// nolint
func logger() zlog.Logger {
	return zlog.GetLogger()
}
