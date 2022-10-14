//go:build windows

package main

import (
	"github.com/Trisia/gosysproxy"
	"golang.org/x/sys/windows/registry"
)

func ReadSystemStatus() ([]ProxyState, error) {
	var proxyEnableRaw uint64
	var proxyServerRaw string
	proxyReg, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.QUERY_VALUE)
	if err == registry.ErrNotExist {
		//This system has never enabled system proxy.
		return []ProxyState{{
			Enabled: false,
		}}, nil
	} else if err != nil && err != registry.ErrNotExist {
		return nil, err
	}
	defer proxyReg.Close()
	proxyEnableRaw, _, err = proxyReg.GetIntegerValue("ProxyEnable")
	if err != nil {
		return nil, err
	}
	proxyServerRaw, _, err = proxyReg.GetStringValue("ProxyServer")
	if err != nil {
		return nil, err
	}

	result := make([]ProxyState, 1)
	result[0] = ProxyState{
		Enabled: proxyEnableRaw == 1,
		Server:  proxyServerRaw,
		Device:  "System",
		Type:    "",
	}
	return result, nil
}

func ApplyProxyStatus(states []ProxyState) error {
	if states[0].Enabled {
		err := gosysproxy.SetGlobalProxy(states[0].Server)
		if err != nil {
			return err
		}
	} else {
		err := gosysproxy.Off()
		if err != nil {
			return err
		}
	}
	return nil
}

func SetSystemProxy(url string) error {
	return gosysproxy.SetGlobalProxy(url)
}
