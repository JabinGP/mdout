package cmdrun

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/JabinGP/mdout/dir"
	"github.com/JabinGP/mdout/parse"
	"github.com/JabinGP/mdout/types"
)

// HTMLToPdfE 输入html，输入pdf
func HTMLToPdfE(inSource string, cmdParmas *types.CmdParmas) error {
	// 检验输入合法性
	if !dir.IsExists(inSource) { //输入文件合法性
		return errors.New("非法的输入文件，文件" + inSource + "不存在")
	}

	// 路径绝对化
	absInSource, err := filepath.Abs(inSource)
	if err != nil {
		log.Println(err)
		return err
	}
	absOut, err := filepath.Abs(cmdParmas.Out)
	if err != nil {
		log.Println(err)
		return err
	}

	// 构建输出路径参数
	var outDir, outName, outExt string

	if cmdParmas.Out == "" { // 如果没有输出路径，默认输出到源文件所在文件夹
		outDir, outName, _, err = dir.GetDirNameExt(absInSource)
		if err != nil {
			log.Println(err)
			return err
		}
		outExt = ".pdf"
	} else { // 如果有输出路径，分为文件夹和文件两种情况
		isDir, err := dir.IsDir(absOut)
		if err != nil {
			log.Println(err)
			return err
		}
		if isDir { // 是文件夹
			_, outName, _, err = dir.GetDirNameExt(absInSource) //只取名字
			outDir = absOut
			outExt = ".pdf"
		} else { // 不是文件夹，认为是一个文件路径
			outDir, outName, outExt, err = dir.GetDirNameExt(absOut)
		}
	}
	// pdf文件输出路径
	outFinal := outDir + "/" + outName + outExt

	// 路径符合chrome要求，替换 # 为 %23
	chromeInSource := strings.Replace(absInSource, "#", "%23", -1)

	// 将html文件转换成pdf byte
	pdfBts, err := parse.Print("file://"+chromeInSource, cmdParmas.PageFormat, cmdParmas.PageOrientation, cmdParmas.PageMargin)
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
