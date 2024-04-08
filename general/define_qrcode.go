/*
File: define_qrcode.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-26 09:25:11

Description: 操作二维码
*/

package general

import (
	"image"

	"github.com/skip2/go-qrcode"
)

// QrCodeImage 生成二维码图片
//
// 参数：
//   - content: 二维码内容
//
// 返回：
//   - 图像对象
//   - 错误信息
func QrCodeImage(content string) (image.Image, error) {
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	qr.DisableBorder = false
	return qr.Image(256), nil
}

// QrCodeString 生成二维码字符串
//
// 参数：
//   - content: 二维码内容
//
// 返回：
//   - 二维码字符串
//   - 错误信息
func QrCodeString(content string) (string, error) {
	qt, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return "", err
	}
	qrString := qt.ToSmallString(false)

	return qrString, nil
}
