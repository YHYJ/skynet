/*
File: qr_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-26 09:25:11

Description: 二维码操作
*/

package function

import (
	"image"

	"github.com/skip2/go-qrcode"
)

// 生成二维码图片
func QrCodeImage(content string) (image.Image, error) {
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	qr.DisableBorder = false
	return qr.Image(256), nil
}

// 生成二维码字符串
func QrCodeString(content string) (string, error) {
	qt, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return "", err
	}
	qrString := qt.ToSmallString(false)

	return qrString, nil
}
