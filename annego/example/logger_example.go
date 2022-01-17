package main

import (
	"io"

	"goBase/annego/logger"
)

func logPanic() {
	panic("log panic test")
}

func main() {
	logger.Notice("default log level: %d", logger.GetLogLevel())

	m := "log writer example\nline 1\r\nline 2\nline 3"
	logger.LogLines(logger.LOG_INFO, m)

	ctx := logger.NewLogContent(
		logger.Int("int", 123),
		logger.String("str", "log content"),
	)
	ctx.Log(logger.LOG_INFO)

	ctx.Append(
		logger.Int("int", 100),
		logger.ErrorField(io.EOF),
	)
	ctx.LogFormat(logger.LOG_INFO, "log context")

	defer func() {
		if r := recover(); r != nil {
			logger.LogPanic(logger.LOG_ERR, r)
		}
	}()
	logPanic()
}
