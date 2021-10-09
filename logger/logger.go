package logger

import (
	"fmt"
	"sync"
	"time"
)

type LogType string

const (
	Info    LogType = "\033[1;34m%s\033[0m"
	Notice  LogType = "\033[1;36m%s\033[0m"
	Warning LogType = "\033[1;33m%s\033[0m"
	Error   LogType = "\033[1;31m%s\033[0m"
	Debug   LogType = "\033[0;36m%s\033[0m"
)

const greenColor = "\033[1;32m%s\033[0m"

type logger struct {
	headerName string
	showDate   bool
	showTime   bool
	showLogs   bool
}

type Logger interface {
	Log(LogType, string)
}

var lock = &sync.Mutex{}
var l *logger

func GetLogger() *logger {
	if l == nil {
		lock.Lock()
		defer lock.Unlock()
		if l == nil {
			l = &logger{
				headerName: "Panda",
				showDate:   true,
				showTime:   true,
				showLogs:   true,
			}
		}
	}
	return l
}

func (l *logger) SetShowLogs(doPrint bool) {
	l.showLogs = doPrint
}

func (l *logger) SetName(name string) {
	l.headerName = name
}

func (l *logger) SetShowDate(show bool) {
	l.showDate = show
}

func (l *logger) SetShowTime(show bool) {
	l.showTime = show
}

func (l *logger) Log(lType LogType, message string) {
	if !l.showLogs {
		return
	}
	var dt string = ""
	var header string = ""
	t := time.Now()
	if l.showDate {
		dt += fmt.Sprintf(greenColor, t.Format("2006-01-02"))
		if l.showTime {
			dt += " "
		} else {
			dt += ":"
		}
	}

	if l.showTime {
		dt += fmt.Sprintf(greenColor, t.Format("15:04:05"))
		dt += ":"
	}

	header = fmt.Sprintf(string(lType), "["+l.headerName+"]")

	fmt.Printf("%s %s %s\n", dt, header, message)
}
