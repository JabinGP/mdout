package requester

import (
	"net/url"

	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/tool"
)

func buildURLReq(inPath string, parmas model.Parmas) (*Request, error) {
	escapedURL := url.QueryEscape(inPath)

	absOutPath, err := tool.GetOutFullName(inPath, parmas)
	if err != nil {
		return nil, err
	}
	var req = Request{
		Parmas:    parmas,
		InType:    "url",
		InPath:    inPath,
		AbsInPath: escapedURL,
	}
	req.OutType = "pdf" // Must be pdf when input type is url
	req.AbsOutPath = absOutPath

	return &req, nil
}
