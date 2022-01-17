package main

import (
	"goBase/annego/config"
	"goBase/annego/logger"
	"goBase/annego/s2s"
)

func s2sHandle(proxy *s2s.ProxyInfo) {
	logger.Info("serverId %x name %s groupId %d statu %d prot %d",
		proxy.ServerID, proxy.Name, proxy.GroupID, proxy.Statu, proxy.Port)
	for isp, ip := range proxy.IPList {
		logger.Info("isp %d ip %s", isp, ip.String())
	}
	for key, val := range proxy.Property {
		logger.Info("key %s val %s", key, val)
	}
}

func main() {
	if err := config.InitDefaultHostInfo(); err != nil {
		logger.Warning("load hostinfo error: %v", err)
		return
	}
	if !s2s.Init("zhy_test1", "3496c9ae9c838ec6db8d2c53bb736a2f8cd693a58551cfa3bb166ebca51780f5") {
		logger.Warning("s2s init fail")
		return
	}
	logger.Info("s2s test init")
	s2sch := s2s.Subscribe("zhy_test2")
	property := make(map[string]string)
	property["name"] = "testname"
	property["addr"] = "0"
	if err := s2s.Config(config.DefaultHostInfo.IPList, 1245, property); err != nil {
		logger.Warning("s2s config error %v", err)
		return
	}
	go s2s.Start()

	for {
		proxy := <-s2sch
		s2sHandle(proxy)
	}
}
