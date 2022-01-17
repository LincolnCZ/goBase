package util

import (
	"encoding/binary"
	"net"
)

// Ltonl uint32 local order to network order
func Ltonl(u32 uint32) uint32 {
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], u32)
	return binary.BigEndian.Uint32(bytes[:])
}

// Ntoll uint32 network order to local order
func Ntoll(u32 uint32) uint32 {
	var bytes [4]byte
	binary.BigEndian.PutUint32(bytes[:], u32)
	return binary.LittleEndian.Uint32(bytes[:])
}

// Ltons uint16 local order to network order
func Ltons(u16 uint16) uint16 {
	var bytes [2]byte
	binary.LittleEndian.PutUint16(bytes[:], u16)
	return binary.BigEndian.Uint16(bytes[:])
}

// Ntols uint16 network order to local order
func Ntols(u16 uint16) uint16 {
	var bytes [2]byte
	binary.BigEndian.PutUint16(bytes[:], u16)
	return binary.LittleEndian.Uint16(bytes[:])
}

// InetNtoa uint32 network order ip to net.IP
func InetNtoa(nip uint32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(nip >> 24 & 0xFF)
	bytes[1] = byte(nip >> 16 & 0xFF)
	bytes[2] = byte(nip >> 8 & 0xFF)
	bytes[3] = byte(nip & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// InetAton net.IP to uint32 network order ip
// if ip not ipv4 return 0
func InetAton(ip net.IP) uint32 {
	ip = ip.To4()
	if len(ip) == 0 {
		return 0
	}

	var nip uint32
	nip |= uint32(ip[3]) << 24
	nip |= uint32(ip[2]) << 16
	nip |= uint32(ip[1]) << 8
	nip |= uint32(ip[0])
	return nip
}

// InetNtos uint32 network order ip to ipv4 string
func InetNtos(nip uint32) string {
	ip := InetNtoa(nip)
	return ip.String()
}

// InetSton string ipv4 to uint32 network order ip
// if ip not ipv4 return 0
func InetSton(s string) uint32 {
	ip := net.ParseIP(s).To4()
	if len(ip) == 0 {
		return 0
	}

	var nip uint32
	nip |= uint32(ip[3]) << 24
	nip |= uint32(ip[2]) << 16
	nip |= uint32(ip[1]) << 8
	nip |= uint32(ip[0])
	return nip
}
