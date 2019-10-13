package cmd

import (
  // 输出日志
  "log"
  // 错误生成
  "errors"

  // 路径处理
  "path/filepath"
  // 获取系统路径信息
  "os"

  // 字符串处理
  "strconv"
  "strings"
  "net/url"
  "regexp"

  // 读写文件
  "io/ioutil"

  // 命令行框架
  "github.com/spf13/cobra"
 
  // 获取home目录，判断路径是否存在
  . "../dir"
  // md解析，html转pdf
  . "../parse"
)

var (

  // 用户home目录，用于获取模板资源文件，存储中间产物
  homeDir string  = HomeDir()
  // 命令行可变参数

  // 文件输出路径
  outPath string
  // 目标文件类型
  targetType string
  // 主题
  theme string
  // 打印页面格式
  pageFormat string
  // 打印页面方向
  pageOrientation string
  // 运行环境设置
  chromePath string
  // 渲染超时设置
  timeout int


  // 根命令
  rootCmd = &cobra.Command{
    Use: "mdout [资源路径]",
    Version: "0.5",
    Short: "将markdown、html、url网络资源转换成pdf",
    Long: "读取对应路径的资源，在内部转换成html文件，在将html文件渲染为pdf文件保存，默认保存在markdown文件所在路径，默认输出为pdf",
    Args: cobra.MinimumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string)error {      
      // 判断是否成功获取模板、主题资源
      if(homeDir==""){
        err := errors.New("获取模板、主题资源失败")
        log.Fatal(err)
        return err
      }

      // 获得输出路径，输出文件夹路径
      outPath,err:= filepath.Abs(outPath)
      if err!= nil{
        log.Fatal(err)
        return err
      }

      match, err := regexp.MatchString(`(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`, args[0])
      if(!match){ //如果不属于url路径，按文件路径处理
        // 获取输入文件路径，文件夹路径，文件名
        sourceFilePath,err := filepath.Abs(args[0])
        if err!= nil{
          log.Fatal(err)
          return err
        }
        // 解析文件路径
        sourceFileName := filepath.Base(sourceFilePath)// 源文件全名带后缀
        sourceExt := filepath.Ext(sourceFileName)  // 源文件后缀

        // 根据文件不同类型处理
        switch sourceExt {

          case ".html":

            return extCaseHtml(sourceFilePath,outPath)

          default:

            return extCaseMd(sourceFilePath,targetType,outPath)

        }

      }else{ // 为url

        return extCaseUrl(args[0],outPath)

      }
    },
  }
)

// 当输入为url时
func extCaseUrl(sourceUrl string,outPath string) error{
  // 目标pdf文件路径
  targetPdfFilePath := outPath+"/"+url.QueryEscape(sourceUrl)+".pdf"
  
  // 将html文件转换成pdf byte
  err,pdfByteArr := PrintHtmlToPdf(sourceUrl);

  // 将得到的pdf byte写入文件
  err =ioutil.WriteFile(targetPdfFilePath,*pdfByteArr,0644)
  if err!=nil{
    log.Fatal(err)
    return err
  }

  // 输出成功信息
  log.Println("Successfully saved pdf in ：" +targetPdfFilePath)
  return nil
}

// 当输入为一个html文件时
func extCaseHtml(sourceFilePath string,outPath string) error {

  // 获取路径
  sourceFileFullName := filepath.Base(sourceFilePath)// 源文件全名带后缀
  sourceExt := filepath.Ext(sourceFileFullName)  // 源文件后缀
  sourceFileName := strings.TrimSuffix(sourceFileFullName, sourceExt) //获取文件名不带后缀

  // 将html文件转换成pdf byte
  err,pdfByteArr := PrintHtmlToPdf("file://"+sourceFilePath);

  // 目标pdf文件路径
  targetPdfFilePath := outPath+"/"+sourceFileName+".pdf"

  // 将得到的pdf byte写入文件
  err =ioutil.WriteFile(targetPdfFilePath,*pdfByteArr,0644)
  if err!=nil{
    log.Fatal(err)
    return err
  }

  // 输出成功信息
  log.Println("Successfully saved pdf in ：" +targetPdfFilePath)
  return nil
}


