package parser

import (
	"fmt"

	"github.com/JabinGP/mdout/requester"
)

// Parser interface that contain Parse function
type Parser interface {
	Parse(*requester.Request) error
}

// Parse parse entry
func Parse(req *requester.Request) error {
	defer func() {
		// 释放资源，defFunc 执行顺序与 defFunc 加入数组的顺序相反
		for _, defFunc := range req.DeferFuncs {
			defer defFunc()
		}
	}()

	for req.InType != req.OutType+"-bytes" {
		parser, err := NewParser(req.InType)
		if err != nil {
			return err
		}
		err = parser.Parse(req)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewParser return parser according to input type
func NewParser(inType string) (Parser, error) {
	switch inType {
	case "md-file":
		return &MDFileParser{}, nil
	case "md-bytes":
		return &MDBytesParser{}, nil
	case "tag-file":
		return &TagFileParser{}, nil
	case "tag-bytes":
		return &TagBytesParser{}, nil
	case "html-bytes":
		return &HTMLBytesParser{}, nil
	case "html-file":
		return &HTMLFileParser{}, nil
	case "html-source":
		return &HTMLSourceParser{}, nil
	case "url":
		return &HTMLSourceParser{}, nil
	default:
		return nil, fmt.Errorf("无法为输入类型 %s 找到对应的 Parser 。", inType)
	}
}
