package yyencode

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
	"math"
	"testing"
	"time"
)

// Nested Array- and ObjectMarshalers.
type turducken struct{}

func (t turducken) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return enc.AddArray("ducks", zapcore.ArrayMarshalerFunc(func(arr zapcore.ArrayEncoder) error {
		for i := 0; i < 2; i++ {
			arr.AppendObject(zapcore.ObjectMarshalerFunc(func(inner zapcore.ObjectEncoder) error {
				inner.AddString("in", "chicken")
				return nil
			}))
		}
		return nil
	}))
}

type turduckens int

func (t turduckens) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	var err error
	tur := turducken{}
	for i := 0; i < int(t); i++ {
		err = multierr.Append(err, enc.AppendObject(tur))
	}
	return err
}

type loggable struct{ bool }

func (l loggable) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if !l.bool {
		return errors.New("can't marshal")
	}
	enc.AddString("loggable", "yes")
	return nil
}

func (l loggable) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	if !l.bool {
		return errors.New("can't marshal")
	}
	enc.AppendBool(true)
	return nil
}

type noJSON struct{}

func (nj noJSON) MarshalJSON() ([]byte, error) {
	return nil, errors.New("no")
}

func assertJSON(t *testing.T, expected string, enc *yyEncoder) {
	assert.Equal(t, expected, enc.buf.String(), "Encoded JSON didn't match expectations.")
}

func assertOutput(t testing.TB, expected string, f func(zapcore.Encoder)) {
	enc := &yyEncoder{buf: Get(), EncoderConfig: &zapcore.EncoderConfig{
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}}
	f(enc)
	assert.Equal(t, expected, enc.buf.String(), "Unexpected encoder output after adding.")

	enc.truncate()
	enc.AddString("foo", "bar")
	f(enc)
	expectedPrefix := `"foo":"bar"`
	if expected != "" {
		// If we expect output, it should be comma-separated from the previous
		// field.
		expectedPrefix += ","
	}
	assert.Equal(t, expectedPrefix+expected, enc.buf.String(), "Unexpected encoder output after adding as a second field.")
}

func TestYYJsonClone(t *testing.T) {
	// The parent encoder is created with plenty of excess capacity.
	parent := &yyEncoder{buf: Get()}
	clone := parent.Clone()

	// Adding to the parent shouldn't affect the clone, and vice versa.
	parent.AddString("foo", "bar")
	clone.AddString("baz", "bing")

	assertJSON(t, `"foo":"bar"`, parent)
	assertJSON(t, `"baz":"bing"`, clone.(*yyEncoder))
}

func TestEscaping(t *testing.T) {
	enc := &yyEncoder{buf: Get()}
	// Test all the edge cases of JSON escaping directly.
	cases := map[string]string{
		// ASCII.
		`foo`: `foo`,
		// Special-cased characters.
		`"`: `\"`,
		`\`: `\\`,
		// Special-cased characters within everyday ASCII.
		`foo"foo`: `foo\"foo`,
		"foo\n":   `foo\n`,
		// Special-cased control characters.
		"\n": `\n`,
		"\r": `\r`,
		"\t": `\t`,
		// \b and \f are sometimes backslash-escaped, but this representation is also
		// conformant.
		"\b": `\u0008`,
		"\f": `\u000c`,
		// The standard lib special-cases angle brackets and ampersands by default,
		// because it wants to protect users from browser exploits. In a logging
		// context, we shouldn't special-case these characters.
		"<": "<",
		">": ">",
		"&": "&",
		// ASCII bell - not special-cased.
		string(byte(0x07)): `\u0007`,
		// Astral-plane unicode.
		`☃`: `☃`,
		// Decodes to (RuneError, 1)
		"\xed\xa0\x80":    `\ufffd\ufffd\ufffd`,
		"foo\xed\xa0\x80": `foo\ufffd\ufffd\ufffd`,
	}

	t.Run("String", func(t *testing.T) {
		for input, output := range cases {
			enc.truncate()
			enc.safeAddString(input)
			assertJSON(t, output, enc)
		}
	})

	t.Run("ByteString", func(t *testing.T) {
		for input, output := range cases {
			enc.truncate()
			enc.safeAddByteString([]byte(input))
			assertJSON(t, output, enc)
		}
	})
}

