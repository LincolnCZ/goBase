package yyserver

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufferRead(t *testing.T) {
	buffer := newReadBuffer()
	reader := new(bytes.Buffer)

	reader.WriteString("123456")
	n, err := buffer.ReadIO(reader)
	assert.Nil(t, err)
	assert.Equal(t, int(n), 6)
	assert.Equal(t, buffer.Len(), 6)
	assert.Equal(t, buffer.Seek(), []byte("123456"))

	buffer.HasRead(3)
	assert.Equal(t, buffer.Len(), 3)
	assert.Equal(t, buffer.Seek(), []byte("456"))

	buffer.HasRead(3)
	assert.Equal(t, buffer.Len(), 0)
	assert.Equal(t, len(buffer.Seek()), 0)
}

func TestBufferGrow(t *testing.T) {
	buffer := newReadBuffer()
	const total = 10 * 1024
	b := make([]byte, total)
	for i := 0; i < total; i++ {
		b[i] = byte(i % 256)
	}
	reader := bytes.NewBuffer(b)
	var err error

	n1, err := buffer.ReadIO(reader)
	assert.Nil(t, err)
	assert.Equal(t, int(n1), buffer.readsize)
	assert.Equal(t, buffer.Len(), buffer.readsize)
	assert.Equal(t, buffer.Seek(), b[:buffer.readsize])

	n2, err := buffer.ReadIO(reader)
	assert.Nil(t, err)
	assert.Equal(t, int(n2), buffer.readsize)
	assert.Equal(t, buffer.Len(), int(n1+n2))
	assert.Equal(t, buffer.Seek(), b[:n1+n2])

	buffer.HasRead(1024)
	assert.Equal(t, buffer.Len(), int(n1+n2-1024))
	assert.Equal(t, buffer.Seek(), b[1024:n1+n2])

	n3, err := buffer.ReadIO(reader)
	assert.Nil(t, err)
	assert.Equal(t, n3, total-n1-n2)
	assert.Equal(t, buffer.Len(), int(total-1024))
}
