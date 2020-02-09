package cmdrun

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/JabinGP/mdout/tool"
	"github.com/JabinGP/mdout/model"
)

// Distribute 根据不同的输入类型处理，定位到不同的执行函数
func Distribute(in string, parmas model.Parmas) error {
	// 获取输入参数类型
	inType, err := tool.GetType(in)
	if err != nil {
		log.Println(err)
		return err
	}

	// 根据不同的输入类型处理，定位到不同的执行函数
	switch inType {

	case "url":
		return inURL(in, parmas)

	case "file":
		return inFile(in, parmas)

	default:
		return inDefault(in, parmas)
	}
}

// inDefault 不符合已知归类时
func inDefault(in string, parmas model.Parmas) error {
	return errors.New("未能正确识别输入，无法找到匹配的预设输入类型，请检查输入。")
}

// inURL 输入为URL时
func inURL(in string, parmas model.Parmas) error {
	return URLToPdf(in, parmas)
}

// inFile 输入为一个文件
func inFile(in string, parmas model.Parmas) error {
	if !tool.IsExists(in) { //输入文件合法性
		return errors.New("非法的输入文件，文件 " + in + " 不存在")
	}

	inExt := filepath.Ext(filepath.Base(in)) // 获取文件后缀

	// 输入文件类型分类
	switch inExt {

	case ".html":
		return inFileHTML(in, parmas)

	case ".md":
		return inFileMd(in, parmas)

	default:
		return inFileDefault(in, parmas)
	}
}
