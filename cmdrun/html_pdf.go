package cmdrun

import (
	"io/ioutil"
	"strings"

	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/parse"
	"github.com/JabinGP/mdout/tool"
)

// HTMLToPdf 输入html，输入pdf
func HTMLToPdf(in string, parmas model.Parmas) error {
	// 路径绝对化
	absIn, err := tool.Abs(in)
	if err != nil {
		return err
	}

	outFileName, err := tool.GetOutFullName(in, parmas)
	if err != nil {
		return err
	}

	// 路径符合chrome要求，替换 # 为 %23
	chromeURI := strings.Replace(absIn, "#", "%23", -1)

	// 将html文件转换成pdf byte
	pdfBts, err := parse.Print(parmas.ExecPath, "file://"+chromeURI, parmas.PageFormat, parmas.PageOrientation, parmas.PageMargin)
	if err != nil {
		return err
	}

	// 将得到的pdf byte写入文件
	err = ioutil.WriteFile(outFileName, *pdfBts, 0644)
	if err != nil {
		return err
	}

	// 输出成功信息
	log.Infof("成功保存pdf文件在：%s", outFileName)
	return nil
}