// 当输入为一个md文件时
func extCaseMd(sourceFilePath string,targetType string,outPath string) error {

  // 获取路径
  sourceDir := filepath.Dir(sourceFilePath)   // 源文件绝对路径
  sourceFileFullName := filepath.Base(sourceFilePath)// 源文件全名带后缀
  sourceExt := filepath.Ext(sourceFileFullName)  // 源文件后缀
  sourceFileName := strings.TrimSuffix(sourceFileFullName, sourceExt) //获取文件名不带后缀

  // 读取源文件
  sourceFileByteArr, err := ioutil.ReadFile(sourceFilePath)
  if err!=nil{
    log.Fatal(err)
    return err
  }

  // md解析
  err,mdByteArr := ParseMd(sourceFileByteArr)
  if err!=nil{
    log.Fatal(err)
    return err
  }

  switch targetType {

    case "tag":

      // 目标html标签文件路径
      targetTagFilePath := sourceDir+"/"+sourceFileName+".html"
      // 将得到的tag写入文件
      err =ioutil.WriteFile(outPath+"/"+sourceFileName+".html",*mdByteArr,0644)
      if err!=nil{
        log.Fatal(err)
        return err
      }

      // 输出成功信息
      log.Println("Successfully saved html in ：" +targetTagFilePath)
      return nil

    case "html":

      // tag拼接
      err,htmlByteArr := AssembleTag(theme,mdByteArr)
      if err!=nil{
        log.Fatal(err)
        return err
      }

      // 目标html文件路径
      targetHtmlFilePath := sourceDir+"/"+sourceFileName+".html"
      // 将得到的html写入文件
      err =ioutil.WriteFile(targetHtmlFilePath,*htmlByteArr,0644)
      if err!=nil{
        log.Fatal(err)
        return err
      }

      // 输出成功信息
      log.Println("Successfully saved html in ：" +targetHtmlFilePath)
      return nil

    case "pdf":

      // tag拼接
      err,htmlByteArr := AssembleTag(theme,mdByteArr)
      if err!=nil{
        log.Fatal(err)
        return err
      }

      // 中间Html文件路径
      tmpHtmlFilePath := sourceDir+"/."+sourceFileName+".tmp.html"
      // 将得到的html写入文件
      err =ioutil.WriteFile(tmpHtmlFilePath,*htmlByteArr,0644)
      if err!=nil{
        log.Fatal(err)
        return err
      }

      // 将html文件转换成pdf byte
      err,pdfByteArr := PrintHtmlToPdf("file://"+tmpHtmlFilePath);

      // 目标pdf文件路径
      targetPdfFilePath := outPath+"/"+sourceFileName+".pdf"
      // 将得到的pdf byte写入文件
      err =ioutil.WriteFile(targetPdfFilePath,*pdfByteArr,0644)
      if err!=nil{
        log.Fatal(err)
        return err
      }

      // 清除中间文件
      if(IsExists(tmpHtmlFilePath)){
        err := os.Remove(tmpHtmlFilePath)
        if err !=nil{
          log.Fatal(err)
          return err
        }
      }

      // 输出成功信息
      log.Println("Successfully saved pdf in ：" +targetPdfFilePath)
      return nil
  }

  return nil
}

// 输出flag信息调试
func showParams(){
  log.Println("outPath","("+outPath+")",
    "targetType","("+targetType+")",
    "theme","("+theme+")",
    "pageFormat","("+pageFormat+")",
    "pageOrientation","("+pageOrientation+")",
    "chromePath","("+chromePath+")",
    "timeout","("+strconv.Itoa(timeout)+")",
  )
}

func init() {
  // 添加 flag
  // 变量 长名 短名 默认值 帮助说明
  rootCmd.Flags().StringVarP(&outPath,"outpath", "o", "./", "文件输出的路径")
  //rootCmd.Flags().StringVarP(&pageFormat,"format", "f", "A4", "打印的页面大小：A5-A3、Legal、Letter、Tabloid")
  rootCmd.Flags().StringVarP(&targetType,"type", "t", "pdf", "输出的文件类型:tag、html、pdf")
  rootCmd.Flags().StringVarP(&theme,"theme", "e", "default", "界面的主题，可放入自定义主题包后修改")
 // rootCmd.Flags().StringVarP(&pageOrientation,"orientation", "r", "portrait", "打印的页面方向,可选portrait（纵向）、landscape（横向）")
 // rootCmd.Flags().IntVarP(&timeout,"timeout", "i", 10000, "超时时间（ms）")
  //rootCmd.Flags().StringVarP(&chromePath,"chrome", "c", "", "chrome所在路径")
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    log.Println(err)
    os.Exit(1)
  }
}
