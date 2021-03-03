package requester

import (
	"net/url"

	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/tool"
)

func buildURLReq(inPath string, params model.Params) (*Request, error) {
	escapedURL := url.QueryEscape(inPath)

	absOutPath, err := tool.GetOutFullName(escapedURL, params)
	if err != nil {
		return nil, err
	}
	var req = Request{
		Params:    params,
		InType:    "url",
		InPath:    inPath,
		AbsInPath: inPath,
	}

	req.OutType = "pdf" // Must be pdf when input type is url
	req.AbsOutPath = absOutPath

	return &req, nil
}
