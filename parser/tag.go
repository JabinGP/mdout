package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/requester"
	"github.com/JabinGP/mdout/static"
	"github.com/JabinGP/mdout/tool"
	"github.com/PuerkitoBio/goquery"
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
	tagBytes := req.Data.([]byte)

	if !tool.IsExists(filepath.FromSlash(static.ThemeFolderFullName + "/" + req.ThemeName)) {
		return fmt.Errorf("无法找到名为 %s 的主题", req.ThemeName)
	}
	log.Debugln("开始生成html...")
	// 获取资源文件夹路径
	var themeDir = filepath.FromSlash(static.ThemeFolderFullName +
		"/" + req.ThemeName)
	// html模板
	var indexHTMLFullName = filepath.FromSlash(themeDir +
		"/index.html")

	// 页面模板
	indexHTMLBytes, err := ioutil.ReadFile(indexHTMLFullName)
	if err != nil {
		return err
	}

	// 获取主体html模板的Reader，用于goquery
	indexHTMLReader := bytes.NewReader(indexHTMLBytes)

	// 获取HtmlDocument对象
	indexHTMLDoc, err := goquery.NewDocumentFromReader(indexHTMLReader)
	if err != nil {
		return err
	}

	// 拼装页面
	indexHTMLDoc.Find(".markdown-body").AppendHtml(string(tagBytes)) // 将markdown标签插入html

	// 将link标签和script标签中的{{homePath}}和{{theme}}替换成为实际路径
	indexHTMLDoc.Find("link").Each(func(i int, selection *goquery.Selection) {
		linkHref, ok := selection.Attr("href") // 获取href属性
		if ok {                                // 如果获取成功
			// 替换
			linkHref = strings.Replace(linkHref, `{{themePath}}`, filepath.ToSlash(themeDir), -1)
			// 生效
			selection.SetAttr("href", linkHref)
		}
	})
	indexHTMLDoc.Find("script").Each(func(i int, selection *goquery.Selection) {
		srcPath, ok := selection.Attr("src") // 获取src属性
		if ok {                              // 如果获取成功
			// 替换
			srcPath = strings.Replace(srcPath, `{{themePath}}`, filepath.ToSlash(themeDir), -1)
			// 生效
			selection.SetAttr("src", srcPath)
		}
	})

	// 获取拼接后的html字符串
	assembledHTML, err := indexHTMLDoc.Html()
	if err != nil {
		return err
	}

	// 构建byte数组
	htmlBytes := []byte(assembledHTML)
	log.Debugln("成功生成html")

	req.Data = htmlBytes
	req.InType = "html-bytes"
	return nil
}
