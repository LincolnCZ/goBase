package packet

import (
	"strconv"
	"testing"
)

type BenchProto struct {
	Flg  uint8
	U32  uint32
	U64  uint64
	List []string
	Map  map[uint32]string
}

func (self *BenchProto) GetURI() uint32 {
	return 0xa
}

func (self *BenchProto) Marshal(pk *Pack) {
	pk.PutUint8(self.Flg)
	pk.PutUint32(self.U32)
	pk.PutUint64(self.U64)
	pk.PutUint32(uint32(len(self.List)))
	for _, item_21 := range self.List {
		pk.PutShortStr(item_21)
	}
	pk.PutUint32(uint32(len(self.Map)))
	for key_25, val_25 := range self.Map {
		pk.PutUint32(key_25)
		pk.PutShortStr(val_25)
	}
}

func (self *BenchProto) MarshalReflect(pk *Pack) {
	pk.PutUint8(self.Flg)
	pk.PutUint32(self.U32)
	pk.PutUint64(self.U64)
	pk.PutSlice(self.List)
	pk.PutMap(self.Map)
}

func (self *BenchProto) MarshalDefault(pk *Pack) {
	DefaultMarshal(self, pk)
}

func (self *BenchProto) Unmarshal(up *Unpack) error {
	var err error
	if self.Flg, err = up.PopUint8(); err != nil {
		return err
	}
	if self.U32, err = up.PopUint32(); err != nil {
		return err
	}
	if self.U64, err = up.PopUint64(); err != nil {
		return err
	}
	var l_43 uint32
	if l_43, err = up.PopUint32(); err != nil {
		return err
	}
	self.List = make([]string, l_43)
	for i := uint32(0); i < l_43; i++ {
		var item_43 string
		if item_43, err = up.PopShortStr(); err != nil {
			return err
		}
		self.List[i] = item_43
	}
	var l_55 uint32
	if l_55, err = up.PopUint32(); err != nil {
		return err
	}
	self.Map = make(map[uint32]string, l_55)
	for i := uint32(0); i < l_55; i++ {
		var key_55 uint32
		var val_55 string
		if key_55, err = up.PopUint32(); err != nil {
			return err
		}
		if val_55, err = up.PopShortStr(); err != nil {
			return err
		}
		self.Map[key_55] = val_55
	}
	return nil
}

func (self *BenchProto) UnmarshalReflect(up *Unpack) error {
	var err error
	if self.Flg, err = up.PopUint8(); err != nil {
		return err
	}
	if self.U32, err = up.PopUint32(); err != nil {
		return err
	}
	if self.U64, err = up.PopUint64(); err != nil {
		return err
	}
	if err = up.PopSlice(&self.List); err != nil {
		return err
	}
	if err = up.PopMap(&self.Map); err != nil {
		return err
	}
	return nil
}

func (self *BenchProto) UnmarshalDefault(up *Unpack) error {
	return DefaultUnmarshal(self, up)
}

func NewBenchProto() *BenchProto {
	p := &BenchProto{}
	p.Flg = 10
	p.U32 = 0x010203
	p.U64 = 0xfffafbfcfd
	for i := 0; i < 50; i++ {
		p.List = append(p.List, strconv.Itoa(i))
	}
	p.Map = make(map[uint32]string)
	for i := 0; i < 50; i++ {
		p.Map[uint32(i)] = strconv.Itoa(i + 100)
	}
	return p
}

func BenchmarkMarshal(b *testing.B) {
	pk := NewPack()
	oldproto := NewBenchProto()
	var newproto BenchProto
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		oldproto.Marshal(pk)
		up := NewUnpack(pk.BodyBytes())
		newproto.Unmarshal(up)
	}
}

func BenchmarkReflect(b *testing.B) {
	pk := NewPack()
	oldproto := NewBenchProto()
	var newproto BenchProto
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		oldproto.MarshalReflect(pk)
		up := NewUnpack(pk.BodyBytes())
		newproto.UnmarshalReflect(up)
	}
}

func BenchmarkDefault(b *testing.B) {
	pk := NewPack()
	oldproto := NewBenchProto()
	var newproto BenchProto
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		oldproto.MarshalDefault(pk)
		up := NewUnpack(pk.BodyBytes())
		newproto.UnmarshalDefault(up)
	}
}
