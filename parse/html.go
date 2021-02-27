package parse

import (
	"context"

	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/model"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// HTMLParser parse html to pdf
type HTMLParser struct {
	parmas model.Parmas
}

// Parse html to pdf
func (h *HTMLParser) Parse(htmlPath string) ([]byte, error) {
	log.Debugln("准备开始打印...")

	ctx, cancelFunc := h.getChromedpCtx()
	defer cancelFunc()

	log.Debugln("正在加载：" + htmlPath + " ...")

	// 打开html
	err := h.openHTML(ctx, htmlPath)
	if err != nil {
		return nil, err
	}

	// 等待html页面就绪
	err = h.waitForHTMLReady(ctx)
	if err != nil {
		return nil, err
	}

	// 打印
	log.Infoln("正在生成pdf...")
	return h.print(ctx)
}

func (h *HTMLParser) getChromedpCtx() (context.Context, context.CancelFunc) {
	if h.parmas.ExecPath != "" {
		log.Debugln("指定chrome执行路径：" + h.parmas.ExecPath)
		var cancel1, cancel2 context.CancelFunc
		// 定义执行路径以及参数
		opts := []chromedp.ExecAllocatorOption{
			chromedp.ExecPath(h.parmas.ExecPath),
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

func (h *HTMLParser) openHTML(ctx context.Context, htmlPath string) error {
	// 定位到文件
	err := chromedp.Run(ctx, chromedp.Navigate(htmlPath))
	if err != nil {
		return err
	}
	log.Debugln("加载：" + htmlPath + " 完成")
	log.Debugln("正在获取页面渲染状态...")
	return nil
}

func (h *HTMLParser) waitForHTMLReady(ctx context.Context) error {
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

func (h *HTMLParser) print(ctx context.Context) ([]byte, error) {
	var pdfBytes []byte
	// 通过尺寸获取宽高
	paperWidth, paperHeight := getPaperWidthAndHeight(h.parmas.PageFormat)
	// 获取纸张方向
	isLandscape := getPaperOrientation(h.parmas.PageOrientation)
	// 获取页面边距
	marginArr, err := getMargin(h.parmas.PageMargin)
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
