/*
File: define_network.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-26 12:52:23

Description: 网络操作
*/

package gui

import (
	"fmt"
	"net"

	"github.com/yhyj/skynet/general"
)

// GetNetInterfaces 获取网卡信息
//
// 返回：
//   - 网卡信息
//   - 错误信息
func GetNetInterfaces() ([]string, error) {
	netInterfacesInfo, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// 手动添加无法自动获取的0.0.0.0
	netInterfacesData := []string{defaultNic}

	for _, netInterfaceInfo := range netInterfacesInfo {
		addrs, _ := netInterfaceInfo.Addrs()

		if netInterfaceInfo.Flags&net.FlagUp != 0 {
			for _, addr := range addrs {
				ipnet, ok := addr.(*net.IPNet)
				if ok && ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() && !general.IsDockerInterface(netInterfaceInfo) {
					otherNic = fmt.Sprintf("%s - %s", netInterfaceInfo.Name, ipnet.IP.String())
					netInterfacesData = append(netInterfacesData, otherNic)
				}
			}
		}
	}
	return netInterfacesData, nil
}
