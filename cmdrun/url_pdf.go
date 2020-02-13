package cmdrun

import (
	"io/ioutil"
	"net/url"

	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/parse"
	"github.com/JabinGP/mdout/tool"
)

// URLToPdf 输入url，输出pdf
func URLToPdf(in string, parmas model.Parmas) error {
	escapedIn := url.QueryEscape(in)
	outFileName, err := tool.GetOutFullName(escapedIn, parmas)
	if err != nil {
		return err
	}

	// 将html文件转换成pdf byte
	pdfBts, err := parse.Print(parmas.ExecPath, in, parmas.PageFormat, parmas.PageOrientation, parmas.PageMargin)
	if err != nil {
		return err
	}

	// 将得到的pdf byte写入文件
	err = ioutil.WriteFile(outFileName, *pdfBts, 0644)
	if err != nil {
		return err
	}

	// 输出成功信息
	log.Infof("成功保存pdf文件在：%v", outFileName)
	return nil
}
