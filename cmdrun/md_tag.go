package cmdrun

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/JabinGP/mdout/dir"
	"github.com/JabinGP/mdout/parse"
	"github.com/JabinGP/mdout/types"
)

// MdToTagE 输入md，输出md解析后标签
func MdToTagE(inSource string, cmdParmas *types.CmdParmas) error {
	// 检验输入合法性
	if !dir.IsExists(inSource) {
		return errors.New("非法的输入文件，文件不存在")
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

	// 构建输出路径参数
	var outDir, outName, outExt string

	if cmdParmas.Out == "" { // 如果没有输出路径，默认输出到源文件所在文件夹
		outDir, outName, _, err = dir.GetDirNameExt(absInSource)
		if err != nil {
			log.Println(err)
			return err
		}
		outExt = ".html"
	} else { // 如果有输出路径，分为文件夹和文件两种情况
		isDir, err := dir.IsDir(absOut)
		if err != nil {
			log.Println(err)
			return err
		}
		if isDir { // 是文件夹
			_, outName, _, err = dir.GetDirNameExt(absInSource) //只取名字
			outDir = absOut
			outExt = ".html"
		} else { // 不是文件夹，认为是一个文件路径
			outDir, outName, outExt, err = dir.GetDirNameExt(absOut)
		}
	}
	// 目标html标签文件路径
	outFinal := outDir + "/" + outName + outExt

	// 将得到的tag写入文件
	err = ioutil.WriteFile(outFinal, *mdBts, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	// 输出成功信息
	log.Println("成功保存markdown解析标签在：" + outFinal)
	return nil
}
