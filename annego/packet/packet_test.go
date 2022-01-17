package packet

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SimpleProto struct {
	T uint8
	I uint32
	S string
}

func (self *SimpleProto) GetURI() uint32 {
	return 1
}

func (self *SimpleProto) Marshal(pk *Pack) {
	pk.PutUint8(self.T)
	pk.PutUint32(self.I)
	pk.PutShortStr(self.S)
}

func (self *SimpleProto) Unmarshal(up *Unpack) error {
	var err error
	if self.T, err = up.PopUint8(); err != nil {
		return err
	}
	if self.I, err = up.PopUint32(); err != nil {
		return err
	}
	if self.S, err = up.PopShortStr(); err != nil {
		return err
	}
	return nil
}

func TestSimpleMarshal(t *testing.T) {
	proto := &SimpleProto{1, 0x0102, "abcde"}

	// 测试编码
	pk := NewPack()
	proto.Marshal(pk)
	pk.PutHeader(proto.GetURI())
	body := pk.Bytes()

	assert.Equal(t, 22, pk.Len())
	target := []byte{
		22, 0, 0, 0, 1, 0, 0, 0, 200, 0, // header
		1, 2, 1, 0, 0, 5, 0, 97, 98, 99, 100, 101,
	}
	assert.Equal(t, target, body)

	// 测试解码
	body = append(body, 1, 2, 3) // 添加冗余字符
	up := NewUnpack(body)
	header, err := up.PopHeader()
	assert.NoError(t, err)
	assert.Equal(t, header, &Header{
		Length:  22,
		URI:     proto.GetURI(),
		ResCode: ResSuccess,
	})

	rsp := new(SimpleProto)
	err = rsp.Unmarshal(up)
	assert.NoError(t, err)
	assert.Equal(t, proto, rsp)
	assert.Equal(t, 22, up.Offset())
}

func TestSimpleUnmarshalErr(t *testing.T) {
	buff := []byte{
		20, 0, 0, 0, 1, 0, 0, 0, 200, 0, // first byte should be 22
		1, 2, 1, 0, 0, 5, 0, 97, 98, 99, 100, 101,
	}
	up := NewUnpack(buff)
	proto := SimpleProto{}
	header, err := up.PopHeader()
	assert.NoError(t, err)
	assert.Equal(t, header, &Header{
		Length:  20,
		URI:     proto.GetURI(),
		ResCode: ResSuccess,
	})

	err = proto.Unmarshal(up)
	assert.Error(t, err)
}

func TestSimpleMarshalBody(t *testing.T) {
	proto := &SimpleProto{
		T: 1,
		I: 0x0102,
		S: "hello",
	}

	// Marshal
	pk := NewPack()
	proto.Marshal(pk)
	buff := pk.BodyBytes()
	target := []byte{1, 2, 1, 0, 0, 5, 0, 104, 101, 108, 108, 111}
	assert.Equal(t, buff, target)

	// Unmarshal
	up := NewUnpack(buff)
	newproto := &SimpleProto{}
	err := newproto.Unmarshal(up)
	assert.NoError(t, err)
	assert.Equal(t, proto, newproto)
}

type DeepProto struct {
	B    bool
	Flag uint16
	ID   uint64
	S    string
	L    []*SimpleProto
	M    map[uint32]string
}

func (self *DeepProto) GetURI() uint32 {
	return 2
}

func (self *DeepProto) Marshal(pk *Pack) {
	pk.PutBool(self.B)
	pk.PutUint16(self.Flag)
	pk.PutUint64(self.ID)
	pk.PutShortStr(self.S)

	pk.PutUint32(uint32(len(self.L)))
	for _, item := range self.L {
		item.Marshal(pk)
	}

	pk.PutUint32(uint32(len(self.M)))
	for key, val := range self.M {
		pk.PutUint32(key)
		pk.PutShortStr(val)
	}
}

func (self *DeepProto) MarshalReflect(pk *Pack) {
	DefaultMarshal(self, pk)
}

