package packet

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

// UnpackError 解包格式错误
type UnpackError struct {
	uri int
	msg string
}

func (e *UnpackError) Error() string {
	return fmt.Sprintf("unpack error: uri %d %s", e.uri, e.msg)
}

// Marshallable 协议接口，所有协议应该实现
type Marshallable interface {
	GetURI() uint32
	Marshal(pack *Pack)
	Unmarshal(unpack *Unpack) error
}

// Header 包头结构
type Header struct {
	Length  uint32
	URI     uint32
	ResCode uint16
}

// HeaderLength 包头长度: Length 4 + URI 4 + ResCode 2
const HeaderLength = 10

// MaxPacketLength 最大包长度 64MB
const MaxPacketLength = 64 * 1024 * 1024

// ResSuccess 成功完成
const ResSuccess = 200

// Pack 协议marshal到Pack
/* 典型用法如下，可使用GetMarshalPack()简化处理
pack := NewPack()
msg.Marshal(pack)
pack.PutHeader()
io.Write(pack.Bytes())
*/
type Pack struct {
	buf    []byte
	offset int
}

func NewPack() *Pack {
	return &Pack{make([]byte, 256), HeaderLength}
}

func (me *Pack) grow(n int) {
	if me.offset+n > len(me.buf) {
		newsize := me.offset + n + len(me.buf)
		bytes := make([]byte, newsize)
		copy(bytes, me.buf[:me.offset])
		me.buf = bytes
	}
}

func (me *Pack) Len() int {
	return me.offset
}

// Bytes 返回打包后的数据
func (me *Pack) Bytes() []byte {
	return me.buf[:me.offset]
}

// BodyBytes 返回不含包头的数据，嵌套协议封装使用
func (me *Pack) BodyBytes() []byte {
	return me.buf[HeaderLength:me.offset]
}

// Clear 清空数据恢复初始化状态
func (me *Pack) Clear() {
	me.buf = me.buf[0:0]
	me.offset = HeaderLength
}

func (me *Pack) PutBool(b bool) {
	if b {
		me.PutUint8(1)
	} else {
		me.PutUint8(0)
	}
}

func (me *Pack) PutUint8(u8 uint8) {
	me.grow(1)
	me.buf[me.offset] = byte(u8)
	me.offset += 1
}

func (me *Pack) PutUint16(u16 uint16) {
	me.grow(2)
	binary.LittleEndian.PutUint16(me.buf[me.offset:me.offset+2], u16)
	me.offset += 2
}

func (me *Pack) PutUint32(u32 uint32) {
	me.grow(4)
	binary.LittleEndian.PutUint32(me.buf[me.offset:me.offset+4], u32)
	me.offset += 4
}

func (me *Pack) PutUint64(u64 uint64) {
	me.grow(8)
	binary.LittleEndian.PutUint64(me.buf[me.offset:me.offset+8], u64)
	me.offset += 8
}

func (me *Pack) PutByteSlice(bytes []byte) {
	me.grow(len(bytes) + 4)
	binary.LittleEndian.PutUint32(me.buf[me.offset:me.offset+4], uint32(len(bytes)))
	me.offset += 4
	copy(me.buf[me.offset:], bytes)
	me.offset += len(bytes)
}

func (me *Pack) PutShortSlice(bytes []byte) {
	me.grow(len(bytes) + 2)
	binary.LittleEndian.PutUint16(me.buf[me.offset:me.offset+2], uint16(len(bytes)))
	me.offset += 2
	copy(me.buf[me.offset:], bytes)
	me.offset += len(bytes)
}

func (me *Pack) PutShortStr(s string) {
	me.grow(len(s) + 2)
	binary.LittleEndian.PutUint16(me.buf[me.offset:me.offset+2], uint16(len(s)))
	me.offset += 2
	copy(me.buf[me.offset:], s)
	me.offset += len(s)
}

func (me *Pack) PutLongStr(s string) {
	me.grow(len(s) + 4)
	binary.LittleEndian.PutUint32(me.buf[me.offset:me.offset+4], uint32(len(s)))
	me.offset += 4
	copy(me.buf[me.offset:], s)
	me.offset += len(s)
}

func (me *Pack) PutMarshallable(m Marshallable) {
	m.Marshal(me)
}

