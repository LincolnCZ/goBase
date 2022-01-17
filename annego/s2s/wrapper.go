package s2s

/*
#cgo LDFLAGS: -L${SRCDIR} -ls2sclient -luuid -lcrypto -lpthread -lstdc++ -lrt
#include <stdlib.h>
#include "s2sc.h"
*/
import "C"

import (
	"goBase/annego/packet"
	"unsafe"
)

// S2sMetaStatus
const (
	S2SMETA_OK = iota
	S2SMETA_DIED
)

const (
	S2S_ANY_TYPE   = 0
	S2S_TEXTPLAIN  = 128
	S2S_S2SDECODER = 129
)

const (
	S2S_SESSIONOFF = iota
	S2S_SESSIONON
	S2S_SESSIONBIND
	S2S_DNSERROR
	S2S_AUTHFAILURE
	S2S_ERROR
)

type PSubFilter struct {
	Filters []SubFilter
}

func (self *PSubFilter) GetURI() uint32 {
	return 0
}

func (self *PSubFilter) Marshal(pk *packet.Pack) {
	pk.PutUint32(uint32(len(self.Filters)))
	for _, f := range self.Filters {
		f.Marshal(pk)
	}
}

func (self *PSubFilter) Unmarshal(up *packet.Unpack) error {
	size, err := up.PopUint32()
	if err != nil {
		return err
	}
	self.Filters = make([]SubFilter, size)
	for i := uint32(0); i < size; i++ {
		if err = self.Filters[i].Unmarshal(up); err != nil {
			return err
		}
	}
	return nil
}

type PNotifyResult struct {
	Status uint32
	Metas  []S2SMeta
}

func (self *PNotifyResult) GetURI() uint32 {
	return 0
}

func (self *PNotifyResult) Marshal(pk *packet.Pack) {
	pk.PutUint32(self.Status)
	pk.PutUint32(uint32(len(self.Metas)))
	for _, m := range self.Metas {
		m.Marshal(pk)
	}
}

func (self *PNotifyResult) Unmarshal(up *packet.Unpack) error {
	var err error
	if self.Status, err = up.PopUint32(); err != nil {
		return err
	}

	var l uint32
	if l, err = up.PopUint32(); err != nil {
		return err
	}
	self.Metas = make([]S2SMeta, l)
	for i := uint32(0); i < l; i++ {
		if err = self.Metas[i].Unmarshal(up); err != nil {
			return err
		}
	}
	return nil
}

func initialize(name string, key string) int {
	cname := C.CString(name)
	ckey := C.CString(key)
	fd := C.initialize(cname, ckey, S2S_S2SDECODER)
	C.free(unsafe.Pointer(cname))
	C.free(unsafe.Pointer(ckey))
	return int(fd)
}

func subscribe(filters []SubFilter) bool {
	req := PSubFilter{Filters: filters}
	pk := packet.GetMarshalPack(&req)
	bi := pk.BodyBytes()
	input := C.struct_Buffer{
		buffer: C.CBytes(bi),
		size:   C.int(len(bi)),
	}

	res := C.subscribe(input)
	C.free(input.buffer)
	return res == 0
}

func pollNotify() *PNotifyResult {
	rsp := new(PNotifyResult)

	result := C.pollNotify()
	up := packet.NewUnpack(C.GoBytes(result.buffer, result.size))
	err := rsp.Unmarshal(up)
	C.free(result.buffer)
	if err != nil {
		panic(err)
	}
	return rsp
}

func setMine(data []byte) bool {
	cinfo := C.CBytes(data)
	infolen := len(data)
	res := C.setMine((*C.char)(cinfo), C.int(infolen))
	C.free(cinfo)
	return res == 0
}

func delMine() bool {
	res := C.delMine()
	return res == 0
}

func getMine() *S2SMeta {
	rsp := new(S2SMeta)

	result := C.getMine()
	if result.buffer == nil {
		return nil
	}
	up := packet.NewUnpack(C.GoBytes(result.buffer, result.size))
	err := rsp.Unmarshal(up)
	C.free(result.buffer)
	if err != nil {
		panic(err)
	}
	return rsp
}
