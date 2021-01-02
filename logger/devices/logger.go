package devices

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

var (
	LevelPrefix = []string{"[M]", "[A]", "[C]", "[E]", "[W]", "[N]", "[I]", "[D]"}
)

type Logger struct {
	lock        sync.Mutex
	level       int
	init        bool
	contentType string
	devices     map[string]IDevice
}

func NewLogger() *Logger {
	log := new(Logger)
	log.level = LevelDebug
	log.devices = make(map[string]IDevice)
	log.setLogger(DeviceConsole)
	return log
}

func (e *Logger) setLogger(key string, configs ...string) error {
	config := append(configs, "{}")[0]
	if _, ok := e.devices[key]; ok {
		return fmt.Errorf("logs: duplicate devicename %q (you have set this device before)", key)
	}

	device, err := NewDevice(key)
	if err != nil {
		return err
	}
	err = device.Init(config)
	if err == nil {
		e.devices[key] = device
	} else {
		fmt.Fprintln(os.Stderr, "Logger.SetLogger: "+err.Error())
	}
	return err
}

func (e *Logger) formatText(level int, msg string, v ...interface{}) (string, error) {
	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v...)
	}
	return fmt.Sprintf(" %s %s", LevelPrefix[level], msg), nil
}

func (e *Logger) writeMsg(level int, msg string, v ...interface{}) error {
	when := time.Now()
	message, err := e.formatText(level, msg, v...)
	if err == nil {
		e.writeToLoggers(when, message, level)
	}
	return err
}

func (e *Logger) writeToLoggers(when time.Time, msg string, level int) {
	for key, device := range e.devices {
		if err := device.WriteMsg(when, msg, level); err != nil {
			fmt.Fprintf(os.Stderr, "unable to WriteMsg to device: %v, error: %v\n", key, err)
		}
	}
}

func (e *Logger) SetLevel(v int) {
	e.level = v
}

func (e *Logger) SetLogger(key string, configs ...string) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if !e.init {
		e.devices = make(map[string]IDevice)
		e.init = true
	}
	return e.setLogger(key, configs...)
}

func (e *Logger) Debug(format string, v ...interface{}) {
	if LevelDebug > e.level {
		return
	}
	e.writeMsg(LevelDebug, format, v...)
}

func (e *Logger) Info(format string, v ...interface{}) {
	if LevelInformational > e.level {
		return
	}
	e.writeMsg(LevelInformational, format, v...)
}

func (e *Logger) Warn(format string, v ...interface{}) {
	if LevelWarning > e.level {
		return
	}
	e.writeMsg(LevelWarning, format, v...)
}

func (e *Logger) Error(format string, v ...interface{}) {
	if LevelError > e.level {
		return
	}
	e.writeMsg(LevelError, format, v...)
}

func (e *Logger) Alert(format string, v ...interface{}) {
	if LevelAlert > e.level {
		return
	}
	e.writeMsg(LevelAlert, format, v...)
}

func (e *Logger) Critical(format string, v ...interface{}) {
	if LevelCritical > e.level {
		return
	}
	e.writeMsg(LevelCritical, format, v...)
}

func (e *Logger) Notice(format string, v ...interface{}) {
	if LevelNotice > e.level {
		return
	}
	e.writeMsg(LevelNotice, format, v...)
}

func (e *Logger) Close() {
	for _, item := range e.devices {
		item.Flush()
		item.Destroy()
	}
	e.devices = nil
}