func (self *DeepProto) Unmarshal(up *Unpack) error {
	var err error
	if self.B, err = up.PopBool(); err != nil {
		return err
	}
	if self.Flag, err = up.PopUint16(); err != nil {
		return err
	}
	if self.ID, err = up.PopUint64(); err != nil {
		return err
	}
	if self.S, err = up.PopShortStr(); err != nil {
		return err
	}

	var l uint32
	if l, err = up.PopUint32(); err != nil {
		return err
	}
	self.L = make([]*SimpleProto, l)
	for i := uint32(0); i < l; i++ {
		item := &SimpleProto{}
		if err = item.Unmarshal(up); err != nil {
			return err
		}
		self.L[i] = item
	}

	if l, err = up.PopUint32(); err != nil {
		return err
	}
	self.M = make(map[uint32]string, l)
	for i := uint32(0); i < l; i++ {
		key, err := up.PopUint32()
		if err != nil {
			return err
		}
		val, err := up.PopShortStr()
		if err != nil {
			return err
		}
		self.M[key] = val
	}
	return nil
}

func (self *DeepProto) UnmarshalReflect(up *Unpack) error {
	return DefaultUnmarshal(self, up)
}

func NewDeepProto() *DeepProto {
	proto := &DeepProto{}
	proto.B = true
	proto.Flag = 10
	proto.ID = 0x010203040506
	for i := 0; i < 128; i++ {
		proto.S += strconv.Itoa(i)
	}
	proto.L = make([]*SimpleProto, 20)
	for i := 0; i < 20; i++ {
		sp := &SimpleProto{}
		sp.T = 0
		sp.I = uint32(i)
		sp.S = strconv.Itoa(i)
		proto.L[i] = sp
	}
	proto.M = make(map[uint32]string)
	for i := 0; i < 200; i++ {
		proto.M[uint32(i)] = fmt.Sprint(i)
	}
	return proto
}

func TestDeepProto(t *testing.T) {
	// 生成协议
	var err error
	oldproto := NewDeepProto()

	// Marshal
	pk := NewPack()
	oldproto.Marshal(pk)
	pk.PutHeader(oldproto.GetURI())
	buff := pk.Bytes()
	// Unmarshal
	up := NewUnpack(buff)
	header, err := up.PopHeader()
	if assert.NoError(t, err) {
		assert.Equal(t, header.URI, oldproto.GetURI())
	}

	newproto := &DeepProto{}
	err = newproto.Unmarshal(up)
	if assert.NoError(t, err) {
		assert.Equal(t, oldproto, newproto)
	}
}

func TestReflectMarshal(t *testing.T) {
	// 生成协议
	oldproto := NewDeepProto()

	// Marshal
	pk1 := NewPack()
	oldproto.Marshal(pk1)
	pk1.PutHeader(oldproto.GetURI())
	buf1 := pk1.Bytes()
	// MarshalReflect
	pk2 := NewPack()
	oldproto.MarshalReflect(pk2)
	pk2.PutHeader(oldproto.GetURI())
	buf2 := pk2.Bytes()

	// Go语言map遍历顺序不确定
	// 所以无法直接比较marshal之后二进制格式是否相等
	// 没有想到有什么方法绕过这点，所以直接交叉unmarshal之后进行比较
	assert.Equal(t, len(buf1), len(buf2))

	// UnmarshalReflect
	up1 := NewUnpack(buf1)
	up1.PopHeader()
	proto1 := &DeepProto{}
	assert.NoError(t, proto1.UnmarshalReflect(up1))
	assert.Equal(t, oldproto, proto1)

	// Unmarshal
	up2 := NewUnpack(buf2)
	up2.PopHeader()
	proto2 := &DeepProto{}
	assert.NoError(t, proto2.Unmarshal(up2))
	assert.Equal(t, oldproto, proto2)
}

type BytesProto struct {
	t uint16
	i uint64
	s []byte
	l []byte
}

func (self *BytesProto) GetURI() uint32 {
	return 1
}

func (self *BytesProto) Marshal(pk *Pack) {
	pk.PutUint16(self.t)
	pk.PutUint64(self.i)
	pk.PutShortSlice(self.s)
	pk.PutByteSlice(self.l)
}

func (self *BytesProto) Unmarshal(up *Unpack) error {
	var err error
	if self.t, err = up.PopUint16(); err != nil {
		return err
	}
	if self.i, err = up.PopUint64(); err != nil {
		return err
	}
	if self.s, err = up.PopShortSlice(); err != nil {
		return err
	}
	if self.l, err = up.PopByteSlice(); err != nil {
		return err
	}
	return nil
}

