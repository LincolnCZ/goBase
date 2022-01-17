package main

import (
	"goBase/annego/packet"
)

type PTest struct {
	Int  uint32
	Str  string
	List []uint32
}

func (self *PTest) GetURI() uint32 {
	return 1
}

func (self *PTest) Marshal(pk *packet.Pack) {
	pk.PutUint32(self.Int)
	pk.PutShortStr(self.Str)
	pk.PutUint32(uint32(len(self.List)))
	for _, item := range self.List {
		pk.PutUint32(item)
	}
}

func (self *PTest) Unmarshal(up *packet.Unpack) error {
	var err error
	if self.Int, err = up.PopUint32(); err != nil {
		return err
	}
	if self.Str, err = up.PopShortStr(); err != nil {
		return err
	}
	var l uint32
	if l, err = up.PopUint32(); err != nil {
		return err
	}
	self.List = make([]uint32, l)
	for i := uint32(0); i < l; i++ {
		var item uint32
		if item, err = up.PopUint32(); err != nil {
			return err
		}
		self.List[i] = item
	}
	return nil
}

type PTestRes struct {
	Int uint32
	Str string
	Sum uint32
}

func (self *PTestRes) GetURI() uint32 {
	return 2
}

func (self *PTestRes) Marshal(pk *packet.Pack) {
	pk.PutUint32(self.Int)
	pk.PutShortStr(self.Str)
	pk.PutUint32(self.Sum)
}

func (self *PTestRes) Unmarshal(up *packet.Unpack) error {
	var err error
	if self.Int, err = up.PopUint32(); err != nil {
		return err
	}
	if self.Str, err = up.PopShortStr(); err != nil {
		return err
	}
	if self.Sum, err = up.PopUint32(); err != nil {
		return err
	}
	return nil
}
