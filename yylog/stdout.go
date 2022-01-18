package yylog

import (
	"go.uber.org/zap"
)

type stdZapLogInit struct {
}

func (s *stdZapLogInit) loginit(config *logConf) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		zapconfig zap.Config
		llevel    zap.AtomicLevel
		lzaplog   *zap.Logger
		err       error
	)
	zapconfig = zap.NewProductionConfig()
	zapconfig.OutputPaths = []string{"stdout"}
	zapconfig.ErrorOutputPaths = []string{"stderr"}
	zapconfig.Encoding = config.encodeing
	zapconfig.DisableStacktrace = true
	zapconfig.EncoderConfig.TimeKey = "timestamp"
	zapconfig.EncoderConfig.EncodeTime = epochFullTimeEncoder
	lzaplog, err = zapconfig.Build()
	llevel = zapconfig.Level
	return llevel, lzaplog, err
}
