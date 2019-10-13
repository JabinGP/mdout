package parse


import(
	// 输出日志信息
	"log"
	// 生成错误
	"errors"

	// i/o相关
	"bytes"
	"io/ioutil" // 读写文件

	// markdown 解析为tag标签
	"gitlab.com/golang-commonmark/markdown" 

	// 无头chrome api
	"github.com/chromedp/chromedp"
	"github.com/chromedp/cdproto/page"
	// 上下文关系，用于chromedp
	"context"

	// go语言的JQuery，用于html模板拼接
	"github.com/PuerkitoBio/goquery"

	// 获取家路径
	. "../dir"
)

// 将源文件字节流转为html标签字节流
func ParseMd(sourceByteArr []byte)(error,*[]byte){
	// 将输入的源md文件解析为html标签，存在[]byte中
	md := markdown.New(markdown.XHTMLOutput(true))
	tagByteArr :=[]byte(md.RenderToString(sourceByteArr))
	return nil,&tagByteArr;
}

// 将标签拼接为完整的，可独立渲染的html（不依赖外部css，js文件）
func AssembleTag(theme string ,tagByteArr *[]byte)(error,*[]byte){

	// 获取用户home目录
	homeDir := HomeDir()
	if homeDir==""{
		return errors.New("获取资源目录失败"),nil
	}

	// 资源文件路径
	var templateHtmlPath string = homeDir+"/.mdout/h5/template/template.html" // html模板
	var hljsJsPath string = homeDir+"/.mdout/h5/template/highlight.pack.js"		// hljs库文件
	var initHljsJsPath string = homeDir+"/.mdout/h5/template/initHighlight.js" // 调用hljs文件
	var pageCssPath string = homeDir+"/.mdout/h5/theme/"+theme+"/page.css"	// 页面css
	var hljsCssPath string = homeDir+"/.mdout/h5/theme/"+theme+"/hljs.css"	// hljs css

	// 读取html模板
	templateHtmlByteArr, err := ioutil.ReadFile(templateHtmlPath)
	if err!=nil{
		log.Fatal(err)
		return err,nil
	}

	// 读取hljs库
	hljsJsByteArr, err := ioutil.ReadFile(hljsJsPath)
	if err!=nil{
		log.Fatal(err)
		return err,nil
	}

	// 读取调用hljs文件
	initHljsJsByteArr, err := ioutil.ReadFile(initHljsJsPath)
	if err!=nil{
		log.Fatal(err)
		return err,nil
	}

	// 读取页面css
	pageCssByteArr, err := ioutil.ReadFile(pageCssPath)
	if err!=nil{
		log.Fatal(err)
		return err,nil
	}

	// 读取hljs css
	hljsCssByteArr, err := ioutil.ReadFile(hljsCssPath)
	if err!=nil{
		log.Fatal(err)
		return err,nil
	}

	// 获取html模板的Reader，用于goquery
	templateHtmlReader := bytes.NewReader(templateHtmlByteArr)

	// 获取HtmlDocument对象
	templateHtmlDoc, err := goquery.NewDocumentFromReader(templateHtmlReader)
	if err!=nil{
		log.Fatal(err)
		return err,nil
	}
	
	// 页面头部插入css文件
	templateHtmlDoc.Find("head").AppendHtml(`<style type="text/css">`+string(pageCssByteArr)+string(hljsCssByteArr)+`</style>`)
	// 页面头部插入hljs文件
	templateHtmlDoc.Find("head").AppendHtml(`<script type="text/javascript">`+string(hljsJsByteArr)+`</script>`)
	// 页面头部插入hljs调用文件
	templateHtmlDoc.Find("head").AppendHtml(`<script type="text/javascript">`+string(initHljsJsByteArr)+`</script>`)

	// 将markdown标签插入html
	templateHtmlDoc.Find("body").AppendHtml(string(*tagByteArr))

	// 获取凭借后的html字符串
	assembledHtml ,err :=templateHtmlDoc.Html()
	if err!=nil{
		log.Fatal(err)
		return err,nil
	}

	// 构建byte数组
	assembledHtmlByteArr := []byte(assembledHtml)
	return nil,&assembledHtmlByteArr;
}

// 读取对应路径的html，渲染pdf到pdfByteArr
func PrintHtmlToPdf(htmlPath string)(error, *[]byte){

	var pdfByteArr []byte
	// 创建 context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// chromdp 执行打印任务
	err := chromedp.Run(ctx, chromedpPrintTasks(htmlPath, &pdfByteArr))
	if  err != nil {
		log.Fatal(err)
		return err,nil
	}

	return err,&pdfByteArr
}

// 定义chromed打印任务
func chromedpPrintTasks(htmlPath string, pdfByteArr *[]byte)chromedp.Tasks {
	return chromedp.Tasks{
		// 浏览器定位到对应资源
		chromedp.Navigate(htmlPath),

		// 执行打印操作
		chromedp.ActionFunc(func(ctx context.Context) error {

			// 设置打印参数，A4=8.27*11.69inch
			printToPDFParams := &page.PrintToPDFParams{PrintBackground:true,PaperWidth:8.27,PaperHeight:11.69}

			// 获取pdf字节数组
			pdfTmpByteArr, _,err := (printToPDFParams.WithPrintBackground(true)).Do(ctx)
			if err != nil {
				return err
			}

			// 将pdf字节数组赋值给对应指针
			*pdfByteArr = pdfTmpByteArr;
			return nil
		}),
	}
}