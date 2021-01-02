package logger

import (
	"encoding/json"
	"fmt"
	"github.com/team-zf/framework/config"
	"github.com/team-zf/framework/logger/devices"
)

var (
	logger *devices.Logger
)

func Init(debug bool, logDir string, conf *config.LoggerConfig) {
	logger = devices.NewLogger()

	// 控制台输出
	if debug {
		logger.SetLogger(devices.DeviceConsole)
	}

	// 写文件
	if conf.File != nil {
		settings := conf.File
		Prefix := ""
		if settings.Prefix != "" {
			Prefix = settings.Prefix
		}
		Suffix := ".log"
		if settings.Suffix != "" {
			Suffix = settings.Suffix
		}

		settings.FileName = fmt.Sprintf("%s/%v%s", logDir, Prefix, Suffix)
		bytes, _ := json.Marshal(settings)
		logger.SetLogger(devices.DeviceFile, string(bytes))
	}
}

func Close() {
	getInstance().Close()
}

func getInstance() *devices.Logger {
	if logger == nil {
		logger = devices.NewLogger()
	}
	return logger
}

func Debug(format string, v ...interface{}) {
	getInstance().Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	getInstance().Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	getInstance().Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	getInstance().Error(format, v...)
}

func Alert(format string, v ...interface{}) {
	getInstance().Alert(format, v...)
}

func Critical(format string, v ...interface{}) {
	getInstance().Critical(format, v...)
}

func Notice(format string, v ...interface{}) {
	getInstance().Notice(format, v...)
}
