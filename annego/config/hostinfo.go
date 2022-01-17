package config

import (
	"fmt"
	"net"
	"sort"
	"strings"

	"goBase/annego/logger"
	"gopkg.in/ini.v1"
)

const DEFAULT_HOSTINFO_PATH = "/home/dspeak/yyms/hostinfo.ini"

type HostInfoConfig struct {
	IPList     map[int]net.IP
	Status     int
	Area       int
	CityId     int
	ServerId   int
	SecGroupId int
	PriGroupId int
}

// GetIP 获取根据ISP排序
func (config *HostInfoConfig) GetISP() int {
	if len(config.IPList) > 0 {
		isp := MAX_ISP
		for k := range config.IPList {
			if k < isp {
				isp = k
			}
		}
		return isp
	}
	return 0
}

// GetIP 获取根据ISP排序第一个IP
func (config *HostInfoConfig) GetIP() net.IP {
	if isp := config.GetISP(); isp != 0 {
		return config.IPList[isp]
	}
	return nil
}

// SelectBestIP 获取ISP->IP列表中最合适的IP
func (config *HostInfoConfig) SelectBestIP(iplist map[int]net.IP) net.IP {
	if len(iplist) == 0 {
		return net.IP{}
	}

	isp := config.GetISP()
	ip, ok := iplist[isp]
	if ok {
		return ip
	} else {
		keys := make([]int, 0, len(iplist))
		for k := range iplist {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		return iplist[keys[0]]
	}
}

// LoadHostInfoConfigWithError
func LoadHostInfoConfigWithError(path string) (*HostInfoConfig, error) {
	fp, err := ini.Load(path)
	if err != nil {
		return nil, err
	}
	config := HostInfoConfig{}
	section := fp.Section("")

	config.Status, _ = section.Key("status").Int()
	config.Area, _ = section.Key("area_id").Int()
	config.CityId, _ = section.Key("city_id").Int()
	config.ServerId, _ = section.Key("server_id").Int()
	config.SecGroupId, _ = section.Key("sec_group_id").Int()
	config.PriGroupId, _ = section.Key("pri_group_id").Int()
	config.IPList = parseHostIPList(section.Key("ip_isp_list").String())
	if len(config.IPList) == 0 {
		return nil, fmt.Errorf("HostInfo IPList empty")
	}
	return &config, nil
}

// LoadHostInfoConfig 失败返回nil
func LoadHostInfoConfig(path string) *HostInfoConfig {
	config, err := LoadHostInfoConfigWithError(path)
	if err != nil {
		logger.Error("load HostInfo %s error: %v", path, err)
	}
	return config
}

var ispStringMap map[string]int = map[string]int{
	"CTL":      CTL,
	"CNC":      CNC,
	"EDU":      EDU,
	"WBN":      WBN,
	"MOB":      MOB,
	"BGP":      BGP,
	"HK":       ASIA,
	"BRA":      SA,
	"EU":       EU,
	"NA":       NA,
	"INTRANET": INTRANET,
}

func parseHostIPList(info string) map[int]net.IP {
	iplist := make(map[int]net.IP)
	ipinfos := strings.Split(info, ",")
	for _, ipinfo := range ipinfos {
		ips := strings.Split(ipinfo, ":")
		if len(ips) < 2 {
			logger.Warning("Load ipinfo error: %s", ipinfo)
			continue
		}

		ip := net.ParseIP(ips[0])
		if ip == nil {
			logger.Warning("parse ipinfo error: %s", ipinfo)
			continue
		}
		if isp, ok := ispStringMap[ips[1]]; ok {
			iplist[isp] = ip
		}
	}
	return iplist
}

var DefaultHostInfo *HostInfoConfig

func InitDefaultHostInfo() error {
	var err error
	if DefaultHostInfo == nil {
		DefaultHostInfo, err = LoadHostInfoConfigWithError(DEFAULT_HOSTINFO_PATH)
	}
	return err
}

func init() {
	InitDefaultHostInfo()
}
