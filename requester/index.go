package requester

import (
	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/tool"
)

// Request parsed by parser
type Request struct {
	model.Params
	InType     string
	InPath     string
	AbsInPath  string
	AbsOutPath string
	Data       interface{}
	DeferFuncs []func() // 用于释放临时资源
}

// NewRequest return a request
func NewRequest(inPath string, params model.Params) (*Request, error) {
	inType, err := tool.GetType(inPath)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	if inType == "url" {
		return buildURLReq(inPath, params)
	}
	return buildFileReq(inPath, params)
}
