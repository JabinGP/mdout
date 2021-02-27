package requester

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/tool"
)

type Request struct {
	model.Parmas
	InType     string
	InPath     string
	AbsInPath  string
	AbsOutPath string
	Data       interface{}
}

func NewRequest(inPath string, parmas model.Parmas) (*Request, error) {
	// Judge input type
	inType, err := tool.GetType(inPath)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	if inType == "url" {
		return buildURLReq(inPath, parmas)
	}

	return buildFileReq(inPath, parmas)
}

func buildURLReq(inPath string, parmas model.Parmas) (*Request, error) {
	escapedURL := url.QueryEscape(inPath)

	absOutPath, err := tool.GetOutFullName(inPath, parmas)
	if err != nil {
		return nil, err
	}
	var req = Request{
		Parmas:    parmas,
		InType:    "url",
		InPath:    inPath,
		AbsInPath: escapedURL,
	}
	req.OutType = "pdf" // Must be pdf when input type is url
	req.AbsOutPath = absOutPath

	return &req, nil
}

func buildFileReq(inPath string, parmas model.Parmas) (*Request, error) {
	// 路径绝对化
	absInPath, err := tool.Abs(inPath)
	if err != nil {
		return nil, err
	}
	absOutPath, err := tool.GetOutFullName(inPath, parmas)
	if err != nil {
		return nil, err
	}

	// 检验文件有效性
	if !tool.IsExists(absInPath) {
		return nil, errors.New("非法的输入文件，文件 " + absInPath + " 不存在！")
	}

	// 获取输入文件类型
	inExt := filepath.Ext(filepath.Base(inPath))
	inType := strings.ReplaceAll(inExt, ".", "")
	return &Request{
		Parmas:     parmas,
		InPath:     inPath,
		AbsInPath:  absInPath,
		AbsOutPath: absOutPath,
		InType:     inType,
	}, nil
}
