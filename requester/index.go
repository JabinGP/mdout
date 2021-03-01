package requester

import (
	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/tool"
)

// Request parsed by parser
type Request struct {
	model.Parmas
	InType     string
	InPath     string
	AbsInPath  string
	AbsOutPath string
	Data       interface{}
	DeferFuncs []func() // 用于释放临时资源
}

// NewRequest return a request
func NewRequest(inPath string, parmas model.Parmas) (*Request, error) {
	inType, err := tool.GetType(inPath)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	if inType == "url" {
		return buildURLReq(inPath, parmas)
	}
	return buildFileReq(inPath, parmas)
}
