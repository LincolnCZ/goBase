package logger

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type LogContent struct {
	mut       sync.Mutex
	fields    []Field
	fieldsMap map[string]int
}

func NewLogContent(fields ...Field) *LogContent {
	l := &LogContent{
		fields:    make([]Field, 0),
		fieldsMap: make(map[string]int),
	}
	l.Append(fields...)
	return l
}

func (slog *LogContent) Append(fields ...Field) {
	slog.mut.Lock()
	idx := len(slog.fields)
	for _, f := range fields {
		if f.ftype == skipType {
			continue
		}
		old, exist := slog.fieldsMap[f.key]
		if exist {
			slog.fields[old] = f
		} else {
			slog.fields = append(slog.fields, f)
			slog.fieldsMap[f.key] = idx
			idx++
		}
	}
	slog.mut.Unlock()
}

func (slog *LogContent) Copy() *LogContent {
	slog.mut.Lock()
	l := len(slog.fields)
	logNew := &LogContent{
		fields:    make([]Field, l),
		fieldsMap: make(map[string]int, l),
	}
	for i, f := range slog.fields {
		logNew.fields[i] = f
	}
	for k, v := range slog.fieldsMap {
		logNew.fieldsMap[k] = v
	}
	slog.mut.Unlock()
	return logNew
}

func (slog *LogContent) String() string {
	slog.mut.Lock()
	var oss strings.Builder
	oss.WriteByte('{')
	for idx, f := range slog.fields {
		if idx > 0 {
			oss.WriteByte(',')
		}
		oss.WriteString(f.String())
	}
	oss.WriteByte('}')
	slog.mut.Unlock()
	return oss.String()
}

func (slog *LogContent) Log(level int) {
	if level > logLevel {
		return
	}
	logf := GetLoggerFunc(level)
	logf(slog.String())
}

func (slog *LogContent) LogFormat(level int, format string, args ...interface{}) {
	if level > logLevel {
		return
	}
	logf := GetLoggerFunc(level)
	l := fmt.Sprintf(format, args...)
	logf(slog.String() + " " + l)
}

type loggerKey struct{}

func FromContext(ctx context.Context) *LogContent {
	val := ctx.Value(loggerKey{})
	if val == nil {
		return nil
	}
	return val.(*LogContent)
}

func ToContext(ctx context.Context, slog *LogContent) context.Context {
	return context.WithValue(ctx, loggerKey{}, slog)
}
