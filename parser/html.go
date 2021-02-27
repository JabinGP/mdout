package parser

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/requester"
	"github.com/JabinGP/mdout/tool"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// HTMLBytesParser parse html to pdf
type HTMLBytesParser struct {
}

func (h *HTMLBytesParser) Parse(req *requester.Request) error {
	htmlBytes := req.Data.([]byte)
	// 构建临时html文件路径
	tmpDir, tmpName, _, err := tool.GetDirNameExt(req.AbsInPath)
	tmpExt := "html"
	tmpFullName, err := tool.Abs(tmpDir + "/" + "tmp." + tmpName + "." + tmpExt)
	if err != nil {
		return err
	}
	// 将中间html写入文件
	err = ioutil.WriteFile(tmpFullName, htmlBytes, 0644)
	if err != nil {
		return err
	}

	// 更改浏览器打开路径为临时文件
	req.AbsInPath = tmpFullName

	// 清除临时html文件任务加入 defer 队列
	req.DeferFuncs = append(req.DeferFuncs, func() {
		if tool.IsExists(tmpFullName) {
			log.Debugf("清除临时html文件 %s", tmpFullName)
			err := os.Remove(tmpFullName)
			if err != nil {
				log.Errorln(err)
			}
		}
	})

	req.InType = "html-file"
	return nil
}

// HTMLFileParser parse html bytes to tmp html bytes
type HTMLFileParser struct {
}

// Parse html bytes to pdf bytes
func (h *HTMLFileParser) Parse(req *requester.Request) error {
	log.Debugln("准备开始打印...")
	htmlPath := req.AbsInPath
	ctx, cancelFunc := h.getChromedpCtx(*req)
	defer cancelFunc()

	log.Debugln("正在加载：" + htmlPath + " ...")

	// 打开html
	err := h.openHTML(ctx, htmlPath)
	if err != nil {
		return err
	}

	// 等待html页面就绪
	err = h.waitForHTMLReady(ctx)
	if err != nil {
		return err
	}

	// 打印
	log.Infoln("正在生成pdf...")
	pdfBytes, err := h.print(ctx, *req)
	if err != nil {
		return err
	}

	req.Data = pdfBytes
	req.InType = "pdf-bytes"
	return nil
}

func (h *HTMLFileParser) getChromedpCtx(req requester.Request) (context.Context, context.CancelFunc) {
	if req.ExecPath != "" {
		log.Debugln("指定chrome执行路径：" + req.ExecPath)
		var cancel1, cancel2 context.CancelFunc
		// 定义执行路径以及参数
		opts := []chromedp.ExecAllocatorOption{
			chromedp.ExecPath(req.ExecPath),
			chromedp.Headless,
			chromedp.DisableGPU,
		}
		// 创建 context
		allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
		// defer cancel1()
		ctx, cancel2 := chromedp.NewContext(allocCtx)
		// defer cancel2()
		cancelFunc := func() {
			defer cancel1()
			defer cancel2()
		}
		return ctx, cancelFunc
	}

	log.Debugln("未指定执行路径，自动寻找chrome执行路径")
	return chromedp.NewContext(context.Background())
}

func (h *HTMLFileParser) openHTML(ctx context.Context, htmlPath string) error {
	// 定位到文件
	err := chromedp.Run(ctx, chromedp.Navigate(htmlPath))
	if err != nil {
		return err
	}
	log.Debugln("加载：" + htmlPath + " 完成")
	log.Debugln("正在获取页面渲染状态...")
	return nil
}

func (h *HTMLFileParser) waitForHTMLReady(ctx context.Context) error {
	var isJabinGP bool
	// eval是下策，因为官方的查找元素api都无法在元素不存在的时候正常运行，会一直卡住
	err := chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector("#jabingp")!=null`, &isJabinGP))
	if err != nil {
		return err
	}

	// 无同步渲染标记
	if !isJabinGP {
		log.Infoln("页面不支持同步渲染进度，跳过同步渲染并打印pdf")
		return nil
	}

	// 有同步渲染标记
	log.Infoln("页面支持同步渲染进度，开始同步...")

	// 渲染消息
	var renderInfo string
	var renderInfoSaved string

	// 等待并获取渲染状态
	for {
		renderInfoSaved = renderInfo
		err = chromedp.Run(ctx, chromedp.InnerHTML(`#jabingp`, &renderInfo, chromedp.ByID))
		if err != nil {
			return err
		}

		// 消息去重
		if renderInfoSaved != renderInfo {
			log.Infoln(renderInfo)
		}
		if renderInfo == "渲染完成！" {
			break
		}
	}

	return nil
}

func (h *HTMLFileParser) print(ctx context.Context, req requester.Request) ([]byte, error) {
	var pdfBytes []byte
	// 通过尺寸获取宽高
	paperWidth, paperHeight := getPaperWidthAndHeight(req.PageFormat)
	// 获取纸张方向
	isLandscape := getPaperOrientation(req.PageOrientation)
	// 获取页面边距
	marginArr, err := getMargin(req.PageMargin)
	if err != nil {
		log.Errorln(err)
		log.Infoln("转换页面边距出错，将以默认边距打印")
	}

	// 开始打印
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(
			func(ctx context.Context) error {
				// 设置打印参数，A4=8.27*11.69inch
				printToPDFParams := &page.PrintToPDFParams{PrintBackground: true,
					PaperWidth:   paperWidth,
					PaperHeight:  paperHeight,
					Landscape:    isLandscape,
					MarginTop:    marginArr[0],
					MarginRight:  marginArr[1],
					MarginBottom: marginArr[2],
					MarginLeft:   marginArr[3],
				}

				// 获取pdf字节数组
				pdfTmpBytes, _, err := (printToPDFParams.WithPrintBackground(true)).Do(ctx)
				if err != nil {
					return err
				}

				// 将pdf字节数组赋值给对应指针
				pdfBytes = pdfTmpBytes
				return nil
			}),
	)

	if err != nil {
		return nil, err
	}
	return pdfBytes, err
}
