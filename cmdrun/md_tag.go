package cmdrun

import (
	"io/ioutil"

	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/parse"
	"github.com/JabinGP/mdout/tool"
)

// MdToTag 输入md，输出md解析后标签
func MdToTag(in string, parmas model.Parmas) error {
	// 路径绝对化
	absIn, err := tool.Abs(in)
	if err != nil {
		return err
	}

	// 读取源文件
	sourceBts, err := ioutil.ReadFile(absIn)
	if err != nil {
		return err
	}

	// md解析
	mdBts, err := parse.Md(sourceBts)
	if err != nil {
		return err
	}

	outFullName, err := tool.GetOutFullName(in, parmas)
	if err != nil {
		return err
	}

	// 将得到的tag写入文件
	err = ioutil.WriteFile(outFullName, *mdBts, 0644)
	if err != nil {
		return err
	}

	// 输出成功信息
	log.Infof("成功保存markdown解析标签在：%v", outFullName)
	return nil
}
