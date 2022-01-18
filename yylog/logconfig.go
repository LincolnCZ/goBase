package yylog

import (
	"os"
	"path/filepath"
)

const logpath = "/log"

type logConf struct {
	processName   string
	withPid       bool   // default true
	encodeing     string // json console [todo add more]
	targetName    string // stdout asyncfile
	logFileName   string // asyncfile: default processName
	logFilePath   string // asyncfile: default "../log"
	logFileRotate string // asyncfile: default no option [date,hour]
	HostName      string
}

var defaultLogOptions logConf = logConf{
	processName: filepath.Base(os.Args[0]),
	withPid:     true,
	encodeing:   "json",
	targetName:  "stdout",
	logFileName: filepath.Base(os.Args[0]),
	logFilePath: "../log",
}

// LogOption 日志选项接口
type LogOption interface {
	apply(*logConf)
}

type logOptionFunc func(*logConf)

func (lf logOptionFunc) apply(option *logConf) {
	lf(option)
}

// SetTarget 设置输出类型 stdout asyncfile
func SetTarget(name string) LogOption {
	return logOptionFunc(func(option *logConf) {
		option.targetName = name
	})
}

// LogFilePath 设置输出日志名 name.log，当 target = asyncfile 时需要此选项
func LogFileName(name string) LogOption {
	return logOptionFunc(func(option *logConf) {
		option.logFileName = name
	})
}

// LogFilePath 设置输出日志路径，当 target = asyncfile 时需要此选项
func LogFilePath(path string) LogOption {
	return logOptionFunc(func(option *logConf) {
		option.logFilePath = path
	})
}

// LogFileRotate 设置日志文件归档周期，仅支持 hour date.
func LogFileRotate(r string) LogOption {
	return logOptionFunc(func(option *logConf) {
		if r == "hour" || r == "date" {
			option.logFileRotate = r
		}
	})
}

// SetEncode 设置日志编码格式：yyjson json console.
// 默认为json
func SetEncode(enc string) LogOption {
	return logOptionFunc(func(option *logConf) {
		if enc == "json" || enc == "console" || enc == "yyjson" {
			//fmt.Println("set encode",enc)
			option.encodeing = enc
		} else {
			option.encodeing = "json"
		}
	})
}

// WithPid 设置日志输出中是否加入pid的项。
// 默认为true。
func WithPid(yes bool) LogOption {
	return logOptionFunc(func(option *logConf) {
		option.withPid = yes
	})
}

// ProcessName 设置输出的进程名字。
// 默认去当前执行文件的名字。
func ProcessName(pname string) LogOption {
	return logOptionFunc(func(option *logConf) {
		option.processName = pname
	})
}

// HostName 设置日志机器的ip地址,方便定位。
// 默认不输出。
func HostName(hostname string) LogOption {
	return logOptionFunc(func(option *logConf) {
		option.HostName = hostname
	})
}