func TestBytesProto(t *testing.T) {
	// 生成协议
	oldproto := BytesProto{}
	oldproto.t = 0xdeaf
	oldproto.i = 0xdeaf000001012
	oldproto.s = make([]byte, 128)
	oldproto.l = make([]byte, 128)
	for i := 0; i < 128; i++ {
		oldproto.s[i] = byte(i)
		oldproto.l[i] = byte(i)
	}

	var err error
	// Marshal
	pk := NewPack()
	oldproto.Marshal(pk)
	pk.PutHeader(oldproto.GetURI())
	buff := pk.Bytes()

	// Unmarshal
	up := NewUnpack(buff)
	if _, err = up.PopHeader(); err != nil {
		t.Error("unmarshal header", err)
	}
	newproto := BytesProto{}
	if err = newproto.Unmarshal(up); err != nil {
		t.Error("unmarshal protocol", err)
	}
	// 对比结果
	assert.Equal(t, oldproto, newproto)
}

type ContainProto struct {
	B []byte
	L []string
	M map[uint32]*SimpleProto
}

func (self *ContainProto) GetURI() uint32 {
	return 0
}

func (self *ContainProto) Marshal(pk *Pack) {
	pk.PutSlice(self.B)
	pk.PutSlice(self.L)
	pk.PutMap(self.M)
}

func (self *ContainProto) MarshalReflect(pk *Pack) {
	DefaultMarshal(self, pk)
}

func (self *ContainProto) Unmarshal(up *Unpack) error {
	var err error
	if err = up.PopSlice(&self.B); err != nil {
		return err
	}
	if err = up.PopSlice(&self.L); err != nil {
		return err
	}
	if err = up.PopMap(&self.M); err != nil {
		return err
	}
	return nil
}

func (self *ContainProto) UnmarshalReflect(up *Unpack) error {
	return DefaultUnmarshal(self, up)
}

func TestContainProto(t *testing.T) {
	// 生成协议
	oldproto := &ContainProto{
		B: make([]byte, 32),
		L: make([]string, 32),
		M: make(map[uint32]*SimpleProto),
	}
	for i := 0; i < 32; i++ {
		oldproto.B[i] = byte(i)
		oldproto.L[i] = fmt.Sprint(i)
		oldproto.M[uint32(i)] = &SimpleProto{
			T: uint8(i),
			I: uint32(i + 100),
			S: "abcde",
		}
	}

	// 基本编解码测试
	pk := NewPack()
	oldproto.Marshal(pk)
	up1 := NewUnpack(pk.BodyBytes())
	newproto1 := &ContainProto{}
	if assert.NoError(t, newproto1.Unmarshal(up1)) {
		assert.Equal(t, oldproto, newproto1)
	}

	// 反射编解码测试
	pk.Clear()
	oldproto.MarshalReflect(pk)
	up2 := NewUnpack(pk.BodyBytes())
	newproto2 := &ContainProto{}
	if assert.NoError(t, newproto2.UnmarshalReflect(up2)) {
		assert.Equal(t, oldproto, newproto2)
	}
}

type TagProto struct {
	U8   uint   `yyp:"uint8"`
	U16  uint   `yyp:"uint16"`
	U32  uint   `yyp:"uint32"`
	U64  uint   `yyp:"uint64"`
	Skip uint64 `yyp:"-"`
	S16  string `yyp:"str"`
	S32  string `yyp:"str32"`
	B16  []byte `yyp:"str"`
	B32  []byte `yyp:"str32"`
}

func (self *TagProto) GetURI() uint32 {
	return 0
}

func (self *TagProto) Marshal(pk *Pack) {
	DefaultMarshal(self, pk)
}

func (self *TagProto) Unmarshal(up *Unpack) error {
	return DefaultUnmarshal(self, up)
}

func TestTagMarshal(t *testing.T) {
	msg := &TagProto{
		U8:   8,
		U16:  16,
		U32:  32,
		U64:  64,
		Skip: 1,
		S16:  "ss16",
		S32:  "ss32",
		B16:  []byte("bs16"),
		B32:  []byte("bs32"),
	}

	// marshal packet
	pk := NewPack()
	msg.Marshal(pk)
	body := pk.BodyBytes()
	assert.Equal(t, 43, len(body))

	// unmarshal packet
	rsp := new(TagProto)
	rsp.Skip = 2
	up := NewUnpack(body)
	err := rsp.Unmarshal(up)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), rsp.Skip)
	rsp.Skip = 1
	assert.Equal(t, msg, rsp)
}
