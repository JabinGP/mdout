package cmdrun

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/JabinGP/mdout/parse"
	"github.com/JabinGP/mdout/tool"
	"github.com/JabinGP/mdout/types"
)

// MdToPdf 输入md，输出pdf
func MdToPdf(in string, parmas types.Parmas) error {
	// 路径绝对化
	absInSource, err := filepath.Abs(in)
	if err != nil {
		log.Println(err)
		return err
	}

	// 读取源文件
	sourceBts, err := ioutil.ReadFile(absInSource)
	if err != nil {
		log.Println(err)
		return err
	}
	// md解析
	mdBts, err := parse.Md(sourceBts)
	if err != nil {
		log.Println(err)
		return err
	}
	// tag拼接
	htmlBts, err := parse.AssembleTag(parmas.Theme, mdBts)
	if err != nil {
		log.Println(err)
		return err
	}

	// 构建临时html文件路径
	tmpDir, tmpName, _, err := tool.GetDirNameExt(absInSource)
	tmpExt := "html"
	tmpFullName, err := filepath.Abs(tmpDir + "/" + "tmp." + tmpName + "." + tmpExt)
	if err != nil {
		log.Println("拼接临时html文件全名错误！", err)
		return err
	}
	// 将中间html写入文件
	err = ioutil.WriteFile(tmpFullName, *htmlBts, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	// 清除临时html文件
	defer func() {
		if tool.IsExists(tmpFullName) {
			err := os.Remove(tmpFullName)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	// 路径符合chrome要求，替换 # 为 %23
	chromeTmpURI := strings.Replace(tmpFullName, "#", "%23", -1)

	// 将html文件转换成pdf byte
	pdfBts, err := parse.Print("file://"+chromeTmpURI, parmas.PageFormat, parmas.PageOrientation, parmas.PageMargin)
	log.Println("开始保存文件")

	outFullName, err := tool.GetOutFullName(in, parmas)
	if err != nil {
		log.Println("获取输出路径失败！", err)
		return err
	}

	// 将得到的pdf写入文件
	err = ioutil.WriteFile(outFullName, *pdfBts, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	// 输出成功信息
	log.Println("成功保存pdf文件在：" + outFullName)
	return nil
}
