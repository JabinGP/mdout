package parse

import (
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/JabinGP/mdout/model"

	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/tool"
)

// FileRequest ...
type FileRequest struct {
	InPath         string
	Parmas         model.Parmas
	AbsInPath      string
	AbsOutPath     string
	InType         string
	FileServer     bool
	FileServerAddr string
}

// NewFileRequest ..
func NewFileRequest(inPath string, parmas model.Parmas) (FileRequest, error) {
	// 路径绝对化
	absInPath, err := tool.Abs(inPath)
	if err != nil {
		return FileRequest{}, err
	}
	absOutPath, err := tool.GetOutFullName(inPath, parmas)
	if err != nil {
		return FileRequest{}, err
	}

	// 检验文件有效性
	if !tool.IsExists(absInPath) {
		return FileRequest{}, errors.New("非法的输入文件，文件 " + absInPath + " 不存在！")
	}

	// 获取输入文件类型
	inExt := filepath.Ext(filepath.Base(inPath))
	inType := strings.ReplaceAll(inExt, ".", "")
	return FileRequest{
		InPath:     inPath,
		Parmas:     parmas,
		AbsInPath:  absInPath,
		AbsOutPath: absOutPath,
		InType:     inType,
	}, nil
}

// SaveAsFile ...
func (r *FileRequest) SaveAsFile(bts *[]byte) error {
	err := ioutil.WriteFile(r.AbsOutPath, *bts, 0644)
	if err != nil {
		return err
	}
	log.Infof("成功保存文件：%s", r.AbsOutPath)
	return nil
}

// GetBts ...
func (r *FileRequest) GetBts() (*[]byte, error) {
	switch r.InType {
	case "md":
		return r.parseMd()
	case "html":
		return r.parseHtml()
	}
	return nil, errors.New("未找到匹配的输入类型：" + r.InType)
}

// ParseMd ...
func (r *FileRequest) parseMd() (*[]byte, error) {
	// 读取源文件
	sourceBts, err := ioutil.ReadFile(r.AbsInPath)
	if err != nil {
		return nil, err
	}

	// md解析
	mdBts, err := Md(sourceBts)
	if err != nil {
		return nil, err
	}
	if r.Parmas.Type == "tag" {
		return mdBts, nil
	}

	// tag拼接
	htmlBts, err := AssembleTag(r.Parmas.Theme, mdBts)
	if err != nil {
		return nil, err
	}
	if r.Parmas.Type == "html" {
		return htmlBts, nil
	}

	// 构建临时html文件路径
	tmpDir, tmpName, _, err := tool.GetDirNameExt(r.AbsInPath)
	tmpExt := "html"
	tmpFullName, err := tool.Abs(tmpDir + "/" + "tmp." + tmpName + "." + tmpExt)
	if err != nil {
		return nil, err
	}

	// 将中间html写入文件
	err = ioutil.WriteFile(tmpFullName, *htmlBts, 0644)
	if err != nil {
		return nil, err
	}

	// 清除临时html文件
	defer func() {
		if tool.IsExists(tmpFullName) {
			log.Debugf("清除临时html文件 %s", tmpFullName)
			err := os.Remove(tmpFullName)
			if err != nil {
				log.Errorln(err)
			}
		}
	}()

	// 路径符合chrome要求，替换 # 为 %23
	chromeTmpURI := strings.Replace(tmpFullName, "#", "%23", -1)

	// 将html文件转换成pdf byte
	pdfBts, err := Print(r.Parmas.ExecPath,
		"file://"+chromeTmpURI,
		r.Parmas.PageFormat,
		r.Parmas.PageOrientation,
		r.Parmas.PageMargin)

	return pdfBts, nil
}

// ParseHtml ...
func (r *FileRequest) parseHtml() (*[]byte, error) {
	// 路径符合chrome要求，替换 # 为 %23
	chromeURI := strings.Replace(r.AbsInPath, "#", "%23", -1)
	// 将html文件转换成pdf byte
	pdfBts, err := Print(r.Parmas.ExecPath,
		"file://"+chromeURI,
		r.Parmas.PageFormat,
		r.Parmas.PageOrientation,
		r.Parmas.PageMargin)
	if err != nil {
		return nil, err
	}
	return pdfBts, nil
}

// URLRequest ...
type URLRequest struct {
	URL        string
	Parmas     model.Parmas
	EscapedURL string
	AbsOutPath string
}

func NewURLRequest(inPath string, parmas model.Parmas) (URLRequest, error) {
	escapedURL := url.QueryEscape(inPath)
	absOutPath, err := tool.GetOutFullName(escapedURL, parmas)
	if err != nil {
		return URLRequest{}, err
	}
	return URLRequest{
		URL:        inPath,
		Parmas:     parmas,
		EscapedURL: escapedURL,
		AbsOutPath: absOutPath,
	}, nil
}

// SaveAsFile ...
func (u *URLRequest) SaveAsFile(bts *[]byte) error {
	err := ioutil.WriteFile(u.AbsOutPath, *bts, 0644)
	if err != nil {
		return err
	}
	log.Infof("成功保存文件：%s", u.AbsOutPath)
	return nil
}

// GetBts ...
func (u *URLRequest) GetBts() (*[]byte, error) {
	return u.parseURL()
}

// ParseURL ...
func (u *URLRequest) parseURL() (*[]byte, error) {
	pdfBts, err := Print(u.Parmas.ExecPath,
		u.URL, u.Parmas.PageFormat,
		u.Parmas.PageOrientation,
		u.Parmas.PageMargin)
	if err != nil {
		return nil, err
	}
	return pdfBts, nil
}
