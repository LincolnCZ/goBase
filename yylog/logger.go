package yylog

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"goBase/yylog/yyencode"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type yyloger struct {
	logger *zap.Logger
	Level  zap.AtomicLevel
}

var (
	defaultLog *yyloger = nil
	stdlog     *log.Logger
	once       = &sync.Once{}
)

// 启动后自动初始化日志对象
func init() {
	//初始化日志为stdout
	procname := filepath.Base(os.Args[0])
	InitLog(ProcessName(procname), SetTarget("stdout"))

	// 注入 yyjson 日志格式模块，满足YY业务通用的输出要求
	zap.RegisterEncoder("yyjson", func(encoderConfig zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return yyencode.NewYYEncoder(encoderConfig), nil
	})
}

// InitLog 根据options的设置,初始化日志系统。
// Example :
//     yylog.InitLog(yylog.ProcessName("udbapp_xxx"), yylog.SetTarget("asyncfile"), yylog.SetEncode("yyjson"), yylog.LogFilePath("../log/"))
func InitLog(options ...LogOption) error {
	var (
		err    error
		logger *zap.Logger
		level  zap.AtomicLevel
	)
	config := defaultLogOptions
	for _, option := range options {
		option.apply(&config)
	}

	if level, logger, err = zapLogInit(&config); err != nil {
		fmt.Printf("ZapLogInit err:%v", err)
		return err
	}

	logger = logger.WithOptions(zap.AddCallerSkip(2))
	if defaultLog == nil {
		defaultLog = &yyloger{logger, level}
	} else {
		defaultLog.logger = logger
		defaultLog.Level = level
	}
	// redirect go log to defaultLog.logger
	zap.RedirectStdLog(defaultLog.logger)
	stdlog = zap.NewStdLog(defaultLog.logger)

	return nil
}

// InitYYServerLog yy服务器使用的默认日志配置
func InitYYServerLog() error {
	procname := filepath.Base(os.Args[0])
	return InitLog(SetTarget("asyncfile"), SetEncode("json"),
		ProcessName(procname),
		LogFileName(procname+".gfy"),
		LogFilePath("/data/yy/log/"+procname))
}

// GetLogger 获取日志对象，一般用于注入第三方库
func GetLogger() *yyloger {
	return defaultLog
}

// GetStdLogger 获取标准的 log.Logger 对象，用于注入第三方库
func GetStdLogger() *log.Logger {
	return stdlog
}

// GetZlog 获取底层的zap.Logger 对象（共享输出对象），用于业务二次封装
func (l *yyloger) GetZLog(opts ...zap.Option) *zap.Logger {
	return l.logger.WithOptions(opts...)
}

// Clone 复制yyloger 对象，用于业务二次封装
// Example : GetLogger().Clone(zap.AddCallerSkip(2))
func (l *yyloger) Clone(opts ...zap.Option) *yyloger {
	nl := &yyloger{
		logger: l.logger,
		Level:  l.Level,
	}

	nl.logger = l.logger.WithOptions(opts...)
	return nl
}

// Write  实现 io.Writer
func (l *yyloger) Write(p []byte) (n int, err error) {
	l.writelog("info", string(p))
	return len(p), nil
}

// Log 实现 github.com/go-log/log.logger 接口
func (l *yyloger) Log(v ...interface{}) {
	msg := fmt.Sprint(v...)
	l.writelog("info", msg)
}

//Logf 实现 github.com/go-log/log.logger 接口
func (l *yyloger) Logf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.writelog("info", msg)
}

// Error  输出error级别 的日志
func (l *yyloger) Error(msg string) {
	l.writelog("error", msg)
}

// Infof  输出info 级别日志
func (l *yyloger) Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.writelog("info", msg)
}

func (l *yyloger) writelog(level, msg string, fields ...zapcore.Field) {
	switch level {
	case "info":
		l.logger.Info(msg, fields...)
	case "debug":
		l.logger.Debug(msg, fields...)
	case "warn":
		l.logger.Warn(msg, fields...)
	case "error":
		l.logger.Error(msg, fields...)
	case "panic":
		l.logger.Panic(msg, fields...)
	case "dpanic":
		l.logger.DPanic(msg, fields...)
	case "fatal":
		l.logger.Fatal(msg, fields...)
	default:
		l.logger.Info(msg, fields...)
	}
}

// Log 设置不同日志级别的日志
// level 日志级别: debug info warn error panic fatal
func Log(level string, v ...interface{}) {
	msg := fmt.Sprint(v...)
	defaultLog.writelog(level, msg)
}

// LogF 设置不同日志级别的日志, 支持自定义format
// level 日志级别: debug info warn error panic fatal
func LogF(level string, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)

	defaultLog.writelog(level, msg)
}

// Debug 设置debug 级别的日志 .
//  msg 日志关键描述信息
//  fields 由k-v组成的日志信息集合
func Debug(msg string, fields ...zapcore.Field) {
	defaultLog.writelog("debug", msg, fields...)
}

// Info 设置 info 级别的日志 .
//  msg 日志关键描述信息
//  fields 由k-v组成的日志信息集合
func Info(msg string, fields ...zapcore.Field) {
	defaultLog.writelog("info", msg, fields...)
}

// Warn 设置 warn 级别的日志 .
//  msg 日志关键描述信息
//  fields 由k-v组成的日志信息集合
func Warn(msg string, fields ...zapcore.Field) {
	defaultLog.writelog("warn", msg, fields...)
}

// Error 设置 error 级别的日志 .
//  msg 日志关键描述信息
//  fields 由k-v组成的日志信息集合
func Error(msg string, fields ...zapcore.Field) {
	defaultLog.writelog("error", msg, fields...)
}

// Panic 设置 panic 级别的日志 . 输出日志后，触发panic.
//  msg 日志关键描述信息
//  fields 由k-v组成的日志信息集合
func Panic(msg string, fields ...zapcore.Field) {
	defaultLog.writelog("panic", msg, fields...)
}

// Fatal 设置 fatal级别的日志 . 输出日志后，触发 os.Exit(1).
//  msg 日志关键描述信息
//  fields 由k-v组成的日志信息集合
func Fatal(msg string, fields ...zapcore.Field) {
	defaultLog.writelog("fatal", msg, fields...)
}

// Sync 手工触发日志模块sync
func Sync() error {
	return defaultLog.logger.Sync()
}

// SetLogLevel 设置全局日志模块的 可输出级别
// level : debug(all) info warn error fatal(off,none)
func SetLogLevel(level string) error {
	switch strings.ToLower(level) {
	case "debug", "info", "warn", "error", "fatal":
		level = strings.ToLower(level)
	case "all":
		level = "debug"
	case "off", "none":
		level = "fatal"
	default:
		return errors.New("not support level")
	}
	defaultLog.Level.UnmarshalText([]byte(level))
	return nil
}
