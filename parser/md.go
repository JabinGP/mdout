package parser

import (
	"io/ioutil"

	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/requester"
	"gitlab.com/golang-commonmark/markdown"
)

// MDFileParser markdown file parser
type MDFileParser struct {
}

// Parse markdown file to markdown bytes
func (m *MDFileParser) Parse(req *requester.Request) error {
	// 读取源文件
	srcBytes, err := ioutil.ReadFile(req.AbsInPath)
	if err != nil {
		return err
	}

	req.Data = srcBytes
	req.InType = "md-bytes"
	return nil
}

// MDBytesParser markdown parser
type MDBytesParser struct {
}

// Parse markdown bytes to tag bytes
func (m *MDBytesParser) Parse(req *requester.Request) error {
	mdBytes := req.Data.([]byte)
	log.Debugln("开始解析markdown...")
	// 将输入的源md文件解析为html标签，存在[]byte中
	md := markdown.New(markdown.XHTMLOutput(true), markdown.HTML(true))
	tagBytes := []byte(md.RenderToString(mdBytes))
	log.Debugln("解析markdown成功")

	req.Data = tagBytes
	req.InType = "tag-bytes"
	return nil
}
