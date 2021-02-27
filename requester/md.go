package requester

import (
	"errors"

	"github.com/JabinGP/mdout/tool"
)

func buildMDReq(req *Request) error {
	// 检查输出类型
	if !tool.CheckType(req.OutType, []string{"tag", "html", "pdf"}) {
		return errors.New("非法的输出类型：" + req.OutType)
	}

	req.InType = "md-file"
	return nil
}
