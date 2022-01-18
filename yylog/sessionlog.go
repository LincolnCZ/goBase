package yylog

import (
	"context"
	"sync"

	"go.uber.org/zap/zapcore"
)

type logKey struct{}

type logVal struct {
	mut       sync.Mutex
	fields    []zapcore.Field
	fieldsMap map[string]int
}

func newLogVal(fields []zapcore.Field) *logVal {
	slog := &logVal{
		fields:    make([]zapcore.Field, 0),
		fieldsMap: make(map[string]int),
	}
	slog.logAppend(fields)
	return slog
}

func (slog *logVal) logAppend(fields []zapcore.Field) {
	slog.mut.Lock()
	idx := len(slog.fields)
	for _, f := range fields {
		old, exist := slog.fieldsMap[f.Key]
		if exist {
			slog.fields[old] = f
		} else {
			slog.fields = append(slog.fields, f)
			slog.fieldsMap[f.Key] = idx
			idx++
		}
	}
	slog.mut.Unlock()
}

func (slog *logVal) logFlush(fields []zapcore.Field) []zapcore.Field {
	slog.mut.Lock()
	r := make([]zapcore.Field, 0, len(fields)+len(slog.fields))
	r = append(r, slog.fields...)
	for _, f := range fields {
		old, exist := slog.fieldsMap[f.Key]
		if exist {
			r[old] = f
		} else {
			r = append(r, f)
		}
	}
	slog.mut.Unlock()
	return r
}

// LogStart start a session log
func LogStart(ctx context.Context, fields ...zapcore.Field) context.Context {
	slog, ok := ctx.Value(logKey{}).(*logVal)
	if ok {
		if len(fields) > 0 {
			slog.logAppend(fields)
		}
	} else {
		slog = newLogVal(fields)
		ctx = context.WithValue(ctx, logKey{}, slog)
	}
	return ctx
}

// LogAppend append fields to session log
func LogAppend(ctx context.Context, fields ...zapcore.Field) {
	slog, ok := ctx.Value(logKey{}).(*logVal)
	if ok {
		if len(fields) > 0 {
			slog.logAppend(fields)
		}
	}
}

// LogFlush flush a session log, do not append fields to session log
func LogFlush(ctx context.Context, mkey string, fields ...zapcore.Field) {
	slog, ok := ctx.Value(logKey{}).(*logVal)
	logs := fields
	if ok {
		logs = slog.logFlush(fields)
	}

	defaultLog.writelog("info", mkey, logs...)
}
