package parser

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"

	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/requester"
	"github.com/JabinGP/mdout/static"
	"github.com/JabinGP/mdout/tool"
)

// TagFileParser Parse tag file to tag bytes
type TagFileParser struct {
}

// Parse tag file  to tag bytes
func (t *TagFileParser) Parse(req *requester.Request) error {
	// 读取源文件
	srcBytes, err := ioutil.ReadFile(req.AbsInPath)
	if err != nil {
		return err
	}

	req.Data = srcBytes
	req.InType = "tag-bytes"
	return nil
}

// TagBytesParser markdown parser
type TagBytesParser struct {
}

// Parse markdown to html
func (t *TagBytesParser) Parse(req *requester.Request) error {
	// tagBytes := req.Data.([]byte)

	if !tool.IsExists(filepath.FromSlash(static.ThemeFolderFullName + "/" + req.ThemeName)) {
		return fmt.Errorf("无法找到名为 %s 的主题", req.ThemeName)
	}
	log.Debugln("开始生成html...")
	htmlBytes, err := t.buildHTMLFromTemplate(req)
	if err != nil {
		log.Errorf("生成html失败: %v", err)
		return err
	}
	log.Debugln("成功生成html")

	req.Data = htmlBytes
	req.InType = "html-bytes"
	return nil
}

func (t *TagBytesParser) buildHTMLFromTemplate(req *requester.Request) ([]byte, error) {
	// 获取资源文件夹路径
	var themeDir = filepath.FromSlash(static.ThemeFolderFullName +
		"/" + req.ThemeName)
	// html模板
	var indexHTMLFullName = filepath.FromSlash(themeDir +
		"/index.html")

	// 读取模板
	tmpl, err := template.ParseFiles(indexHTMLFullName)
	if err != nil {
		log.Errorf("解析主题模板 %s 失败: %v", indexHTMLFullName, err)
		return nil, err
	}

	// 将数据放入模板
	buf := new(bytes.Buffer)

	tmpl.Execute(buf, model.TemplateData{
		ThemePath:    template.URL(filepath.ToSlash(themeDir)),
		MarkdownBody: template.HTML(req.Data.([]byte)),
	})

	return buf.Bytes(), nil
}
