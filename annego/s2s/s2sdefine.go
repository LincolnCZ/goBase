package s2s

import (
	"fmt"
	"math"
	"net"

	"goBase/annego/packet"
	"goBase/annego/util"

	"gopkg.in/mgo.v2/bson"
)

// SubFilter S2S订阅信息
type SubFilter struct {
	Name    string // 服务名称
	GroupID int32  // 机房id
	S2SType int32  // 编解码协议，默认S2S_ANY_TYPE
}

func (self *SubFilter) GetURI() uint32 {
	return 0
}

func (self *SubFilter) Marshal(pk *packet.Pack) {
	pk.PutShortStr(self.Name)
	pk.PutUint32(uint32(self.GroupID))
	pk.PutUint32(uint32(self.S2SType))
}

func (self *SubFilter) Unmarshal(up *packet.Unpack) error {
	var err error
	var u uint32
	if self.Name, err = up.PopShortStr(); err != nil {
		return err
	}

	if u, err = up.PopUint32(); err != nil {
		return err
	}
	self.GroupID = int32(u)

	if u, err = up.PopUint32(); err != nil {
		return err
	}
	self.S2SType = int32(u)
	return nil
}

// S2SMeta s2s分配的元信息
type S2SMeta struct {
	ServerID  int64  // server分配的惟一标识id;
	MetaType  int    // MetaType
	Name      string // 服务名称；
	GroupID   int32  // 机房id
	Data      []byte
	Timestamp int64
	Statu     int // S2sMetaStatus
}

func (self *S2SMeta) GetURI() uint32 {
	return 0
}

func (self *S2SMeta) Marshal(pk *packet.Pack) {
	pk.PutUint64(uint64(self.ServerID))
	pk.PutUint32(uint32(self.MetaType))
	pk.PutShortStr(self.Name)
	pk.PutUint32(uint32(self.GroupID))
	pk.PutByteSlice(self.Data)
	pk.PutUint64(uint64(self.Timestamp))
	pk.PutUint32(uint32(self.Statu))
}

func (self *S2SMeta) Unmarshal(up *packet.Unpack) error {
	var err error
	var u32 uint32
	var u64 uint64

	if u64, err = up.PopUint64(); err != nil {
		return err
	}
	self.ServerID = int64(u64)

	if u32, err = up.PopUint32(); err != nil {
		return err
	}
	self.MetaType = int(u32)

	if self.Name, err = up.PopShortStr(); err != nil {
		return err
	}

	if u32, err = up.PopUint32(); err != nil {
		return err
	}
	self.GroupID = int32(u32)

	if self.Data, err = up.PopByteSlice(); err != nil {
		return err
	}

	if u64, err = up.PopUint64(); err != nil {
		return err
	}
	self.Timestamp = int64(u64)

	if u32, err = up.PopUint32(); err != nil {
		return err
	}
	self.Statu = int(u32)

	return nil
}

// ProxyInfo 获取的节点信息
type ProxyInfo struct {
	ServerID int64
	Name     string
	GroupID  int32
	Statu    int            // S2sMetaStatus
	IPList   map[int]net.IP // common.ISP -> net.IP
	Port     int
	Property map[string]string
}

func (p *ProxyInfo) IP() net.IP {
	if len(p.IPList) == 0 {
		return net.IP{}
	}

	isp := math.MaxInt32
	for s := range p.IPList {
		if s < isp {
			isp = s
		}
	}
	return p.IPList[isp]
}

func getProxyInfo(meta *S2SMeta) (*ProxyInfo, error) {
	if meta.MetaType != S2S_S2SDECODER {
		return nil, fmt.Errorf("s2s meta type error: %d", meta.MetaType)
	}

	info := ProxyInfo{}
	info.ServerID = meta.ServerID
	info.Name = meta.Name
	info.GroupID = meta.GroupID
	info.Statu = meta.Statu
	if err := decodeRegInfo(meta.Data, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

// S2S 默认封装
type reginfo struct {
	IPList      []int64  `bson:"iplist"`
	TCPPort     int      `bson:"tcp_port"`
	ExPropKey   []string `bson:"exPropKey"`
	ExPropValue []string `bson:"exPropValue"`
}

func encodeRegInfo(iplist map[int]net.IP, port int, property map[string]string) []byte {
	reg := reginfo{make([]int64, 0), 0, make([]string, 0), make([]string, 0)}
	for k, v := range iplist {
		var ispip int64
		ispip = (int64(k) << 32) | int64(util.InetAton(v))
		reg.IPList = append(reg.IPList, ispip)
	}
	reg.TCPPort = port
	for k, v := range property {
		reg.ExPropKey = append(reg.ExPropKey, k)
		reg.ExPropValue = append(reg.ExPropValue, v)
	}

	bsoninfo, err := bson.Marshal(reg)
	if err != nil {
		return nil
	} else {
		return bsoninfo
	}
}

// 设置ProxyInfo的IPList, Port, Property
func decodeRegInfo(bsondata []byte, info *ProxyInfo) error {
	reg := reginfo{make([]int64, 0), 0, make([]string, 0), make([]string, 0)}

	var err error
	if err = bson.Unmarshal(bsondata, &reg); err != nil {
		return err
	}

	info.IPList = make(map[int]net.IP)
	for _, val := range reg.IPList {
		isp := int(val >> 32)
		ip := uint32(val)
		info.IPList[isp] = util.InetNtoa(ip)
	}
	info.Port = reg.TCPPort
	info.Property = make(map[string]string)
	for idx, val := range reg.ExPropKey {
		info.Property[val] = reg.ExPropValue[idx]
	}
	return nil
}
