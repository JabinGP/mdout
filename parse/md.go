package parse

import (
	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/model"
	"gitlab.com/golang-commonmark/markdown"
)

// MDParser markdown parser
type MDParser struct {
	parmas model.Parmas
}

// Parse markdown to html
func (m *MDParser) Parse(mdBytes []byte) ([]byte, error) {
	log.Debugln("开始解析markdown...")
	// 将输入的源md文件解析为html标签，存在[]byte中
	md := markdown.New(markdown.XHTMLOutput(true), markdown.HTML(true))
	tagBytes := []byte(md.RenderToString(mdBytes))
	log.Debugln("解析markdown成功")
	return tagBytes, nil
}
