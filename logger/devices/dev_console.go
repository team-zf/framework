package devices

import (
	"encoding/json"
	"os"
	"time"
)

type ConsoleDevice struct {
	writer *DeviceWriter
	Level  int `json:"level"`
}

func (e *ConsoleDevice) Init(jsonConfig string) error {
	if jsonConfig == "" {
		return nil
	}
	return json.Unmarshal([]byte(jsonConfig), e)
}

func (e *ConsoleDevice) WriteMsg(when time.Time, msg string, level int) error {
	if level > e.Level {
		return nil
	}
	e.writer.Println(when, msg)
	return nil
}

func (e *ConsoleDevice) Destroy() {
}

func (e *ConsoleDevice) Flush() {
}

func init() {
	Register(DeviceConsole, func() IDevice {
		return &ConsoleDevice{
			writer: &DeviceWriter{writer: os.Stdout},
			Level:  LevelDebug,
		}
	})
}