// PutValue 基于反射实现的任意类型marshal函数
func (me *Pack) PutValue(v reflect.Value) {
	switch v.Kind() {
	case reflect.Bool:
		me.PutBool(v.Bool())
	case reflect.Uint8:
		me.PutUint8(uint8(v.Uint()))
	case reflect.Uint16:
		me.PutUint16(uint16(v.Uint()))
	case reflect.Uint32:
		me.PutUint32(uint32(v.Uint()))
	case reflect.Uint64:
		me.PutUint64(v.Uint())
	case reflect.String:
		me.PutShortStr(v.String())

	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			me.PutShortSlice(v.Bytes())
		} else {
			me.putSliceImpl(v)
		}

	case reflect.Map:
		me.putMapImpl(v)

	case reflect.Struct:
		m, ok := v.Addr().Interface().(Marshallable)
		if !ok {
			panic(fmt.Sprintf("Pack.PutValue put struct %v not marshallable", v.Type()))
		}
		me.PutMarshallable(m)

	case reflect.Ptr:
		me.PutValue(v.Elem())

	default:
		panic(fmt.Sprintf("Pack.PutValue not support type %v", v.Kind()))
	}
}

func (me *Pack) putSliceImpl(v reflect.Value) {
	me.PutUint32(uint32(v.Len()))
	for i := 0; i < v.Len(); i++ {
		me.PutValue(v.Index(i))
	}
}

func (me *Pack) putMapImpl(v reflect.Value) {
	me.PutUint32(uint32(v.Len()))
	keys := v.MapKeys()
	for _, k := range keys {
		me.PutValue(k)
		me.PutValue(v.MapIndex(k))
	}
}

func (me *Pack) PutSlice(l interface{}) {
	val := reflect.ValueOf(l)
	if val.Kind() != reflect.Slice {
		panic(fmt.Sprintf("Pack.PutSlice put type %v", val.Kind()))
	}
	me.PutValue(val)
}

func (me *Pack) PutMap(m interface{}) {
	val := reflect.ValueOf(m)
	if val.Kind() != reflect.Map {
		panic(fmt.Sprintf("Pack.PutMap put type %v", val.Kind()))
	}
	me.PutValue(val)
}

// 包头封装使用
func (me *Pack) replaceUint16(pos int, u16 uint16) {
	binary.LittleEndian.PutUint16(me.buf[pos:pos+2], u16)
}

func (me *Pack) replaceUint32(pos int, u32 uint32) {
	binary.LittleEndian.PutUint32(me.buf[pos:pos+4], u32)
}

// PutHeader 应该在Marshal后调用，用来生成包头
func (me *Pack) PutHeader(uri uint32) {
	me.replaceUint32(0, uint32(me.offset))
	me.replaceUint32(4, uri)
	me.replaceUint16(8, ResSuccess)
}

// Unpack 协议从Unpack中unmarshal
/* 典型用法如下（无任何错误处理），可使用YYRegister.Unmarshal简化处理
unpack := NewUnpack(buf)
header, _ = unpack.PopHeader()
msg := Marshallable{}	// 根据header.URI选取具体协议
msg.Unmarshal(unpack)
*/
type Unpack struct {
	buf    []byte
	offset int
	valid  int
	header Header
}

// NewUnpack 从待解码数据生成
func NewUnpack(buf []byte) *Unpack {
	return &Unpack{buf, 0, len(buf), Header{0, 0, ResSuccess}}
}

// Header 返回包头数据，为空返回nil
func (me *Unpack) Header() *Header {
	if me.header.Length == 0 {
		return nil
	}
	return &me.header
}

// Offset Unpack已读取buf中数据长度
func (me *Unpack) Offset() int {
	return me.offset
}

// Length Unpack所包含buf的长度
func (me *Unpack) Length() int {
	return len(me.buf)
}

func (me *Unpack) checkSpace(size int) bool {
	return me.valid >= size+me.offset
}

