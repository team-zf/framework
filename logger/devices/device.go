package devices

import (
	"fmt"
	"io"
	"sync"
	"time"
)

const (
	DeviceConsole   = "console"
	DeviceFile      = "file"
	DeviceMultiFile = "multifile"
)

var (
	devices = make(map[string]func() IDevice)
)

type IDevice interface {
	Init(config string) error
	WriteMsg(when time.Time, msg string, level int) error
	Destroy()
	Flush()
}

type DeviceWriter struct {
	sync.Mutex
	writer io.Writer
}

func (e *DeviceWriter) Println(when time.Time, msg string) {
	e.Lock()
	h, _ := FormatTimeHeader(when)
	e.writer.Write(append(append(h, msg...), '\n'))
	e.Unlock()
}

func FormatTimeHeader(when time.Time) ([]byte, int) {
	_, _, d := when.Date()
	return []byte(when.Format("2006-01-02 15:04:05")), d
}

func Register(key string, _func func() IDevice) {
	if _func == nil {
		panic("devices: Register provide is nil")
	}
	if _, dup := devices[key]; dup {
		panic("devices: Register called twice for provider " + key)
	}
	devices[key] = _func
}

func NewDevice(name string) (IDevice, error) {
	_func, ok := devices[name]
	if !ok {
		return nil, fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", name)
	}
	return _func(), nil
}
