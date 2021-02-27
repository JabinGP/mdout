package parse

import (
	"errors"
	"strconv"
	"strings"

	"github.com/JabinGP/mdout/log"
)

// 通过格式获取打印纸张的宽高
func getPaperWidthAndHeight(pageFormat string) (float64, float64) {
	var paperWidth, paperHeight float64
	switch strings.ToLower(pageFormat) {
	case "legal":
		paperWidth, paperHeight = splitWidthXHeight("216x356")
	case "letter":
		paperWidth, paperHeight = splitWidthXHeight("216x279")
	case "tabloid":
		paperWidth, paperHeight = splitWidthXHeight("279x356")
	case "ledger":
		paperWidth, paperHeight = splitWidthXHeight("279x356")
	case "a5":
		paperWidth, paperHeight = splitWidthXHeight("148x210")
	case "a4":
		paperWidth, paperHeight = splitWidthXHeight("210x297")
	case "a3":
		paperWidth, paperHeight = splitWidthXHeight("297x420")
	case "a2":
		paperWidth, paperHeight = splitWidthXHeight("420x594")
	case "a1":
		paperWidth, paperHeight = splitWidthXHeight("594x841")
	case "a0":
		paperWidth, paperHeight = splitWidthXHeight("841x1189")
	default:
		paperWidth, paperHeight = splitWidthXHeight("841x1189")
	}
	return paperWidth / 25.4, paperHeight / 25.4
}

func splitWidthXHeight(widthXHeight string) (float64, float64) {
	numberArr := strings.Split(strings.ToLower(widthXHeight), "x")
	width, err := strconv.ParseFloat(numberArr[0], 32)
	height, err := strconv.ParseFloat(numberArr[1], 32)
	if err != nil {
		log.Errorln(err)
		log.Errorln("建立" + widthXHeight + "纸张尺寸表时类型转换失败，将尺寸定位A4纸210cm*297cm大小")
		return 210, 297
	}
	return width, height
}

func getPaperOrientation(pageOrientation string) bool {
	if strings.ToLower(pageOrientation) == "landscape" {
		return true
	}
	return false
}

func getMargin(pageMargin string) ([4]float64, error) {

	// 替换所有中文逗号为英文逗号
	pageMargin = strings.Replace(pageMargin, "，", ",", -1)

	// 去除空格
	pageMargin = strings.Replace(pageMargin, " ", "", -1)

	// 获取配置数组
	marginArr := strings.Split(pageMargin, ",")

	// 判断输入类型
	switch len(marginArr) {
	case 0:
		return [4]float64{1, 1, 1, 1}, errors.New("无效的输入边距值！")
	case 1:
		marginAll, err := strconv.ParseFloat(marginArr[0], 64)
		if err != nil {
			return [4]float64{1, 1, 1, 1}, err
		}

		// 输出0会被chromedp重设为0.4默认值
		if marginAll == 0 {
			marginAll = 0.0000000000000001
		}
		return [4]float64{marginAll, marginAll, marginAll, marginAll}, nil
	default:
		return [4]float64{1, 1, 1, 1}, errors.New("无法识别的输入边距类型！")
	}
}
