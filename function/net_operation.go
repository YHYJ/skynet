/*
File: net_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-26 12:52:23

Description: 网络操作
*/

package function

import "net"

func GetNetInterfaces() (map[int]map[string]string, error) {
	netInterfacesInfo, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	netInterfacesData := make(map[int]map[string]string)
	count := 2

	// 手动添加无法自动获取的0.0.0.0
	netInterfacesData[1] = map[string]string{
		"name": "any",
		"ip":   "0.0.0.0",
	}

	for _, netInterfaceInfo := range netInterfacesInfo {
		addrs, err := netInterfaceInfo.Addrs()
		if err != nil {
			println(err)
			continue
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ok && ipnet.IP.To4() != nil {
				netInterfacesData[count] = map[string]string{
					"name": netInterfaceInfo.Name,
					"ip":   ipnet.IP.String(),
				}
			} else if !ok && ipnet.IP.To4() == nil {
				// 不符合if条件的网卡虽然不显示也会使count加1，所以这里要减1
				count -= 1
			}
		}
		count += 1
	}
	return netInterfacesData, err
}
