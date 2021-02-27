package cmd

import (
	"os"

	"github.com/JabinGP/mdout/parser"
	"github.com/JabinGP/mdout/requester"

	"github.com/JabinGP/mdout/tool"

	"github.com/JabinGP/mdout/config"
	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/model"
	"github.com/spf13/cobra"
)

var (
	// 命令行输入参数，与cobra命令行绑定
	cmdParmas model.Parmas
	// 根命令
	rootCmd = &cobra.Command{
		Use:     "mdout",
		Version: "0.6",
		Short:   "将markdown、html、url转换成pdf",
		Long:    "读取输入的文件，在内部转换成html，并将html渲染为pdf保存",
		Args:    cobra.MinimumNArgs(1),
		RunE:    rootRunE,
	}
)

// init 包初始化
func init() {
	initRootCmdFlags()
	addCmdToRoot()
	setConfigLoggerLevel()
}

// Execute 程序执行入口
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}
}

// 根目录运行函数
func rootRunE(cmd *cobra.Command, args []string) error {

	// 运行时日志等级
	setRuntimeLoggerLevel()

	// 输出调试参数
	showParams()

	// 获取用户输入路径
	inPath := args[0]

	// Build request
	req, err := requester.NewRequest(inPath, cmdParmas)
	if err != nil {
		return err
	}

	// Parse request
	err = parser.Parse(req)
	if err != nil {
		return err
	}

	// Save
	err = tool.SaveFile(req.Data.([]byte), req.AbsOutPath)
	if err != nil {
		return err
	}

	log.Infof("成功保存文件：%s", req.AbsOutPath)
	return nil
}

func setRuntimeLoggerLevel() {
	if cmdParmas.Verbose {
		level, _ := tool.TransformToLogrusLevel("debug")
		log.SetStdoutLevel(level)
	}
}

// func distribute(inPath string, cmdParmas model.Parmas) error {
// 	// 获取输入参数类型
// 	inType, err := tool.GetType(inPath)
// 	if err != nil {
// 		log.Errorln(err)
// 		return err
// 	}

// 	// 根据不同的输入类型处理，定位到不同的执行函数
// 	switch inType {
// 	case "url":
// 		return inURL(inPath, cmdParmas)
// 	default:
// 		return inFile(inPath, cmdParmas)
// 	}
// }

// func inURL(inPath string, cmdParmas model.Parmas) error {
// 	req, err := parse.NewURLRequest(inPath, cmdParmas)
// 	if err != nil {
// 		return err
// 	}

// 	bts, err := req.GetBts()
// 	if err != nil {
// 		return err
// 	}

// 	return req.SaveAsFile(bts)
// }

// func inFile(inPath string, cmdParmas model.Parmas) error {
// 	var allowTypes = []string{"tag", "html", "pdf"}
// 	var allow = false

// 	for _, allowType := range allowTypes {
// 		if cmdParmas.OutType == allowType {
// 			allow = true
// 		}
// 	}
// 	if !allow {
// 		return errors.New("非法的输出类型：" + cmdParmas.OutType)
// 	}

// 	req, err := parse.NewFileRequest(inPath, cmdParmas)
// 	if err != nil {
// 		return err
// 	}

// 	bts, err := req.GetBts()
// 	if err != nil {
// 		return err
// 	}

// 	return req.SaveAsFile(bts)
// }

func setConfigLoggerLevel() {
	stdoutLevel, err := tool.TransformToLogrusLevel(config.Obj.Runtime.StdoutLogLevel)
	if err != nil {
		panic(err)
	}
	fileLevel, err := tool.TransformToLogrusLevel(config.Obj.Runtime.FileLogLevel)
	if err != nil {
		panic(err)
	}
	log.SetStdoutLevel(stdoutLevel)
	log.SetFileLevel(fileLevel)
}

func initRootCmdFlags() {
	rootFlags := rootCmd.Flags()
	confParmas := config.Obj.Parmas

	// 添加Flags：变量 长名 短名 默认值 帮助说明
	rootFlags.StringVarP(&cmdParmas.OutPath, "out", "o", confParmas.OutPath, "文件输出的路径")
	rootFlags.StringVarP(&cmdParmas.OutType, "type", "t", confParmas.OutType, "输出的文件类型:tag、html、pdf")
	rootFlags.StringVarP(&cmdParmas.ThemeName, "theme", "e", confParmas.ThemeName, "界面的主题，可放入自定义主题包后修改")
	rootFlags.StringVarP(&cmdParmas.PageFormat, "format", "f", confParmas.PageFormat, "打印的页面格式：A5-A1、Legal、Letter、Tabloid")
	rootFlags.StringVarP(&cmdParmas.PageOrientation, "orientation", "r", confParmas.PageOrientation, "打印的页面方向,可选portrait（纵向）、landscape（横向）")
	rootFlags.StringVarP(&cmdParmas.PageMargin, "margin", "m", confParmas.PageMargin, "打印的页面边距大小，以英寸为单位")
	rootFlags.StringVarP(&cmdParmas.ExecPath, "exec-path", "p", confParmas.ExecPath, "Chrome的执行路径")
	rootFlags.BoolVarP(&cmdParmas.Verbose, "verbose", "v", false, "控制台输出详细日志")
}

func addCmdToRoot() {
	rootCmd.AddCommand(getInstallCmd())
	rootCmd.AddCommand(getConfigCmd())
}

// 输出参数信息调试
func showParams() {
	log.Debugf("---这是你的合计输入参数---")
	log.Debugf("输出路径：%s\n", cmdParmas.OutPath)
	log.Debugf("输出格式：%s\n", cmdParmas.OutType)
	log.Debugf("选择主题：%s\n", cmdParmas.ThemeName)
	log.Debugf("打印页面格式：%s\n", cmdParmas.PageFormat)
	log.Debugf("打印页面方向：%s\n", cmdParmas.PageOrientation)
	log.Debugf("打印页面边距：%s\n", cmdParmas.PageMargin)
	log.Debugf("Chrome的执行路径：%s\n", cmdParmas.ExecPath)
	log.Debugf("--------------------------")
}
