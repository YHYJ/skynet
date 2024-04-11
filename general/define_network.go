/*
File: define_network.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-08 13:09:08

Description: 网络操作
*/

package general

import (
	"net"

	"github.com/gookit/color"
)

// GetNetInterfacesForCli 为 CLI 获取网卡信息
//
// 返回：
//   - 网卡信息
//   - 错误信息
func GetNetInterfacesForCli() (map[int]map[string]string, error) {
	netInterfacesInfo, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	netInterfacesData := make(map[int]map[string]string)
	// 手动添加无法自动获取的0.0.0.0
	netInterfacesData[1] = map[string]string{
		"name": "any",
		"ip":   "0.0.0.0",
	}
	count := 1 // 网卡编号

	for _, netInterfaceInfo := range netInterfacesInfo {
		addrs, _ := netInterfaceInfo.Addrs()

		if netInterfaceInfo.Flags&net.FlagUp != 0 {
			for _, addr := range addrs {
				ipnet, ok := addr.(*net.IPNet)
				if ok && ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() && !IsDockerInterface(netInterfaceInfo) {
					count += 1
					netInterfacesData[count] = map[string]string{
						"name": netInterfaceInfo.Name,
						"ip":   ipnet.IP.String(),
					}
				}
			}
		}
	}
	return netInterfacesData, nil
}

// GetNetInterfacesForGui 为 GUI 获取网卡信息
//
// 返回：
//   - 网卡信息
//   - 错误信息
func GetNetInterfacesForGui() ([]string, error) {
	netInterfacesInfo, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// 手动添加无法自动获取的0.0.0.0
	netInterfacesData := []string{DefaultNic}

	for _, netInterfaceInfo := range netInterfacesInfo {
		addrs, _ := netInterfaceInfo.Addrs()

		if netInterfaceInfo.Flags&net.FlagUp != 0 {
			for _, addr := range addrs {
				ipnet, ok := addr.(*net.IPNet)
				if ok && ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() && !IsDockerInterface(netInterfaceInfo) {
					OtherNic = color.Sprintf("%s - %s", netInterfaceInfo.Name, ipnet.IP.String())
					netInterfacesData = append(netInterfacesData, OtherNic)
				}
			}
		}
	}
	return netInterfacesData, nil
}
