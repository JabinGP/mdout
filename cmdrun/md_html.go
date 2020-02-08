package cmdrun

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/JabinGP/mdout/parse"
	"github.com/JabinGP/mdout/tool"
	"github.com/JabinGP/mdout/types"
)

// MdToHTML 输入md，输出完整html页面
func MdToHTML(in string, parmas types.Parmas) error {
	// 路径绝对化
	absIn, err := filepath.Abs(in)
	if err != nil {
		log.Println(err)
		return err
	}

	// 读取源文件
	sourceBts, err := ioutil.ReadFile(absIn)
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

	// 拼接输出文件全名
	outFullName, err := tool.GetOutFullName(in, parmas)
	if err != nil {
		log.Println("获取输出路径失败！", err)
		return err
	}
	// 将得到的tag写入文件
	err = ioutil.WriteFile(outFullName, *htmlBts, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	// 输出成功信息
	log.Println("成功保存html文件在：" + outFullName)

	return nil
}
