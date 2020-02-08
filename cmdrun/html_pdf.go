package cmdrun

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/JabinGP/mdout/parse"
	"github.com/JabinGP/mdout/tool"
	"github.com/JabinGP/mdout/types"
)

// HTMLToPdf 输入html，输入pdf
func HTMLToPdf(in string, parmas types.Parmas) error {
	// 路径绝对化
	absIn, err := filepath.Abs(in)
	if err != nil {
		log.Println(err)
		return err
	}

	outFileName, err := tool.GetOutFullName(in, parmas)
	if err != nil {
		log.Println("获取输出路径失败！", err)
		return err
	}

	// 路径符合chrome要求，替换 # 为 %23
	chromeURI := strings.Replace(absIn, "#", "%23", -1)

	// 将html文件转换成pdf byte
	pdfBts, err := parse.Print("file://"+chromeURI, parmas.PageFormat, parmas.PageOrientation, parmas.PageMargin)
	if err != nil {
		log.Println(err)
		return err
	}

	// 将得到的pdf byte写入文件
	err = ioutil.WriteFile(outFileName, *pdfBts, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	// 输出成功信息
	log.Println("成功保存pdf文件在：" + outFileName)
	return nil
}
