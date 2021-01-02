package devices

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FileDevice struct {
	sync.RWMutex

	// The opened file
	FileName string `json:"filename"`
	file     *os.File
	writer   *bufio.Writer

	// Rotate at size
	MaxSize        int `json:"maxsize"`
	maxSizeCurSize int

	// Rotate at line
	MaxLines         int `json:"maxlines"`
	maxLinesCurLines int

	// Rotate daily
	Daily         bool  `json:"daily"`
	MaxDays       int64 `json:"maxdays"`
	dailyOpenDate int
	dailyOpenTime time.Time
	Rotate        bool   `json:"rotate"`
	Level         int    `json:"level"`
	MinLevel      int    `json:"minlevel"`
	Perm          string `json:"perm"`
	RotatePerm    string `json:"rotateperm"`

	fileNameOnly, suffix string
}

func (e *FileDevice) Init(jsonConfig string) error {
	if jsonConfig == "" {
		return nil
	}
	if err := json.Unmarshal([]byte(jsonConfig), e); err != nil {
		return err
	}
	if e.FileName == "" {
		return errors.New("jsonconfig must have filename")
	}
	e.suffix = filepath.Ext(e.FileName)
	e.fileNameOnly = strings.TrimSuffix(e.FileName, e.suffix)
	if e.suffix == "" {
		e.suffix = ".log"
	}
	return e.startLogger()
}

func (e *FileDevice) WriteMsg(when time.Time, msg string, level int) error {
	if level > e.Level || level < e.MinLevel {
		return nil
	}
	h, d := FormatTimeHeader(when)
	msg = string(h) + msg + "\n"
	if e.Rotate {
		e.RLock()
		if e.needRotate(d) {
			e.RUnlock()
			e.Lock()
			if e.needRotate(d) {
				if err := e.doRotate(when); err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", e.FileName, err)
				}
			}
			e.Unlock()
		} else {
			e.RUnlock()
		}
	}
	e.Lock()
	_, err := e.writer.Write([]byte(msg))
	if err == nil {
		e.maxLinesCurLines++
		e.maxSizeCurSize += len(msg)
	}
	e.Unlock()
	return err
}

func (e *FileDevice) Destroy() {
	e.writer.Flush()
	e.file.Sync()
	e.file.Close()
}

func (e *FileDevice) Flush() {
	e.writer.Flush()
	e.file.Sync()
}

func (e *FileDevice) startLogger() error {
	file, writer, err := e.createLogFile()
	if err != nil {
		return err
	}
	if e.writer != nil {
		e.writer.Flush()
		e.file.Close()
	}
	e.file = file
	e.writer = writer
	return e.initFile()
}

func (e *FileDevice) createLogFile() (*os.File, *bufio.Writer, error) {
	perm, err := strconv.ParseInt(e.Perm, 8, 64)
	if err != nil {
		return nil, nil, err
	}
	file, err := os.OpenFile(e.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err == nil {
		os.Chmod(e.FileName, os.FileMode(perm))
	}
	writer := bufio.NewWriterSize(file, 8192)
	return file, writer, err
}

func (e *FileDevice) initFile() error {
	if err := e.writer.Flush(); err != nil {
		return err
	}
	info, err := e.file.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s", err)
	}
	e.maxSizeCurSize = int(info.Size())
	e.maxLinesCurLines = 0
	e.dailyOpenTime = time.Now()
	e.dailyOpenDate = e.dailyOpenTime.Day()
	if e.Daily {
		go e.dailyRotate(e.dailyOpenTime)
	}
	if info.Size() > 0 && e.MaxLines > 0 {
		count, err := e.lines()
		if err != nil {
			return err
		}
		e.maxLinesCurLines = count
	}
	return nil
}

func (e *FileDevice) dailyRotate(openTime time.Time) {
	y, m, d := openTime.Add(time.Hour * 24).Date()
	nextDay := time.Date(y, m, d, 0, 0, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextDay.UnixNano() - openTime.UnixNano() + 100))
	<-tm.C
	e.Lock()
	if e.needRotate(time.Now().Day()) {
		if err := e.doRotate(time.Now()); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", e.FileName, err)
		}
	}
	e.Unlock()
}

func (e *FileDevice) lines() (int, error) {
	file, err := os.Open(e.FileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buf := make([]byte, 32768) // 32KB
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}
		count += bytes.Count(buf[:c], lineSep)
		if err == io.EOF {
			break
		}
	}
	return count, nil
}

func (e *FileDevice) needRotate(day int) bool {
	v1 := e.MaxLines > 0 && e.maxLinesCurLines >= e.MaxLines
	v2 := e.MaxSize > 0 && e.maxSizeCurSize >= e.MaxSize
	v3 := e.Daily && day != e.dailyOpenDate
	return v1 && v2 && v3
}

func (e *FileDevice) doRotate(logTime time.Time) error {
	num := 1
	fileName := ""
	rotatePerm, err := strconv.ParseInt(e.RotatePerm, 8, 64)
	if err != nil {
		return err
	}

	if _, err := os.Lstat(e.FileName); err != nil {
		goto RESTART_LOGGER
	}

	if e.MaxLines > 0 || e.MaxSize > 0 {
		for ; err == nil && num <= 999; num++ {
			fileName = e.fileNameOnly + fmt.Sprintf(".%s.%03d%s", logTime.Format("2006-01-02"), num, e.suffix)
			_, err = os.Lstat(fileName)
		}
	} else {
		fileName = fmt.Sprintf("%s.%s%s", e.fileNameOnly, e.dailyOpenTime.Format("2006-01-02"), e.suffix)
		_, err = os.Lstat(fileName)
		for ; err == nil && num <= 999; num++ {
			fileName = e.fileNameOnly + fmt.Sprintf(".%s.%03d%s", e.dailyOpenTime.Format("2006-01-02"), num, e.suffix)
			_, err = os.Lstat(fileName)
		}
	}
	if err == nil {
		return fmt.Errorf("Rotate: Cannot find free log number to rename %s", e.FileName)
	}

	e.writer.Flush()
	e.file.Close()

	err = os.Rename(e.FileName, fileName)
	if err != nil {
		goto RESTART_LOGGER
	}

	err = os.Chmod(fileName, os.FileMode(rotatePerm))

RESTART_LOGGER:

	startLoggerErr := e.startLogger()
	go e.deleteOldLog()

	if startLoggerErr != nil {
		return fmt.Errorf("Rotate StartLogger: %s", startLoggerErr)
	}
	if err != nil {
		return fmt.Errorf("Rotate: %s", err)
	}
	return nil
}

func (e *FileDevice) deleteOldLog() {
	dir := filepath.Dir(e.FileName)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Unable to delete old log '%s', error: %v\n", path, r)
			}
		}()
		return nil
	})
}

func init() {
	Register(DeviceFile, func() IDevice {
		return &FileDevice{
			Daily:      true,
			MaxDays:    365,
			Rotate:     true,
			RotatePerm: "0440",
			Level:      LevelDebug,
			MinLevel:   LevelEmergency,
			Perm:       "0660",
		}
	})
}
