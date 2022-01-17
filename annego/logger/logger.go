package logger

import (
	"bytes"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

const (
	LOG_EMERG   int = int(syslog.LOG_EMERG)
	LOG_ERR     int = int(syslog.LOG_ERR)
	LOG_WARNING int = int(syslog.LOG_WARNING)
	LOG_NOTICE  int = int(syslog.LOG_NOTICE)
	LOG_INFO    int = int(syslog.LOG_INFO)
	LOG_DEBUG   int = int(syslog.LOG_DEBUG)
)

var logLevel int

// Ilogger log接口
type Ilogger interface {
	Debug(str string) error
	Info(str string) error
	Notice(str string) error
	Warning(str string) error
	Err(str string) error
	Emerg(str string) error
}

// GetProgramName 获取程序执行名称
func GetProgramName() string {
	_, file := path.Split(os.Args[0])
	return file
}

// GSyslogName syslog 写日志名称，需要满足下列结尾格式才能写入 /data/yy/log
// _d 带有时间戳，进程名等前置信息
// _rd 没有任何前置信息
var syslogName string

func init() {
	syslogName = GetProgramName()
	logLevel = LOG_INFO
}

type defaultLogger struct{}

func (l defaultLogger) Log(level string, str string) error {
	tm := time.Now()
	fmt.Println(tm.Format("2006-01-02 15:04:05"), level, syslogName+":", str)
	return nil
}

func (l defaultLogger) Debug(str string) error {
	return l.Log("debug", str)
}

func (l defaultLogger) Info(str string) error {
	return l.Log("info", str)
}

func (l defaultLogger) Notice(str string) error {
	return l.Log("notice", str)
}

func (l defaultLogger) Warning(str string) error {
	return l.Log("warning", str)
}

func (l defaultLogger) Err(str string) error {
	return l.Log("err", str)
}

func (l defaultLogger) Emerg(str string) error {
	return l.Log("emerg", str)
}

// 默认打印到标准输出
var logger Ilogger = defaultLogger{}
var initOnce sync.Once

// InitLog 初始化日志到syslog
// 初始化之前的日志会打印到标准输出
func InitLog() error {
	var err error
	initOnce.Do(func() {
		var syslogger Ilogger
		syslogger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_USER, syslogName)
		if err != nil {
			return
		}
		logger = syslogger
	})
	return err
}

func SetLogLevel(level int) {
	logLevel = level
}

func GetLogLevel() int {
	return logLevel
}

func GetLogger() Ilogger {
	return logger
}

// Debug logger
func Debug(format string, args ...interface{}) {
	if LOG_DEBUG <= logLevel {
		_, fn, line, ok := runtime.Caller(1)
		var str string
		if ok {
			_, fn = path.Split(fn)
			str = fmt.Sprintf("[%s:%d] ", fn, line)
		}
		logger.Debug(str + fmt.Sprintf(format, args...))
	}
}

// Info logger
func Info(format string, args ...interface{}) {
	if LOG_INFO <= logLevel {
		_, fn, line, ok := runtime.Caller(1)
		var str string
		if ok {
			_, fn = path.Split(fn)
			str = fmt.Sprintf("[%s:%d] ", fn, line)
		}
		logger.Info(str + fmt.Sprintf(format, args...))
	}
}

// Notice logger
func Notice(format string, args ...interface{}) {
	if LOG_NOTICE <= logLevel {
		_, fn, line, ok := runtime.Caller(1)
		var str string
		if ok {
			_, fn = path.Split(fn)
			str = fmt.Sprintf("[%s:%d] ", fn, line)
		}
		logger.Notice(str + fmt.Sprintf(format, args...))
	}
}

// Warning logger
func Warning(format string, args ...interface{}) {
	if LOG_WARNING <= logLevel {
		_, fn, line, ok := runtime.Caller(1)
		var str string
		if ok {
			_, fn = path.Split(fn)
			str = fmt.Sprintf("[%s:%d] ", fn, line)
		}
		logger.Warning(str + fmt.Sprintf(format, args...))
	}
}

// Error logger
func Error(format string, args ...interface{}) {
	if LOG_ERR <= logLevel {
		_, fn, line, ok := runtime.Caller(1)
		var str string
		if ok {
			_, fn = path.Split(fn)
			str = fmt.Sprintf("[%s:%d] ", fn, line)
		}
		logger.Err(str + fmt.Sprintf(format, args...))
	}
}

// Emerg logger
func Emerg(format string, args ...interface{}) {
	if LOG_EMERG <= logLevel {
		_, fn, line, ok := runtime.Caller(1)
		var str string
		if ok {
			_, fn = path.Split(fn)
			str = fmt.Sprintf("[%s:%d] ", fn, line)
		}
		logger.Emerg(str + fmt.Sprintf(format, args...))
	}
}

type LogHandle func(string) error

// GetLoggerFunc 获取level日志级别的输出函数
func GetLoggerFunc(level int) LogHandle {
	var handle LogHandle
	switch level {
	case LOG_DEBUG:
		handle = logger.Debug
	case LOG_INFO:
		handle = logger.Info
	case LOG_NOTICE:
		handle = logger.Notice
	case LOG_WARNING:
		handle = logger.Warning
	case LOG_ERR:
		handle = logger.Err
	case LOG_EMERG:
		handle = logger.Emerg
	default:
		panic("GetLoggerFunc error level")
	}
	return handle
}

// 实现io.Writer接口，用来替换标准库中log行为
type LogWriter struct {
	level  int
	prefix string
	handle LogHandle
}

func NewLogWrite(level int, prefix string) *LogWriter {
	writer := &LogWriter{
		level:  level,
		prefix: prefix,
		handle: GetLoggerFunc(level),
	}
	return writer
}

func (l *LogWriter) Write(p []byte) (int, error) {
	if l.level > logLevel {
		return len(p), nil
	}

	buff := bytes.NewBuffer(p)
	for {
		line, err := buff.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimRight(line, "\n\r")
		if len(l.prefix) > 0 {
			l.handle(fmt.Sprintf("[%s] %s", l.prefix, line))
		} else {
			l.handle(fmt.Sprintf("%s", line))
		}
	}
	return len(p), nil
}

var stdWriter *LogWriter = nil

// InitRedirectStdlog 替换标准库log行为，使用level日志等级输出
func InitRedirectStdlog(level int) {
	if stdWriter != nil {
		panic("InitStdlogAsSyslog is setting")
	}
	stdWriter := NewLogWrite(level, "stdlog")

	log.SetFlags(log.Lshortfile)
	log.SetOutput(stdWriter)
}

// LogPanic 打印panic信息和堆栈
func LogPanic(level int, r interface{}) {
	if r == nil {
		return
	}
	writer := NewLogWrite(level, "panic")
	s := fmt.Sprintln("recover message:", r)
	m := append([]byte(s), debug.Stack()...)
	writer.Write(m)
}

// LogLines 打印多行日志，根据换行符进行换行
func LogLines(level int, s string) {
	handle := GetLoggerFunc(level)

	l := s
	for len(l) > 0 {
		idx := strings.IndexAny(l, "\r\n")
		if idx == -1 {
			handle(l)
			break
		} else {
			if idx > 0 {
				handle(l[:idx])
			}
			l = l[idx+1:]
		}
	}
}
