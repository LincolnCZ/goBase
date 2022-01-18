package asyncLog

import (
	"fmt"
	"math/rand"
)

// Priority 日志优先级
type Priority int

const (
	// LevelAll 所有日志
	LevelAll Priority = iota
	// LevelDebug  debug日志
	LevelDebug
	// LevelInfo  info日志
	LevelInfo
	// LevelWarn  warn 日志
	LevelWarn
	// LevelError  error 日志
	LevelError
	// LevelFatal  fatal 日志
	LevelFatal
	// LevelOff  关闭日志
	LevelOff
)

var (
	// 日志等级
	levelTitle = map[Priority]string{
		LevelDebug: "[DEBUG]",
		LevelInfo:  "[INFO]",
		LevelWarn:  "[WARN]",
		LevelError: "[ERROR]",
		LevelFatal: "[FATAL]",
	}
)

// NewLevelLog 写入等级日志
// 级别高于logLevel才会被写入
func NewLevelLog(filename string, logLevel Priority) *LogFile {
	lf := NewLogFile(filename)
	lf.level = logLevel

	return lf
}

// SetLevel 设置日志等级
func (lf *LogFile) SetLevel(logLevel Priority) {
	lf.level = logLevel
}

// Debug  输出debug 日志
func (lf *LogFile) Debug(format string, a ...interface{}) error {
	return lf.writeLevelMsg(LevelDebug, format, a...)
}

// Info  输出info 日志
func (lf *LogFile) Info(format string, a ...interface{}) error {
	return lf.writeLevelMsg(LevelInfo, format, a...)
}

// Warn  输出warn 日志
func (lf *LogFile) Warn(format string, a ...interface{}) error {
	return lf.writeLevelMsg(LevelWarn, format, a...)
}

// Error  输出error 日志
func (lf *LogFile) Error(format string, a ...interface{}) error {
	return lf.writeLevelMsg(LevelError, format, a...)
}

// Fatal  输出fatal 日志
func (lf *LogFile) Fatal(format string, a ...interface{}) error {
	return lf.writeLevelMsg(LevelFatal, format, a...)
}

func (lf *LogFile) writeLevelMsg(level Priority, format string, a ...interface{}) error {
	if lf.probability < 1.0 && rand.Float32() > lf.probability {
		// 按照概率写入
		return nil
	}

	if level >= lf.level {
		msg := fmt.Sprintf(format, a...)
		return lf.Write(levelTitle[level] + " " + msg)
	}

	return nil
}
