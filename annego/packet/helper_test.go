package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type simpleProto struct {
	i uint32
	s string
}

func (self *simpleProto) GetURI() uint32 {
	return 1
}

func (self *simpleProto) Marshal(pk *Pack) {
	pk.PutUint32(self.i)
	pk.PutShortStr(self.s)
}

func (self *simpleProto) Unmarshal(up *Unpack) error {
	var err error
	if self.i, err = up.PopUint32(); err != nil {
		return err
	}
	if self.s, err = up.PopShortStr(); err != nil {
		return err
	}
	return nil
}

func TestUnmarshalBytes(t *testing.T) {
	register := NewYYRegister()
	register.Register(&simpleProto{})

	var sendbuf []byte
	for i := 0; i < 2; i++ {
		newmsg := &simpleProto{uint32(i), "abcdefg123456789"}
		pack := GetMarshalPack(newmsg)
		sendbuf = append(sendbuf, pack.Bytes()...)
	}

	var msg Marshallable
	var readsize int
	var err error
	// 缓存区比包头小
	msg, readsize, err = register.UnmarshalBytes(sendbuf[0:8])
	assert.Equal(t, err, ErrInputNotEnough)

	// 缓存区包含包头，但是不够包长度
	msg, readsize, err = register.UnmarshalBytes(sendbuf[0:12])
	assert.Equal(t, err, ErrInputNotEnough)

	// 缓存区包含一个完整包以上的数据
	msg, readsize, err = register.UnmarshalBytes(sendbuf)
	assert.NotNil(t, msg)
	assert.Equal(t, readsize, len(sendbuf)/2)
	assert.NoError(t, err)

	// 包头长度异常
	sendbuf[0] = 54
	msg, readsize, err = register.UnmarshalBytes(sendbuf)
	assert.Error(t, err)
	assert.NotEqual(t, err, ErrInputNotEnough)
	t.Logf("unmarshal error %v", err)
}

func TestUnmarshalBody(t *testing.T) {
	msg1 := &simpleProto{1234, "abcdefg123456789"}
	msg2 := new(simpleProto)
	body := MarshalBody(msg1)
	err := UnmarshalBody(body, msg2)
	assert.NoError(t, err)
	assert.Equal(t, msg1.i, msg2.i)
	assert.Equal(t, msg1.s, msg2.s)
}
