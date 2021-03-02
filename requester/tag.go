package requester

import (
	"errors"

	"github.com/JabinGP/mdout/tool"
)

func buildTagReq(req *Request) error {
	// 检查输出类型
	if !tool.CheckType(req.OutType, []string{"html", "pdf"}) {
		return errors.New("非法的输出类型：" + req.OutType)
	}

	req.InType = "tag-file"
	return nil
}
