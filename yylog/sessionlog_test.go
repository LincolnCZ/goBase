package yylog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func checkLogVal(slog *logVal) bool {
	for idx, field := range slog.fields {
		v, ok := slog.fieldsMap[field.Key]
		if !ok {
			return false
		}
		if v != idx {
			return false
		}
	}
	return true
}

func TestLogAppend(t *testing.T) {
	slog := newLogVal([]zap.Field{zap.String("key1", "val1"), zap.String("key2", "val2")})
	assert.True(t, checkLogVal(slog))
	assert.True(t, slog.fields[0].Equals(zap.String("key1", "val1")))
	assert.True(t, slog.fields[1].Equals(zap.String("key2", "val2")))

	slog.logAppend([]zap.Field{zap.Int("key1", 1), zap.String("key3", "val3")})
	assert.True(t, checkLogVal(slog))
	assert.True(t, slog.fields[0].Equals(zap.Int("key1", 1)))
	assert.True(t, slog.fields[1].Equals(zap.String("key2", "val2")))
	assert.True(t, slog.fields[2].Equals(zap.String("key3", "val3")))
}

func TestLogFlush(t *testing.T) {
	slog := newLogVal([]zap.Field{zap.String("key1", "val1"), zap.String("key2", "val2")})
	assert.True(t, checkLogVal(slog))

	fields := slog.logFlush([]zap.Field{zap.Int("key1", 1), zap.String("key3", "val3")})
	assert.True(t, checkLogVal(slog))
	assert.True(t, slog.fields[0].Equals(zap.String("key1", "val1")))
	assert.True(t, slog.fields[1].Equals(zap.String("key2", "val2")))

	assert.True(t, fields[0].Equals(zap.Int("key1", 1)))
	assert.True(t, fields[1].Equals(zap.String("key2", "val2")))
	assert.True(t, fields[2].Equals(zap.String("key3", "val3")))
}
