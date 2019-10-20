package cmdrun

import (
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"

	"github.com/JabinGP/mdout/dir"
	"github.com/JabinGP/mdout/parse"
	"github.com/JabinGP/mdout/types"
)

// URLToPdfE 输入url，输出pdf
func URLToPdfE(inSource string, cmdParmas *types.CmdParmas) error {

	// 路径绝对化
	absOut, err := filepath.Abs(cmdParmas.Out)
	if err != nil {
		log.Println(err)
		return err
	}

	// 构建输出路径参数
	var outDir, outName, outExt string

	if cmdParmas.Out == "" { // 如果没有输出路径，默认输出到源文件所在文件夹
		outDir, err = filepath.Abs("./")
		if err != nil {
			log.Println(err)
			return err
		}
		outName = url.QueryEscape(inSource)
		outExt = ".pdf"
	} else { // 如果有输出路径，分为文件夹和文件两种情况
		isDir, err := dir.IsDir(absOut)
		if err != nil {
			log.Println(err)
			return err
		}
		if isDir { // 是文件夹
			outName = url.QueryEscape(inSource)
			outDir = absOut
			outExt = ".pdf"
		} else { // 不是文件夹，认为是一个文件路径
			outDir, outName, outExt, err = dir.GetDirNameExt(absOut)
		}
	}
	// pdf文件输出路径
	outFinal := outDir + "/" + outName + outExt

	// 将html文件转换成pdf byte
	pdfBts, err := parse.Print(inSource, cmdParmas.PageFormat, cmdParmas.PageOrientation, cmdParmas.PageMargin)
	if err != nil {
		log.Println(err)
		return err
	}

	// 将得到的pdf byte写入文件
	err = ioutil.WriteFile(outFinal, *pdfBts, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	// 输出成功信息
	log.Println("成功保存pdf文件在：" + outFinal)
	return nil
}
