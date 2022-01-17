package util

import (
	"net"
	"testing"
)

func TestOrder(t *testing.T) {
	var u32 uint32 = 0x01020304
	if Ltonl(u32) != 0x04030201 {
		t.Error("Ltonl error")
	}
	if Ntoll(u32) != 0x04030201 {
		t.Error("Ntoll error")
	}

	var u16 uint16 = 0x0a0b
	if Ltons(u16) != 0x0b0a {
		t.Error("Ltons error")
	}
	if Ntols(u16) != 0x0b0a {
		t.Error("Ntols error")
	}
}

func TestNtoa(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	ips := InetNtoa(16777343)
	if !ip.Equal(ips) {
		t.Error("ip != ips", ips)
	}

	ip6 := net.ParseIP("2001:db8::68")
	if InetAton(ip6) != 0 {
		t.Error("aton ipvt error")
	}

	ipint := InetAton(ip)
	if ipint != 16777343 {
		t.Error("ipint error:", ip, ipint)
	}
}

func TestSton(t *testing.T) {
	if InetSton("127.0.0.1") != 16777343 {
		t.Error("ston right error")
	}
	if InetSton("abcde") != 0 {
		t.Error("ston abcde error")
	}
	if InetSton("2001:db8::68") != 0 {
		t.Error("ston ipv6 error")
	}
	if InetNtos(16777343) != "127.0.0.1" {
		t.Error("ntos error")
	}
}
