package cmdrun

import (
	"errors"

	"github.com/JabinGP/mdout/model"
)

func inFileDefault(in string, parmas model.Parmas) error {
	return errors.New("识别到文件，但未找到匹配的预设输入文件类型，请检查输入。")
}

func inFileHTML(in string, parmas model.Parmas) error {
	return HTMLToPdf(in, parmas)
}

func inFileMd(in string, parmas model.Parmas) error {
	// 输出文件类型分类
	switch parmas.Type {

	case "tag":
		return MdToTag(in, parmas)

	case "html":
		return MdToHTML(in, parmas)

	case "pdf":
		return MdToPdf(in, parmas)

	default:
		return errors.New("识别到md输入，但未能正确识别匹配的指定输出文件类型，请检查输入。")
	}
}
