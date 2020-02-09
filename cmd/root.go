package cmd

import (
	// 日志
	"log"

	// 获取系统路径信息
	"os"
	// 命令行框架
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// 命令执行函数

	// 自定义结构体
	"github.com/JabinGP/mdout/cmdrun"
	"github.com/JabinGP/mdout/model"

	// 配置文件读取
	"github.com/JabinGP/mdout/config"
)

var (
	v *viper.Viper
	// 命令行输入参数，与cobra命令行绑定
	cmdParmas model.Parmas
	// 配置文件参数
	conf model.Config
	// 根命令
	rootCmd = &cobra.Command{
		Use:     "mdout",
		Version: "0.5",
		Short:   "将markdown、html、url转换成pdf",
		Long:    "读取输入的文件，在内部转换成html，并将html渲染为pdf保存",
		Args:    cobra.MinimumNArgs(1),
		RunE:    rootRunE,
	}
)

// Execute 程序执行入口
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// 根目录运行函数
func rootRunE(cmd *cobra.Command, args []string) error {
	// 输出参数
	showParams()
	// 获取用户输入路径
	in := args[0]
	return cmdrun.Distribute(in, cmdParmas)
}

// init 包初始化
func init() {
	initConfig()
	initRootFlags()
	addCommand()
}

func initConfig() {
	v = config.NewViper()
	v.Unmarshal(&conf)
}

func initRootFlags() {
	rootFlags := rootCmd.Flags()
	confParmas := conf.Parmas

	// 添加Flags：变量 长名 短名 默认值 帮助说明
	rootFlags.StringVarP(&cmdParmas.Out, "outpath", "o", confParmas.Out, "文件输出的路径")
	rootFlags.StringVarP(&cmdParmas.Type, "type", "t", confParmas.Type, "输出的文件类型:tag、html、pdf")
	rootFlags.StringVarP(&cmdParmas.Theme, "theme", "e", confParmas.Theme, "界面的主题，可放入自定义主题包后修改")
	rootFlags.StringVarP(&cmdParmas.PageFormat, "format", "f", confParmas.PageFormat, "打印的页面格式：A5-A1、Legal、Letter、Tabloid")
	rootFlags.StringVarP(&cmdParmas.PageOrientation, "orientation", "r", confParmas.PageOrientation, "打印的页面方向,可选portrait（纵向）、landscape（横向）")
	rootFlags.StringVarP(&cmdParmas.PageMargin, "margin", "m", confParmas.PageMargin, "打印的页面边距大小，以英寸为单位")
	rootFlags.StringVarP(&cmdParmas.ExecPath, "exec-path", "p", confParmas.ExecPath, "Chrome的执行路径")
}

func addCommand() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(configCmd)
}

// 输出参数信息调试
func showParams() {
	log.Println("---这是你的合计输入参数---")
	log.Printf("输出路径：%s\n", cmdParmas.Out)
	log.Printf("输出格式：%s\n", cmdParmas.Type)
	log.Printf("选择主题：%s\n", cmdParmas.Theme)
	log.Printf("打印页面格式：%s\n", cmdParmas.PageFormat)
	log.Printf("打印页面方向：%s\n", cmdParmas.PageOrientation)
	log.Printf("打印页面边距：%s\n", cmdParmas.PageMargin)
	log.Printf("Chrome的执行路径：%s\n", cmdParmas.ExecPath)
	log.Println("--------------------------")
}
