package logger

import (
	"fmt"
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

func New() Logger {

	return &logger{
		headerName: "Panda",
		showDate:   true,
		showTime:   true,
		showLogs:   true,
	}
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

func (l *logger) log(lType LogType, message string) {
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

func (l *logger) Info(msg string, kv ...interface{}) {
	l.log(Info, msg)
}

func (l *logger) Warn(msg string, kv ...interface{}) {
	l.log(Warning, msg)
}

func (l *logger) Error(msg string, kv ...interface{}) {
	l.log(Error, msg)
}

func (l *logger) Sync() {

}
