package mylogger

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

//自定义一个日志库

type LogLevel uint16

const (
	_ LogLevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

type ConsoleLogger struct {
	Level LogLevel
}

func NewLog(levelStr string) ConsoleLogger {
	level := parseLogLevel(levelStr)
	return ConsoleLogger{Level: level}
}
func parseLogLevel(s string) LogLevel {
	switch strings.ToLower(s) {
	case "debug":
		return DEBUG
	case "trace":
		return TRACE
	case "info":
		return INFO
	case "warning":
		return WARNING
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "ONFO"
	}
}

func (c ConsoleLogger) enable(msg string) bool {
	return c.Level <= parseLogLevel(msg)
}

func (c ConsoleLogger) logPrint(lv LogLevel, format string, arg ...interface{}) {
	if c.enable(lv.String()) {
		funcName, fileName, lineNo := getInfo(3)
		msg := fmt.Sprintf(format, arg...)
		fmt.Printf("[%s] [%s] [%s:%s:line:%d ]: %s \n", time.Now().Format("2006-01-02 15:04:05"), lv, fileName, funcName, lineNo, msg)
	}
}
func (c ConsoleLogger) Debug(format string, arg ...interface{}) {

		c.logPrint(DEBUG, format, arg...)

}

func (c ConsoleLogger) Trace(format string, arg ...interface{}) {

		c.logPrint(TRACE, format, arg...)

}

func (c ConsoleLogger) Info(format string, arg ...interface{}) {

		c.logPrint(INFO, format, arg...)

}

func (c ConsoleLogger) Warning(format string, arg ...interface{}) {
		c.logPrint(WARNING, format, arg...)
}

func (c ConsoleLogger) Error(format string, arg ...interface{}) {
		c.logPrint(ERROR, format, arg...)
}

func (c ConsoleLogger) Fatal(format string, arg ...interface{}) {
		c.logPrint(FATAL, format, arg...)
}

func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, fileName, lineNo, ok := runtime.Caller(skip)

	if !ok {
		fmt.Printf("runtime.caller() failed/n")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	funcName = strings.Split(funcName, ".")[1]
	fileName = path.Base(fileName)

	return
}
