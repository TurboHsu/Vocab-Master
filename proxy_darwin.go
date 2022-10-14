//go:build darwin

package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func GetName(source string) string {
	pIndex := strings.Index(source, ":")
	for i := pIndex + 1; i < len(source); i++ {
		if source[i] != ' ' {
			return strings.TrimSpace(source[i:])
		}
	}
	return ""
}

func IsInterfaceActive(name string) (bool, error) {
	var out bytes.Buffer
	cmd := exec.Command("ifconfig")
	cmd.Stdout = &out
	_ = cmd.Run()
	ifconfigInfo := strings.Split(out.String(), "\n")

	current := false
	for _, ifconfig := range ifconfigInfo {
		if !strings.HasPrefix(ifconfig, "\t") {
			current = strings.HasPrefix(ifconfig, name+":")
		} else if current && strings.Contains(ifconfig, "status:") {
			activity := GetName(ifconfig)
			return activity == "active", nil
		}
	}
	return false, errors.New("interface not listed")
}

func GetActivePorts() []string {
	var result = make([]string, 0)

	cmd := exec.Command("networksetup", "-listallhardwareports")
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run()
	portsInfo := strings.Split(out.String(), "\n")

	portName := ""
	deviceName := ""
	for i := 0; i < len(portsInfo); i++ {
		port := portsInfo[i]
		if strings.Contains(port, "Hardware Port:") {
			portName = GetName(port)
		} else if strings.Contains(port, "Device:") {
			deviceName = GetName(port)
		}
		if portName != "" && deviceName != "" {
			active, err := IsInterfaceActive(deviceName)
			if err != nil {
				continue
			}
			if active {
				result = append(result, portName)
			}
			portName = ""
			deviceName = ""
			i++
		}
	}

	return result
}

func GatherProxyInfo(t string, dev string) (*ProxyState, error) {
	cmd := exec.Command("networksetup", "-get"+t, dev)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return nil, err
	}

	var result ProxyState
	result.Type = t
	result.Device = dev

	proxyInfo := strings.Split(out.String(), "\n")
	for _, info := range proxyInfo {
		ele := strings.Split(info, ":")
		if len(ele) < 2 {
			continue
		}
		name := strings.TrimSpace(ele[0])
		value := strings.TrimSpace(ele[1])
		switch name {
		case "Enabled":
			result.Enabled = value == "Yes"
		case "Server":
			result.Server = value + result.Server
		case "Port":
			result.Server += ":" + value
		}
	}

	return &result, nil
}

func ReadSystemStatus() ([]ProxyState, error) {
	devices := GetActivePorts()
	results := make([]ProxyState, 2*len(devices))

	for _, device := range devices {
		http, err := GatherProxyInfo("webproxy", device)
		if err != nil {
			return nil, err
		}
		results[0] = *http

		https, err := GatherProxyInfo("securewebproxy", device)
		if err != nil {
			return nil, err
		}
		results[1] = *https
	}

	return results, nil
}

func SetTypedProxyState(t string, dev string, enabled bool) error {
	state := "off"
	if enabled {
		state = "on"
	}
	cmd := exec.Command("networksetup", "-set"+t+"state", dev, state)
	return cmd.Run()
}

func SetTypedProxy(t string, dev string, url string) error {
	err := SetTypedProxyState(t, dev, true)
	if err != nil {
		return err
	}
	splits := strings.Split(url, ":")
	host := splits[0]
	port := splits[1]
	cmd := exec.Command("networksetup", "-set"+t, dev, host, port)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ApplyProxyStatus(states []ProxyState) error {
	for _, state := range states {
		err := SetTypedProxy(state.Type, state.Device, state.Server)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetSystemProxy(url string) error {
	devices := GetActivePorts()
	for _, dev := range devices {
		err1 := SetTypedProxy("webproxy", dev, url)
		err2 := SetTypedProxy("securewebproxy", dev, url)
		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
	}
	return nil
}
