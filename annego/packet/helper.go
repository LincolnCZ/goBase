package packet

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrInputNotEnough = errors.New("packet: input not enought")

// GetMarshalPack 返回msg打包完成的Pack
func GetMarshalPack(msg Marshallable) *Pack {
	pk := NewPack()
	msg.Marshal(pk)
	pk.PutHeader(msg.GetURI())
	return pk
}

// MarshalBody 消息编码，不包含协议头
func MarshalBody(msg Marshallable) []byte {
	pk := NewPack()
	msg.Marshal(pk)
	return pk.BodyBytes()
}

// MarshalBody 消息解码，不包含协议头
func UnmarshalBody(data []byte, msg Marshallable) error {
	up := NewUnpack(data)
	return msg.Unmarshal(up)
}

type YYRegister struct {
	register map[uint32]reflect.Type
}

func NewYYRegister() *YYRegister {
	return &YYRegister{make(map[uint32]reflect.Type)}
}

func (reg *YYRegister) Register(msg Marshallable) bool {
	_, ok := reg.register[msg.GetURI()]
	if !ok {
		reg.register[msg.GetURI()] = reflect.TypeOf(msg).Elem()
	}
	return !ok
}

// Unmarshal 直接解析Unpack
func (reg *YYRegister) Unmarshal(unpack *Unpack) (Marshallable, error) {
	var header *Header
	var err error
	if header = unpack.Header(); header == nil {
		header, err = unpack.PopHeader()
		if err != nil {
			return nil, err
		}
	}
	msgtype, ok := reg.register[header.URI]
	if !ok {
		return nil, fmt.Errorf("not register uri:%d", header.URI)
	}
	msgvalue := reflect.New(msgtype)
	msg := msgvalue.Interface().(Marshallable)
	err = msg.Unmarshal(unpack)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// UnmarshalBytes 从数据中进行解包，返回结果
// msg: 解析的协议包
// readsize: 成功解析的数据长度
// err: 解包错误
func (reg *YYRegister) UnmarshalBytes(data []byte) (msg Marshallable, readsize int, err error) {
	msg = nil
	readsize = 0
	err = nil

	unpack := NewUnpack(data)
	if unpack.Length() <= HeaderLength {
		err = ErrInputNotEnough
		return
	}
	header, _ := unpack.PopHeader()
	if MaxPacketLength < header.Length {
		err = fmt.Errorf("unmarshal header length too long, length %d uri %d", header.Length, header.URI)
		return
	}
	if unpack.Length() < int(header.Length) {
		err = ErrInputNotEnough
		return
	}

	msg, err = reg.Unmarshal(unpack)
	if err != nil {
		return
	}
	// 正常解包但长度错误，可能是header.Length错误。判定为解包失败
	if err == nil && unpack.Offset() != int(header.Length) {
		err = fmt.Errorf("unmarshal error length: %d %d", unpack.Offset(), header.Length)
		msg = nil
		return
	}
	readsize = int(header.Length)
	return
}

var DefaultYYRegister *YYRegister = NewYYRegister()
