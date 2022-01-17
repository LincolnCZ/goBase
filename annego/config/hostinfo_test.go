package config

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostInfo(t *testing.T) {
	config := LoadHostInfoConfig("hostinfo.ini")
	iplist := map[int]net.IP{
		CTL:      net.IPv4(221, 228, 107, 123),
		CNC:      net.IPv4(103, 229, 149, 122),
		MOB:      net.IPv4(112, 25, 251, 123),
		INTRANET: net.IPv4(10, 25, 66, 99),
	}

	assert.NotNil(t, config)
	assert.Equal(t, config.ServerId, 73634)
	assert.Equal(t, config.CityId, 1296)
	assert.Equal(t, config.Area, 3)
	assert.Equal(t, config.PriGroupId, 538)
	assert.Equal(t, config.SecGroupId, 538)
	assert.Equal(t, config.GetISP(), CTL)
	assert.Equal(t, config.GetIP(), net.IPv4(221, 228, 107, 123))
	assert.Equal(t, config.IPList, iplist)
}

func TestSelectBestIP(t *testing.T) {
	config := LoadHostInfoConfig("hostinfo.ini")

	iplist := map[int]net.IP{
		CTL:      net.IPv4(221, 228, 107, 123),
		CNC:      net.IPv4(103, 229, 149, 122),
		MOB:      net.IPv4(112, 25, 251, 123),
		INTRANET: net.IPv4(10, 25, 66, 99),
	}
	assert.True(t, config.SelectBestIP(iplist).Equal(net.IPv4(221, 228, 107, 123)))

	iplist = map[int]net.IP{
		CNC:      net.IPv4(103, 229, 149, 122),
		MOB:      net.IPv4(112, 25, 251, 123),
		INTRANET: net.IPv4(10, 25, 66, 99),
	}
	assert.True(t, config.SelectBestIP(iplist).Equal(net.IPv4(103, 229, 149, 122)))

	iplist = map[int]net.IP{
		ASIA: net.IPv4(103, 229, 149, 122),
		SA:   net.IPv4(10, 25, 66, 99),
	}
	assert.True(t, config.SelectBestIP(iplist).Equal(net.IPv4(103, 229, 149, 122)))
}
