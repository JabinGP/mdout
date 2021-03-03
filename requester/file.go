package requester

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/tool"
)

func buildFileReq(inPath string, params model.Params) (*Request, error) {
	// 获取输入文件类型
	inExt := filepath.Ext(filepath.Base(inPath))
	inType := strings.ReplaceAll(inExt, ".", "")

	// 检查输入类型
	if !tool.CheckType(inType, []string{"md", "tag", "html", "htm"}) {
		return nil, errors.New("非法的输入文件后缀类型：" + inType)
	}

	// 路径绝对化
	absInPath, err := tool.Abs(inPath)
	if err != nil {
		return nil, err
	}
	absOutPath, err := tool.GetOutFullName(inPath, params)
	if err != nil {
		return nil, err
	}

	// 检验文件有效性
	if !tool.IsExists(absInPath) {
		return nil, errors.New("非法的输入文件，文件 " + absInPath + " 不存在！")
	}

	req := &Request{
		Params:     params,
		InPath:     inPath,
		AbsInPath:  absInPath,
		AbsOutPath: absOutPath,
	}

	switch inType {
	case "md":
		err = buildMDReq(req)
	case "tag":
		err = buildTagReq(req)
	default:
		err = buildHTMLReq(req)
	}

	if err != nil {
		return nil, err
	}

	return req, nil
}
