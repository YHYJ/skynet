/*
File: variable.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-23 12:33:31

Description: 预定义变量
*/

package function

var (
	TemplateData = map[string]interface{}{} // i18n模板数据
	Localizer    = CreateLocalizer()        // 语言解析器
)
