package yylog

import (
	"fmt"
	"os"
	"time"

	"goBase/yylog/yyencode"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogIniter interface {
	loginit(conf *logConf) (zap.AtomicLevel, *zap.Logger, error)
}

func zapLogInit(config *logConf) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		zapinit zapLogIniter
		level   zap.AtomicLevel
		llog    *zap.Logger
		err     error
	)

	if config.targetName == "asyncfile" {
		zapinit = &asyncFileZapLogInit{}
	} else {
		zapinit = &stdZapLogInit{}
	}

	if level, llog, err = zapinit.loginit(config); err != nil {
		fmt.Printf("loginit err:%v", err)
		return level, llog, err
	}

	if config.encodeing != "yyjson" {
		if config.withPid {
			llog = llog.With(zap.Int("pid", os.Getpid()))
		}
		if config.processName != "" {
			llog = llog.With(zap.String("procname", config.processName))
		}
	} else {
		if config.processName != "" {
			yyencode.ProcessName = config.processName
		}
	}

	if config.HostName != "" {
		llog = llog.With(zap.String("hostname", config.HostName))
	}

	return level, llog, nil
}

func epochFullTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
