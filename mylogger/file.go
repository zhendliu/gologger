package mylogger

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"
)

//向文件中写日志

type FileLogger struct {
	Level         LogLevel
	filePath      string
	fileName      string
	maxFileSize   int64
	fileHandle    *os.File
	errFileHandle *os.File
	logChan       chan *logMsg
}

type logMsg struct {
	level     LogLevel
	msg       string
	funcName  string
	fileName  string
	timestamp string
	line      int
}

//
func NewFileLogger(levelStr, fp, fn string, maxSize int64) *FileLogger {

	f := &FileLogger{
		Level:       parseLogLevel(levelStr),
		filePath:    fp,
		fileName:    fn,
		maxFileSize: maxSize,
		logChan:     make(chan *logMsg, 50000),
	}
	err := f.initFile()
	if err != nil {
		panic(err)
	}
	go f.writeLogBackground()
	return f
}
func (f *FileLogger) initFile() error {
	fullName := path.Join(f.filePath, f.fileName)
	var err error
	f.fileHandle, err = os.OpenFile(fullName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("init log file error:%s", err)
		return err
	}
	f.errFileHandle, err = os.OpenFile(fullName+"err", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("init log file error:%s", err)
		return err
	}
	return nil
}

//func (f *FileLogger) Close() {
//	f.fileHandle.Close()
//	f.errFileHandle.Close()
//}

func (f *FileLogger) enable(msg string) bool {
	return f.Level <= parseLogLevel(msg)
}
func (f *FileLogger) checkSize(file *os.File) bool {

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("err:", err)
		panic(err)
		return false
	}
	return fileInfo.Size() >= f.maxFileSize
}

func (f *FileLogger) splitFile(file *os.File) (*os.File, error) {

	//cut log
	rand.Seed(time.Now().UnixNano())

	nowStr := time.Now().Format("20060102150405")
	nowStr =fmt.Sprintf("%s-%d",nowStr , rand.Intn(1000))

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	logName := path.Join(f.filePath, fileInfo.Name())
	newLogName := fmt.Sprintf("%s.bak%s", logName, nowStr)
	// closed file
	file.Close()
	// rename file
	os.Rename(logName, newLogName)
	// open new file
	fileObj, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("init log file error:%s", err)
		return nil, err
	}
	return fileObj, nil
}

//后台写入日志

func (f *FileLogger) writeLogBackground() {

	for {
		if f.checkSize(f.fileHandle) {
			newFile, err := f.splitFile(f.fileHandle)
			if err != nil {
				panic(err)
				return
			}
			f.fileHandle = *&newFile
		}
		select {
		case logTmp := <-f.logChan:
			fmt.Println("取出日志")
			fmt.Println(logTmp)
			logInfo := fmt.Sprintf("[%s] [%s] [%s:%s:line:%d ]: %s \n", logTmp.timestamp, logTmp.level, logTmp.fileName, logTmp.funcName, logTmp.line, logTmp.msg)
			fmt.Fprintf(f.fileHandle, logInfo)
			if logTmp.level >= ERROR {
				if f.checkSize(f.errFileHandle) {
					newFile, err := f.splitFile(f.errFileHandle)
					if err != nil {
						panic(err)
						return
					}
					f.errFileHandle = newFile
				}
				//记录日志级别大于error级别，写入err文件中
				fmt.Fprintf(f.errFileHandle, logInfo)
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}
func (f *FileLogger) logPrint(lv LogLevel, format string, arg ...interface{}) {

	if f.enable(lv.String()) {
		funcName, fileName, lineNo := getInfo(3)
		msg := fmt.Sprintf(format, arg...)
		//日志发送到通道中
		logTmp := &logMsg{
			level:     lv,
			msg:       msg,
			funcName:  funcName,
			fileName:  fileName,
			timestamp: time.Now().Format("2006-01-02 15:04:05"),
			line:      lineNo,
		}
		select {
		case f.logChan <- logTmp:
		default:
			//丢弃溢出的阻塞
		}
	}

}

func (f *FileLogger) Debug(format string, arg ...interface{}) {
	if f.enable("DEBUG") {
		f.logPrint(DEBUG, format, arg...)
	}
}

func (f *FileLogger) Trace(format string, arg ...interface{}) {
	if f.enable("TRACE") {
		f.logPrint(TRACE, format, arg...)
	}
}

func (f *FileLogger) Info(format string, arg ...interface{}) {
	if f.enable("INFO") {
		f.logPrint(INFO, format, arg...)
	}
}

func (f *FileLogger) Warning(format string, arg ...interface{}) {
	if f.enable("WARNING") {
		f.logPrint(WARNING, format, arg...)
	}
}

func (f *FileLogger) Error(format string, arg ...interface{}) {
	if f.enable("ERROR") {
		f.logPrint(ERROR, format, arg...)
	}
}

func (f *FileLogger) Fatal(format string, arg ...interface{}) {
	if f.enable("FATAL") {
		f.logPrint(FATAL, format, arg...)
	}
}
