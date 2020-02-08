package cmdrun

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/JabinGP/mdout/parse"
	"github.com/JabinGP/mdout/tool"
	"github.com/JabinGP/mdout/types"
)

// MdToTag 输入md，输出md解析后标签
func MdToTag(in string, parmas types.Parmas) error {
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

	outFullName, err := tool.GetOutFullName(in, parmas)
	if err != nil {
		log.Println("获取输出路径失败！", err)
		return err
	}

	// 将得到的tag写入文件
	err = ioutil.WriteFile(outFullName, *mdBts, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	// 输出成功信息
	log.Println("成功保存markdown解析标签在：" + outFullName)
	return nil
}
