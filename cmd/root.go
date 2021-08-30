package cmd

import (
	"os"

	"github.com/JabinGP/mdout/parser"
	"github.com/JabinGP/mdout/requester"
	"github.com/JabinGP/mdout/static"

	"github.com/JabinGP/mdout/tool"

	"github.com/JabinGP/mdout/config"
	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/model"
	"github.com/spf13/cobra"
)

var (
	// 命令行输入参数，与cobra命令行绑定
	cmdParams model.Params
	// 根命令
	rootCmd = &cobra.Command{
		Use:     "mdout",
		Version: static.Version,
		Short:   "将markdown、html、url转换成pdf",
		Long:    "读取输入的文件，在内部转换成html，并将html渲染为pdf保存",
		Args:    cobra.MinimumNArgs(1),
		RunE:    rootRunE,
	}
)

// init 初始化
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

	// 输出配置文件
	config.ShowConfig()

	// 输出调试参数
	showParams()
	// 构建请求
	req, err := requester.NewRequest(args[0], cmdParams)
	if err != nil {
		return err
	}

	// 执行请求
	err = parser.Parse(req)
	if err != nil {
		return err
	}

	// 保存数据文件
	err = tool.SaveFile(req.Data.([]byte), req.AbsOutPath)
	if err != nil {
		return err
	}

	log.Infof("成功保存文件：%s", req.AbsOutPath)
	return nil
}

func setRuntimeLoggerLevel() {
	if cmdParams.Verbose {
		level, _ := tool.TransformToLogrusLevel("debug")
		log.SetStdoutLevel(level)
	}
}

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
	confParams := config.Obj.Params

	// 添加Flags：变量 长名 短名 默认值 帮助说明
	rootFlags.StringVarP(&cmdParams.OutPath, "out", "o", confParams.OutPath, "文件输出的路径")
	rootFlags.StringVarP(&cmdParams.OutType, "type", "t", confParams.OutType, "输出的文件类型:tag、html、pdf")
	rootFlags.StringVarP(&cmdParams.ThemeName, "theme", "e", confParams.ThemeName, "界面的主题，可放入自定义主题包后修改")
	rootFlags.StringVarP(&cmdParams.PageFormat, "format", "f", confParams.PageFormat, "打印的页面格式：A5-A1、Legal、Letter、Tabloid")
	rootFlags.StringVarP(&cmdParams.PageOrientation, "orientation", "r", confParams.PageOrientation, "打印的页面方向,可选portrait（纵向）、landscape（横向）")
	rootFlags.StringVarP(&cmdParams.PageMargin, "margin", "m", confParams.PageMargin, "打印的页面边距大小，以英寸为单位")
	rootFlags.StringVarP(&cmdParams.ExecPath, "exec-path", "p", confParams.ExecPath, "Chrome的执行路径")
	rootFlags.BoolVarP(&cmdParams.Verbose, "verbose", "v", false, "控制台输出详细日志")
}

func addCmdToRoot() {
	rootCmd.AddCommand(getConfigCmd())
	rootCmd.AddCommand(getInstallCmd())
	rootCmd.AddCommand(getShowCmd())
	rootCmd.AddCommand(getHttpCmd())
}

// 输出参数信息调试
func showParams() {
	log.Debugf("---这是你的合计输入参数---")
	log.Debugf("输出路径：%s\n", cmdParams.OutPath)
	log.Debugf("输出格式：%s\n", cmdParams.OutType)
	log.Debugf("选择主题：%s\n", cmdParams.ThemeName)
	log.Debugf("打印页面格式：%s\n", cmdParams.PageFormat)
	log.Debugf("打印页面方向：%s\n", cmdParams.PageOrientation)
	log.Debugf("打印页面边距：%s\n", cmdParams.PageMargin)
	log.Debugf("Chrome的执行路径：%s\n", cmdParams.ExecPath)
	log.Debugf("--------------------------")
}
