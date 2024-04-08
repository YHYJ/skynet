/*
File: define_filter.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-31 16:45:30

Description: 过滤器
*/

package general

import (
	"net"
	"strings"
)

// IsDockerInterface 通过接口名称前缀判断是否是 Docker 虚拟接口
//
// 参数：
//   - iface: 接口信息
//
// 返回：
//   - 是否是 Docker 虚拟接口
func IsDockerInterface(iface net.Interface) bool {
	ifaceName := strings.ToLower(iface.Name)
	if strings.HasPrefix(ifaceName, "br-") || strings.HasPrefix(ifaceName, "veth") || strings.HasPrefix(ifaceName, "docker") {
		return true
	}
	return false
}
