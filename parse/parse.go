package parse

import (
	// 输出日志信息
	"log"
	"strconv"
	"strings"

	// 生成错误
	"errors"

	// i/o相关
	"bytes"
	"io/ioutil" // 读写文件

	// markdown 解析为tag标签
	"gitlab.com/golang-commonmark/markdown"

	// 无头chrome api
	"github.com/JabinGP/mdout/tool"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	// 上下文关系，用于chromedp
	"context"

	// go语言的JQuery，用于html模板拼接
	"github.com/PuerkitoBio/goquery"
)

// Md 将源文件字节流转为html标签字节流
func Md(sourceByteArr []byte) (*[]byte, error) {
	log.Println("开始解析markdown")
	// 将输入的源md文件解析为html标签，存在[]byte中
	md := markdown.New(markdown.XHTMLOutput(true))
	tagByteArr := []byte(md.RenderToString(sourceByteArr))
	log.Println("解析markdown成功")
	return &tagByteArr, nil
}

// AssembleTag 将标签拼接为完整的，可独立渲染的html（不依赖外部css，js文件）
func AssembleTag(theme string, tagBytes *[]byte) (*[]byte, error) {

	log.Println("开始生成html")

	// 获取用户home目录
	homeDir := tool.GetHomeDir()
	if homeDir == "" {
		return nil, errors.New("获取资源目录失败")
	}

	// 获取资源文件夹路径
	var themeDir = homeDir + `/mdout/theme/` + theme

	// 资源文件路径
	var indexHTMLFullName = themeDir + `/index.html` // html模板

	// 页面模板
	indexHTMLBytes, err := ioutil.ReadFile(indexHTMLFullName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// 获取主体html模板的Reader，用于goquery
	indexHTMLReader := bytes.NewReader(indexHTMLBytes)

	// 获取HtmlDocument对象
	indexHTMLDoc, err := goquery.NewDocumentFromReader(indexHTMLReader)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// 拼装页面
	indexHTMLDoc.Find(".markdown-body").AppendHtml(string(*tagBytes)) // 将markdown标签插入html

	// 将link标签和script标签中的{{homePath}}和{{theme}}替换成为实际路径
	indexHTMLDoc.Find("link").Each(func(i int, selection *goquery.Selection) {
		linkHref, ok := selection.Attr("href") // 获取href属性
		if ok {                                // 如果获取成功
			// 替换
			linkHref = strings.Replace(linkHref, `{{homePath}}`, homeDir, -1)
			linkHref = strings.Replace(linkHref, `{{theme}}`, theme, -1)
			// 生效
			selection.SetAttr("href", linkHref)
		}
	})
	indexHTMLDoc.Find("script").Each(func(i int, selection *goquery.Selection) {
		srcPath, ok := selection.Attr("src") // 获取src属性
		if ok {                              // 如果获取成功
			// 替换
			srcPath = strings.Replace(srcPath, `{{homePath}}`, homeDir, -1)
			srcPath = strings.Replace(srcPath, `{{theme}}`, theme, -1)
			// 生效
			selection.SetAttr("src", srcPath)
		}
	})

	// 获取拼接后的html字符串
	assembledHTML, err := indexHTMLDoc.Html()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// 构建byte数组
	assembledHTMLBytes := []byte(assembledHTML)
	log.Println("成功生成html")
	return &assembledHTMLBytes, nil
}

// Print 读取对应路径的html，渲染pdf保存到pdfBytes
func Print(htmlPath string, pageFormat string, pageOrientation string, pageMargin string) (*[]byte, error) {

	log.Println("准备开始打印")
	var pdfBytes []byte
	// 创建 context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	log.Println("正在加载" + htmlPath)
	// chromdp 执行打印任务
	// 定位到文件
	err := chromedp.Run(ctx, chromedp.Navigate(htmlPath))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("加载" + htmlPath + "完成")
	log.Println("正在获取页面渲染状态")

	// err = chromedp.Run(ctx, chromedp.WaitReady(`.markdown-body`))
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }
	// 检查是否有同步渲染标记
	var isJabinGP bool
	// eval是下策，因为官方的查找元素api都无法在元素不存在的时候正常运行，会一直卡住
	err = chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector("#jabingp")!=null`, &isJabinGP))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if isJabinGP { // 有同步渲染标记

		log.Println("页面支持同步渲染进度，开始同步")
		// 渲染消息
		var renderInfo string
		var renderInfoSaved string

		// 获取渲染状态
		for {
			renderInfoSaved = renderInfo
			err = chromedp.Run(ctx, chromedp.InnerHTML(`#jabingp`, &renderInfo, chromedp.ByID))
			if err != nil {
				log.Println(err)
				return nil, err
			}

			// 消息去重
			if renderInfoSaved != renderInfo {
				log.Println(renderInfo)
			}
			if renderInfo == "渲染完成！" {
				break
			}
		}
	} else {
		log.Println("页面不支持同步渲染进度，跳过同步渲染并打印")
	}

	// 打印
	log.Println("开始打印，正在等待打印机渲染pdf")

	// 通过尺寸获取宽高
	paperWidth, paperHeight := getPaperWidthAndHeight(pageFormat)
	// 获取纸张方向
	isLandscape := getPaperOrientation(pageOrientation)
	// 获取页面边距
	marginTop, marginRight, marginBottom, marginLeft, err := getMargin(pageMargin)
	if err != nil {
		log.Print(err)
		log.Println("，将以默认边距打印")
	}

	// 开始打印
	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		// 设置打印参数，A4=8.27*11.69inch
		printToPDFParams := &page.PrintToPDFParams{PrintBackground: true,
			PaperWidth:   paperWidth,
			PaperHeight:  paperHeight,
			Landscape:    isLandscape,
			MarginTop:    marginTop,
			MarginRight:  marginRight,
			MarginBottom: marginBottom,
			MarginLeft:   marginLeft,
		}

		// 获取pdf字节数组
		pdfTmpByteArr, _, err := (printToPDFParams.WithPrintBackground(true)).Do(ctx)
		if err != nil {
			return err
		}

		// 将pdf字节数组赋值给对应指针
		pdfBytes = pdfTmpByteArr
		return nil
	}),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("渲染pdf成功，准备保存文件")

	return &pdfBytes, err
}

