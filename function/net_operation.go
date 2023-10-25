/*
File: net_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-26 12:52:23

Description: 网络操作
*/

package function

import (
	"net"
	"strings"
)

func GetNetInterfaces() (map[int]map[string]string, error) {
	netInterfacesInfo, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	netInterfacesData := make(map[int]map[string]string)
	// 手动添加无法自动获取的0.0.0.0
	netInterfacesData[1] = map[string]string{
		"name": "any",
		"ip":   "0.0.0.0",
	}
	count := 1 // 网卡编号

	for _, netInterfaceInfo := range netInterfacesInfo {
		addrs, err := netInterfaceInfo.Addrs()
		if err != nil {
			println(err)
			continue
		}

		if netInterfaceInfo.Flags&net.FlagUp != 0 {
			for _, addr := range addrs {
				ipnet, ok := addr.(*net.IPNet)
				if ok && ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() && !isDockerInterface(netInterfaceInfo) {
					count += 1
					netInterfacesData[count] = map[string]string{
						"name": netInterfaceInfo.Name,
						"ip":   ipnet.IP.String(),
					}
				}
			}
		}
	}
	return netInterfacesData, err
}

// 通过接口名称前缀判断是否是虚拟接口
func isDockerInterface(iface net.Interface) bool {
	ifaceName := strings.ToLower(iface.Name)
	if strings.HasPrefix(ifaceName, "br-") || strings.HasPrefix(ifaceName, "veth") || strings.HasPrefix(ifaceName, "docker") {
		return true
	}
	return false
}