func (me *Unpack) PopBool() (bool, error) {
	if !me.checkSpace(1) {
		return false, &UnpackError{int(me.header.URI), "PopBool"}
	}
	u8 := uint8(me.buf[me.offset])
	me.offset += 1

	if u8 == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (me *Unpack) PopUint8() (uint8, error) {
	if !me.checkSpace(1) {
		return 0, &UnpackError{int(me.header.URI), "PopUint8"}
	}
	u8 := uint8(me.buf[me.offset])
	me.offset += 1
	return u8, nil
}

func (me *Unpack) PopUint16() (uint16, error) {
	if !me.checkSpace(2) {
		return 0, &UnpackError{int(me.header.URI), "PopUint16"}
	}
	u16 := binary.LittleEndian.Uint16(me.buf[me.offset : me.offset+2])
	me.offset += 2
	return u16, nil
}

func (me *Unpack) PopUint32() (uint32, error) {
	if !me.checkSpace(4) {
		return 0, &UnpackError{int(me.header.URI), "PopUint32"}
	}
	u32 := binary.LittleEndian.Uint32(me.buf[me.offset : me.offset+4])
	me.offset += 4
	return u32, nil
}

func (me *Unpack) PopUint64() (uint64, error) {
	if !me.checkSpace(8) {
		return 0, &UnpackError{int(me.header.URI), "PopUint64"}
	}
	u64 := binary.LittleEndian.Uint64(me.buf[me.offset : me.offset+8])
	me.offset += 8
	return u64, nil
}

func (me *Unpack) PopShortStr() (string, error) {
	length, err := me.PopUint16()
	if err != nil {
		return "", err
	}
	if !me.checkSpace(int(length)) {
		s := fmt.Sprintf("PopShortStr %d", length)
		return "", &UnpackError{int(me.header.URI), s}
	}
	s := string(me.buf[me.offset : me.offset+int(length)])
	me.offset += int(length)
	return s, nil
}

func (me *Unpack) PopLongStr() (string, error) {
	length, err := me.PopUint32()
	if err != nil {
		return "", err
	}
	if !me.checkSpace(int(length)) {
		s := fmt.Sprintf("PopLongStr %d", length)
		return "", &UnpackError{int(me.header.URI), s}
	}
	s := string(me.buf[me.offset : me.offset+int(length)])
	me.offset += int(length)
	return s, nil
}

func (me *Unpack) PopByteSlice() ([]byte, error) {
	length, err := me.PopUint32()
	if err != nil {
		return make([]byte, 0), err
	}
	if !me.checkSpace(int(length)) {
		s := fmt.Sprintf("PopByteSlice %d", length)
		return make([]byte, 0), &UnpackError{int(me.header.URI), s}
	}
	s := me.buf[me.offset : me.offset+int(length)]
	me.offset += int(length)
	return s, nil
}

func (me *Unpack) PopShortSlice() ([]byte, error) {
	length, err := me.PopUint16()
	if err != nil {
		return make([]byte, 0), err
	}
	if !me.checkSpace(int(length)) {
		s := fmt.Sprintf("PopShortSlice %d", length)
		return make([]byte, 0), &UnpackError{int(me.header.URI), s}
	}
	s := me.buf[me.offset : me.offset+int(length)]
	me.offset += int(length)
	return s, nil
}

func (me *Unpack) PopMarshallable(m Marshallable) error {
	return m.Unmarshal(me)
}

// PopValue 基于反射实现的任意类型unmarshal函数，注意v必须是CanSet
func (me *Unpack) PopValue(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Bool:
		b, err := me.PopBool()
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Uint8:
		u8, err := me.PopUint8()
		if err != nil {
			return err
		}
		v.SetUint(uint64(u8))
	case reflect.Uint16:
		u16, err := me.PopUint16()
		if err != nil {
			return err
		}
		v.SetUint(uint64(u16))
	case reflect.Uint32:
		u32, err := me.PopUint32()
		if err != nil {
			return err
		}
		v.SetUint(uint64(u32))
	case reflect.Uint64:
		u64, err := me.PopUint64()
		if err != nil {
			return err
		}
		v.SetUint(u64)

	case reflect.String:
		str, err := me.PopShortStr()
		if err != nil {
			return err
		}
		v.SetString(str)

	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			bt, err := me.PopShortSlice()
			if err != nil {
				return err
			}
			v.SetBytes(bt)
		} else {
			return me.popSliceImpl(v)
		}

	case reflect.Map:
		return me.popMapImpl(v)

	case reflect.Struct:
		m, ok := v.Addr().Interface().(Marshallable)
		if !ok {
			panic(fmt.Sprintf("Unpack.PopValue put struct %v not marshallable", v.Type()))
		}
		return me.PopMarshallable(m)

	case reflect.Ptr:
		nval := reflect.New(v.Type().Elem())
		if err := me.PopValue(nval.Elem()); err != nil {
			return err
		}
		v.Set(nval)

	default:
		panic(fmt.Sprintf("Unpack.PopValue not support type %v", v.Kind()))
	}
	return nil
}

func (me *Unpack) popSliceImpl(v reflect.Value) error {
	l, err := me.PopUint32()
	if err != nil {
		return err
	}
	count := int(l)

	newval := reflect.MakeSlice(v.Type(), count, count)
	for i := 0; i < count; i++ {
		if err := me.PopValue(newval.Index(i)); err != nil {
			return err
		}
	}
	v.Set(newval)
	return nil
}