// 通过格式，获取打印纸张的宽高
func getPaperWidthAndHeight(pageFormat string) (float64, float64) {
	var paperWidth, paperHeight float64
	switch strings.ToLower(pageFormat) {
	case "legal":
		paperWidth, paperHeight = splitWidthXHeight("216x356")
	case "letter":
		paperWidth, paperHeight = splitWidthXHeight("216x279")
	case "tabloid":
		paperWidth, paperHeight = splitWidthXHeight("279x356")
	case "ledger":
		paperWidth, paperHeight = splitWidthXHeight("279x356")
	case "a5":
		paperWidth, paperHeight = splitWidthXHeight("148x210")
	case "a4":
		paperWidth, paperHeight = splitWidthXHeight("210x297")
	case "a3":
		paperWidth, paperHeight = splitWidthXHeight("297x420")
	case "a2":
		paperWidth, paperHeight = splitWidthXHeight("420x594")
	case "a1":
		paperWidth, paperHeight = splitWidthXHeight("594x841")
	case "a0":
		paperWidth, paperHeight = splitWidthXHeight("841x1189")
	default:
		paperWidth, paperHeight = splitWidthXHeight("841x1189")
	}
	return paperWidth / 25.4, paperHeight / 25.4
}

func splitWidthXHeight(widthXHeight string) (float64, float64) {
	numberArr := strings.Split(strings.ToLower(widthXHeight), "x")
	width, err := strconv.ParseFloat(numberArr[0], 32)
	height, err := strconv.ParseFloat(numberArr[1], 32)
	if err != nil {
		log.Println(err)
		log.Println("建立" + widthXHeight + "纸张尺寸表时类型转换失败，将尺寸定位A4纸210cm*297cm大小")
		return 210, 297
	}
	return width, height
}

func getPaperOrientation(pageOrientation string) bool {
	if strings.ToLower(pageOrientation) == "landscape" {
		return true
	}
	return false
}

func getMargin(pageMargin string) (float64, float64, float64, float64, error) {

	// 替换所有中文逗号为英文逗号
	pageMargin = strings.Replace(pageMargin, "，", ",", -1)

	// 去除空格
	pageMargin = strings.Replace(pageMargin, " ", "", -1)

	// 获取配置数组
	marginArr := strings.Split(pageMargin, ",")

	// 判断输入类型
	switch len(marginArr) {
	case 0:
		return 1, 1, 1, 1, errors.New("无效的输入边距值")
	case 1:
		marginAll, err := strconv.ParseFloat(marginArr[0], 64)
		if err != nil {
			return 1, 1, 1, 1, err
		}

		// 输出0会被chromedp重设为0.4默认值
		if marginAll == 0 {
			marginAll = 0.0000000000000001
		}
		return marginAll, marginAll, marginAll, marginAll, nil
	default:
		return 1, 1, 1, 1, errors.New("无法识别的输入边距类型")
	}
}
