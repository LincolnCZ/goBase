package logger

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogContext(t *testing.T) {
	type TestObject struct {
		A int
		S string
	}
	testObject := &TestObject{
		A: 100,
		S: "test",
	}

	ctx := NewLogContent(
		Int("int", 123),
		Float32("float", 10.0),
		String("str", "hello"),
		Object("object", testObject),
	)
	s := ctx.String()

	obj := struct {
		Int    int
		Float  float64
		Str    string
		Object *TestObject
	}{}
	assert.NoError(t, json.Unmarshal([]byte(s), &obj))
	assert.Equal(t, 123, obj.Int)
	assert.Equal(t, 10.0, obj.Float)
	assert.Equal(t, "hello", obj.Str)
	assert.Equal(t, testObject, obj.Object)
}

func TestCopy(t *testing.T) {
	ctx1 := NewLogContent(
		Int("int", 123),
		Float32("float", 10.0),
		String("str", "hello"),
	)
	ctx2 := ctx1.Copy()

	assert.Equal(t, ctx1.fields[0].String(), ctx2.fields[0].String())
	assert.Equal(t, ctx1.fields[1].String(), ctx2.fields[1].String())
	assert.Equal(t, ctx1.fields[2].String(), ctx2.fields[2].String())

	ctx2.Append(Int("int", 100))
	assert.Equal(t, "\"int\":123", ctx1.fields[0].String())
	assert.Equal(t, "\"int\":100", ctx2.fields[0].String())
}

func TestSkip(t *testing.T) {
	ctx := NewLogContent(
		String("a", "a"),
		Int("b", 1),
		Skip(),
	)
	s := ctx.String()

	v := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal([]byte(s), &v))
	assert.Equal(t, 2, len(v))
}
