package parser

import (
	"fmt"
	"io/ioutil"

	"github.com/JabinGP/mdout/config"
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

	md := markdown.New(
		markdown.XHTMLOutput(config.Obj.Runtime.EnableXHTMLOutput),
		markdown.HTML(config.Obj.Runtime.EnableHTMLTag),
	)

	tokens := md.Parse(mdBytes)

	for i := range tokens {
		token, ok := tokens[i].(*markdown.Fence)
		if ok && token.Params == "mermaid" {
			newToken := &markdown.HTMLBlock{}
			newToken.Content = fmt.Sprintf(`<div class="mermaid">
%s</div>
`, token.Content)
			newToken.Map = token.Map
			newToken.Lvl = token.Lvl
			tokens[i] = newToken
		}
	}

	tagBytes := []byte(md.RenderTokensToString(tokens))
	log.Debugln("解析markdown成功")

	req.Data = tagBytes
	req.InType = "tag-bytes"
	return nil
}
