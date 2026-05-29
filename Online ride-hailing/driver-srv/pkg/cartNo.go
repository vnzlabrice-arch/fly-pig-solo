package pkg

import (
	"errors"
	"strings"
)

// LicenseOCRResult 驾驶证识别结果
type LicenseOCRResult struct {
	CarNo   string // 车牌号
	CarType string // 车型
	IsValid bool   // 是否有效（真假）
}

// OcrLicense 模拟行驶证/驾驶证OCR识别 + 真假验证
// imgUrl: 图片URL, carPlate: 用户输入的车牌号（用于比对）
func OcrLicense(imgUrl string, carPlate string) (*LicenseOCRResult, error) {
	// ======================
	// 模拟规则（演示用）
	// ======================
	// 1. 图片地址包含 "fake"/"invalid"/"jia" → 判定为假证
	// 2. 否则默认识别为真，返回用户输入的车牌号
	// ======================

	lowerUrl := strings.ToLower(imgUrl)

	// 模拟：假证
	if strings.Contains(lowerUrl, "fake") ||
		strings.Contains(lowerUrl, "invalid") ||
		strings.Contains(lowerUrl, "jia") {
		return &LicenseOCRResult{
			CarNo:   "",
			CarType: "",
			IsValid: false,
		}, errors.New("证件伪造/无效，请重新上传")
	}

	// 模拟：真证（返回用户输入的车牌号作为OCR识别结果）
	return &LicenseOCRResult{
		CarNo:   carPlate, // 使用用户输入的车牌号
		CarType: "模拟车型",
		IsValid: true,
	}, nil
}
