/*
File: i18n.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-23 09:52:36

Description: 翻译
*/

package function

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// 创建一个语言解析器
func CreateLocalizer() *i18n.Localizer {
	// 创建一个 Bundle 以在应用程序的整个生命周期中使用
	bundle := i18n.NewBundle(language.English)
	// 在初始化时，将翻译文件加载到 Bundle 中
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile("locales/active.en.toml")
	bundle.LoadMessageFile("locales/active.zh_CN.toml")
	// 定义语言顺序
	langs := strings.Split(GetVariable("LANG"), ".")[0]
	langs += ",en_US"
	// 创建一个 Localizer 以便解析首选语言
	localizer := i18n.NewLocalizer(bundle, langs)

	return localizer
}

func Translate(localizer *i18n.Localizer, id string, templateData map[string]interface{}) string {
	result := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: templateData,
	},
	)

	return result
}