func TestJSONEncoderObjectFields(t *testing.T) {
	tests := []struct {
		desc     string
		expected string
		f        func(zapcore.Encoder)
	}{
		{"binary", `"k":"YWIxMg=="`, func(e zapcore.Encoder) { e.AddBinary("k", []byte("ab12")) }},
		{"bool", `"k\\":true`, func(e zapcore.Encoder) { e.AddBool(`k\`, true) }}, // test key escaping once
		{"bool", `"k":true`, func(e zapcore.Encoder) { e.AddBool("k", true) }},
		{"bool", `"k":false`, func(e zapcore.Encoder) { e.AddBool("k", false) }},
		{"byteString", `"k":"v\\"`, func(e zapcore.Encoder) { e.AddByteString(`k`, []byte(`v\`)) }},
		{"byteString", `"k":"v"`, func(e zapcore.Encoder) { e.AddByteString("k", []byte("v")) }},
		{"byteString", `"k":""`, func(e zapcore.Encoder) { e.AddByteString("k", []byte{}) }},
		{"byteString", `"k":""`, func(e zapcore.Encoder) { e.AddByteString("k", nil) }},
		{"complex128", `"k":"1+2i"`, func(e zapcore.Encoder) { e.AddComplex128("k", 1+2i) }},
		{"complex64", `"k":"1+2i"`, func(e zapcore.Encoder) { e.AddComplex64("k", 1+2i) }},
		{"duration", `"k":0.000000001`, func(e zapcore.Encoder) { e.AddDuration("k", 1) }},
		{"float64", `"k":1`, func(e zapcore.Encoder) { e.AddFloat64("k", 1.0) }},
		{"float64", `"k":10000000000`, func(e zapcore.Encoder) { e.AddFloat64("k", 1e10) }},
		{"float64", `"k":"NaN"`, func(e zapcore.Encoder) { e.AddFloat64("k", math.NaN()) }},
		{"float64", `"k":"+Inf"`, func(e zapcore.Encoder) { e.AddFloat64("k", math.Inf(1)) }},
		{"float64", `"k":"-Inf"`, func(e zapcore.Encoder) { e.AddFloat64("k", math.Inf(-1)) }},
		{"float32", `"k":1`, func(e zapcore.Encoder) { e.AddFloat32("k", 1.0) }},
		{"float32", `"k":10000000000`, func(e zapcore.Encoder) { e.AddFloat32("k", 1e10) }},
		{"float32", `"k":"NaN"`, func(e zapcore.Encoder) { e.AddFloat32("k", float32(math.NaN())) }},
		{"float32", `"k":"+Inf"`, func(e zapcore.Encoder) { e.AddFloat32("k", float32(math.Inf(1))) }},
		{"float32", `"k":"-Inf"`, func(e zapcore.Encoder) { e.AddFloat32("k", float32(math.Inf(-1))) }},
		{"int", `"k":42`, func(e zapcore.Encoder) { e.AddInt("k", 42) }},
		{"int64", `"k":42`, func(e zapcore.Encoder) { e.AddInt64("k", 42) }},
		{"int32", `"k":42`, func(e zapcore.Encoder) { e.AddInt32("k", 42) }},
		{"int16", `"k":42`, func(e zapcore.Encoder) { e.AddInt16("k", 42) }},
		{"int8", `"k":42`, func(e zapcore.Encoder) { e.AddInt8("k", 42) }},
		{"string", `"k":"v\\"`, func(e zapcore.Encoder) { e.AddString(`k`, `v\`) }},
		{"string", `"k":"v"`, func(e zapcore.Encoder) { e.AddString("k", "v") }},
		{"string", `"k":""`, func(e zapcore.Encoder) { e.AddString("k", "") }},
		{"time", `"k":1`, func(e zapcore.Encoder) { e.AddTime("k", time.Unix(1, 0)) }},
		{"uint", `"k":42`, func(e zapcore.Encoder) { e.AddUint("k", 42) }},
		{"uint64", `"k":42`, func(e zapcore.Encoder) { e.AddUint64("k", 42) }},
		{"uint32", `"k":42`, func(e zapcore.Encoder) { e.AddUint32("k", 42) }},
		{"uint16", `"k":42`, func(e zapcore.Encoder) { e.AddUint16("k", 42) }},
		{"uint8", `"k":42`, func(e zapcore.Encoder) { e.AddUint8("k", 42) }},
		{"uintptr", `"k":42`, func(e zapcore.Encoder) { e.AddUintptr("k", 42) }},
		{
			desc:     "object (success)",
			expected: `"k":{"loggable":"yes"}`,
			f: func(e zapcore.Encoder) {
				assert.NoError(t, e.AddObject("k", loggable{true}), "Unexpected error calling MarshalLogObject.")
			},
		},
		{
			desc:     "object (error)",
			expected: `"k":{}`,
			f: func(e zapcore.Encoder) {
				assert.Error(t, e.AddObject("k", loggable{false}), "Expected an error calling MarshalLogObject.")
			},
		},
		{
			desc:     "object (with nested array)",
			expected: `"turducken":{"ducks":[{"in":"chicken"},{"in":"chicken"}]}`,
			f: func(e zapcore.Encoder) {
				assert.NoError(
					t,
					e.AddObject("turducken", turducken{}),
					"Unexpected error calling MarshalLogObject with nested ObjectMarshalers and ArrayMarshalers.",
				)
			},
		},
		{
			desc:     "array (with nested object)",
			expected: `"turduckens":[{"ducks":[{"in":"chicken"},{"in":"chicken"}]},{"ducks":[{"in":"chicken"},{"in":"chicken"}]}]`,
			f: func(e zapcore.Encoder) {
				assert.NoError(
					t,
					e.AddArray("turduckens", turduckens(2)),
					"Unexpected error calling MarshalLogObject with nested ObjectMarshalers and ArrayMarshalers.",
				)
			},
		},
		{
			desc:     "array (success)",
			expected: `"k":[true]`,
			f: func(e zapcore.Encoder) {
				assert.NoError(t, e.AddArray(`k`, loggable{true}), "Unexpected error calling MarshalLogArray.")
			},
		},
		{
			desc:     "array (error)",
			expected: `"k":[]`,
			f: func(e zapcore.Encoder) {
				assert.Error(t, e.AddArray("k", loggable{false}), "Expected an error calling MarshalLogArray.")
			},
		},
		{
			desc:     "reflect (success)",
			expected: `"k":{"loggable":"yes"}`,
			f: func(e zapcore.Encoder) {
				assert.NoError(t, e.AddReflected("k", map[string]string{"loggable": "yes"}), "Unexpected error JSON-serializing a map.")
			},
		},
		{
			desc:     "reflect (failure)",
			expected: "",
			f: func(e zapcore.Encoder) {
				assert.Error(t, e.AddReflected("k", noJSON{}), "Unexpected success JSON-serializing a noJSON.")
			},
		},
		{
			desc: "namespace",
			// EncodeEntry is responsible for closing all open namespaces.
			expected: `"outermost":{"outer":{"foo":1,"inner":{"foo":2,"innermost":{`,
			f: func(e zapcore.Encoder) {
				e.OpenNamespace("outermost")
				e.OpenNamespace("outer")
				e.AddInt("foo", 1)
				e.OpenNamespace("inner")
				e.AddInt("foo", 2)
				e.OpenNamespace("innermost")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assertOutput(t, tt.expected, tt.f)
		})
	}
}

func TestJSONEncoderArrays(t *testing.T) {
	tests := []struct {
		desc     string
		expected string // expect f to be called twice
		f        func(zapcore.ArrayEncoder)
	}{
		{"bool", `[true,true]`, func(e zapcore.ArrayEncoder) { e.AppendBool(true) }},
		{"byteString", `["k","k"]`, func(e zapcore.ArrayEncoder) { e.AppendByteString([]byte("k")) }},
		{"byteString", `["k\\","k\\"]`, func(e zapcore.ArrayEncoder) { e.AppendByteString([]byte(`k\`)) }},
		{"complex128", `["1+2i","1+2i"]`, func(e zapcore.ArrayEncoder) { e.AppendComplex128(1 + 2i) }},
		{"complex64", `["1+2i","1+2i"]`, func(e zapcore.ArrayEncoder) { e.AppendComplex64(1 + 2i) }},
		{"durations", `[0.000000002,0.000000002]`, func(e zapcore.ArrayEncoder) { e.AppendDuration(2) }},
		{"float64", `[3.14,3.14]`, func(e zapcore.ArrayEncoder) { e.AppendFloat64(3.14) }},
		{"float32", `[3.14,3.14]`, func(e zapcore.ArrayEncoder) { e.AppendFloat32(3.14) }},
		{"int", `[42,42]`, func(e zapcore.ArrayEncoder) { e.AppendInt(42) }},
		{"int64", `[42,42]`, func(e zapcore.ArrayEncoder) { e.AppendInt64(42) }},
		{"int32", `[42,42]`, func(e zapcore.ArrayEncoder) { e.AppendInt32(42) }},
		{"int16", `[42,42]`, func(e zapcore.ArrayEncoder) { e.AppendInt16(42) }},
		{"int8", `[42,42]`, func(e zapcore.ArrayEncoder) { e.AppendInt8(42) }},
		{"string", `["k","k"]`, func(e zapcore.ArrayEncoder) { e.AppendString("k") }},
		{"string", `["k\\","k\\"]`, func(e zapcore.ArrayEncoder) { e.AppendString(`k\`) }},
		{"times", `[1,1]`, func(e zapcore.ArrayEncoder) { e.AppendTime(time.Unix(1, 0)) }},
		{"uint", `[42,42]`, func(e zapcore.ArrayEncoder) { e.AppendUint(42) }},
		{"uint64", `[42,42]`, func(e zapcore.ArrayEncoder) { e.AppendUint64(42) }},
		{"uint32", `[42,42]`, func(e zapcore.ArrayEncoder) { e.AppendUint32(42) }},
		{"uint16", `[42,42]`, func(e zapcore.ArrayEncoder) { e.AppendUint16(42) }},
		{"uint8", `[42,42]`, func(e zapcore.ArrayEncoder) { e.AppendUint8(42) }},
		{"uintptr", `[42,42]`, func(e zapcore.ArrayEncoder) { e.AppendUintptr(42) }},
		{
			desc:     "arrays (success)",
			expected: `[[true],[true]]`,
			f: func(arr zapcore.ArrayEncoder) {
				assert.NoError(t, arr.AppendArray(zapcore.ArrayMarshalerFunc(func(inner zapcore.ArrayEncoder) error {
					inner.AppendBool(true)
					return nil
				})), "Unexpected error appending an array.")
			},
		},
		{
			desc:     "arrays (error)",
			expected: `[[true],[true]]`,
			f: func(arr zapcore.ArrayEncoder) {
				assert.Error(t, arr.AppendArray(zapcore.ArrayMarshalerFunc(func(inner zapcore.ArrayEncoder) error {
					inner.AppendBool(true)
					return errors.New("fail")
				})), "Expected an error appending an array.")
			},
		},
		{
			desc:     "objects (success)",
			expected: `[{"loggable":"yes"},{"loggable":"yes"}]`,
			f: func(arr zapcore.ArrayEncoder) {
				assert.NoError(t, arr.AppendObject(loggable{true}), "Unexpected error appending an object.")
			},
		},
		{
			desc:     "objects (error)",
			expected: `[{},{}]`,
			f: func(arr zapcore.ArrayEncoder) {
				assert.Error(t, arr.AppendObject(loggable{false}), "Expected an error appending an object.")
			},
		},
		{
			desc:     "reflect (success)",
			expected: `[{"foo":5},{"foo":5}]`,
			f: func(arr zapcore.ArrayEncoder) {
				assert.NoError(
					t,
					arr.AppendReflected(map[string]int{"foo": 5}),
					"Unexpected an error appending an object with reflection.",
				)
			},
		},
		{
			desc:     "reflect (error)",
			expected: `[]`,
			f: func(arr zapcore.ArrayEncoder) {
				assert.Error(
					t,
					arr.AppendReflected(noJSON{}),
					"Unexpected an error appending an object with reflection.",
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			f := func(enc zapcore.Encoder) error {
				return enc.AddArray("array", zapcore.ArrayMarshalerFunc(func(arr zapcore.ArrayEncoder) error {
					tt.f(arr)
					tt.f(arr)
					return nil
				}))
			}
			assertOutput(t, `"array":`+tt.expected, func(enc zapcore.Encoder) {
				err := f(enc)
				assert.NoError(t, err, "Unexpected error adding array to JSON encoder.")
			})
		})
	}
}