func (me *Unpack) popMapImpl(v reflect.Value) error {
	l, err := me.PopUint32()
	if err != nil {
		return err
	}
	count := int(l)

	tp := v.Type()
	newval := reflect.MakeMap(tp)
	for i := 0; i < count; i++ {
		key := reflect.New(tp.Key()).Elem()
		val := reflect.New(tp.Elem()).Elem()
		if err := me.PopValue(key); err != nil {
			return err
		}
		if err := me.PopValue(val); err != nil {
			return err
		}
		newval.SetMapIndex(key, val)
	}
	v.Set(newval)
	return nil
}

// PopSlice 应传入指向slice的指针
func (me *Unpack) PopSlice(l interface{}) error {
	val := reflect.ValueOf(l).Elem()
	if val.Kind() != reflect.Slice {
		panic(fmt.Sprintf("Unpack.PopSlice put type %v", val.Kind()))
	}
	return me.PopValue(val)
}

// PopMap 应传入指向map的指针
func (me *Unpack) PopMap(m interface{}) error {
	val := reflect.ValueOf(m).Elem()
	if val.Kind() != reflect.Map {
		panic(fmt.Sprintf("Unpack.PopMap put type %v", val.Kind()))
	}
	return me.PopValue(val)
}

// PopHeader 应该在Unmarshal前调用，用来解析包头
func (me *Unpack) PopHeader() (*Header, error) {
	length, err := me.PopUint32()
	if err != nil {
		return nil, err
	}
	uri, err := me.PopUint32()
	if err != nil {
		return nil, err
	}
	resCode, err := me.PopUint16()
	if err != nil {
		return nil, err
	}
	me.header = Header{length, uri, resCode}
	if int(length) < len(me.buf) {
		me.valid = int(length)
	}
	return &me.header, nil
}

// DefaultMarshal 基于反射实现的默认Marshal函数
func DefaultMarshal(proto Marshallable, pack *Pack) {
	v := reflect.ValueOf(proto).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		vi := v.Field(i)
		tag := t.Field(i).Tag.Get("yyp")

		if tag == "" {
			pack.PutValue(vi)
		} else {
			switch tag {
			case "uint8":
				pack.PutUint8(uint8(vi.Uint()))
			case "uint16":
				pack.PutUint16(uint16(vi.Uint()))
			case "uint32":
				pack.PutUint32(uint32(vi.Uint()))
			case "uint64":
				pack.PutUint64(vi.Uint())
			case "str":
				if vi.Kind() == reflect.String {
					pack.PutShortStr(vi.String())
				} else {
					pack.PutShortSlice(vi.Bytes())
				}
			case "str32":
				if vi.Kind() == reflect.String {
					pack.PutLongStr(vi.String())
				} else {
					pack.PutByteSlice(vi.Bytes())
				}
			case "-":
				// do nothing
			default:
				panic(fmt.Sprintf("DefaultMarshal yyp tag unknown: %s", tag))
			}
		}
	}
}

// DefaultUnmarshal 基于反射实现的默认Unmarshal函数
func DefaultUnmarshal(proto Marshallable, unpack *Unpack) error {
	v := reflect.ValueOf(proto).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		vi := v.Field(i)
		tag := t.Field(i).Tag.Get("yyp")

		if tag == "" {
			if err := unpack.PopValue(vi); err != nil {
				return err
			}
		} else {
			switch tag {
			case "uint8":
				u8, err := unpack.PopUint8()
				if err != nil {
					return err
				}
				vi.SetUint(uint64(u8))
			case "uint16":
				u16, err := unpack.PopUint16()
				if err != nil {
					return err
				}
				vi.SetUint(uint64(u16))
			case "uint32":
				u32, err := unpack.PopUint32()
				if err != nil {
					return err
				}
				vi.SetUint(uint64(u32))
			case "uint64":
				u64, err := unpack.PopUint64()
				if err != nil {
					return err
				}
				vi.SetUint(u64)
			case "str":
				if vi.Kind() == reflect.String {
					s, err := unpack.PopShortStr()
					if err != nil {
						return err
					}
					vi.SetString(s)
				} else {
					b, err := unpack.PopShortSlice()
					if err != nil {
						return err
					}
					vi.SetBytes(b)
				}
			case "str32":
				if vi.Kind() == reflect.String {
					s, err := unpack.PopLongStr()
					if err != nil {
						return err
					}
					vi.SetString(s)
				} else {
					b, err := unpack.PopByteSlice()
					if err != nil {
						return err
					}
					vi.SetBytes(b)
				}
			case "-":
				// do nothing
			default:
				panic(fmt.Sprintf("DefaultUnmarshal yyp tag unknown: %s", tag))
			}
		}
	}
	return nil
}
