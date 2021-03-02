package requester

import (
	"errors"

	"github.com/JabinGP/mdout/tool"
)

func buildHTMLReq(req *Request) error {
	// 检查输出类型
	if !tool.CheckType(req.OutType, []string{"pdf"}) {
		return errors.New("非法的输出类型：" + req.OutType)
	}

	req.InType = "html-file"
	return nil
}
