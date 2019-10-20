package cmd

import (
	// 读写文件

	// 输出日志
	"fmt"
	"log"

	// 错误生成
	"errors"
	// 路径处理
	"path/filepath"
	// 获取系统路径信息
	"os"
	// 命令行框架
	"github.com/spf13/cobra"

	// 读取配置文件
	"github.com/spf13/viper"

	// 命令执行函数
	cmdrun "github.com/JabinGP/mdout/cmdrun"

	// 路径相关工具
	dir "github.com/JabinGP/mdout/dir"

	// 自定义结构体
	types "github.com/JabinGP/mdout/types"
)

var (
	// 输入参数
	cmdParmas types.CmdParmas

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

func initConfig() {
	// 获取文件保存地址和文件名
	confFilePath := dir.HomeDir() + "/binmdout"

	// 配置文件的路径和名字
	viper.SetConfigName("conf")
	viper.AddConfigPath(confFilePath)

	// 设置默认参数
	viper.SetDefault("Out", "")
	viper.SetDefault("Type", "pdf")
	viper.SetDefault("Theme", "default")
	viper.SetDefault("PageFormat", "a4")
	viper.SetDefault("PageOrientation", "portrait")
	viper.SetDefault("PageMargin", "0.4")

	// 读取
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("读取配置文件成功:", viper.ConfigFileUsed())
		showConf()
	} else {
		log.Printf("读取配置文件失败: %s \n", err)
		log.Println("将以系统预设值作为参数默认值")
	}
}

// init，自动初始化
func init() {
	// 加载配置文件
	initConfig()

	// 添加命令
	rootCmd.AddCommand(InstallCmd())

	// 添加Flags：变量 长名 短名 默认值 帮助说明
	rootCmd.Flags().StringVarP(&cmdParmas.Out, "outpath", "o", viper.GetString("Out"), "文件输出的路径")
	rootCmd.Flags().StringVarP(&cmdParmas.Type, "type", "t", viper.GetString("Type"), "输出的文件类型:tag、html、pdf")
	rootCmd.Flags().StringVarP(&cmdParmas.Theme, "theme", "e", viper.GetString("Theme"), "界面的主题，可放入自定义主题包后修改")
	rootCmd.Flags().StringVarP(&cmdParmas.PageFormat, "format", "f", viper.GetString("PageFormat"), "打印的页面格式：A5-A1、Legal、Letter、Tabloid")
	rootCmd.Flags().StringVarP(&cmdParmas.PageOrientation, "orientation", "r", viper.GetString("PageOrientation"), "打印的页面方向,可选portrait（纵向）、landscape（横向）")
	rootCmd.Flags().StringVarP(&cmdParmas.PageMargin, "margin", "m", viper.GetString("PageMargin"), "打印的页面边距大小，以英寸为单位")
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
	log.Println("-----------------------")
}

// 输入读取到的配置文件
func showConf() {
	log.Println("---这是你的配置文件参数---")
	log.Printf("输出路径：%s\n", viper.GetString("Out"))
	log.Printf("输出格式：%s\n", viper.GetString("Type"))
	log.Printf("选择主题：%s\n", viper.GetString("Theme"))
	log.Printf("打印页面格式：%s\n", viper.GetString("PageFormat"))
	log.Printf("打印页面方向：%s\n", viper.GetString("PageOrientation"))
	log.Printf("打印页面边距：%s\n", viper.GetString("PageMargin"))
	log.Println("-----------------------")
}

// 根目录运行函数
func rootRunE(cmd *cobra.Command, args []string) error {

	// 输出参数
	showParams()

	// 获取用户输入路径
	inSource := args[0]

	// 获取输入参数类型
	inType, err := dir.GetType(inSource)
	if err != nil {
		log.Println(err)
		return err
	}

	// 根据不同的输入类型处理，定位到不同的执行函数

	// 输入类型分类
	switch inType {

	case "url": // 如果是url
		return cmdrun.URLToPdfE(inSource, &cmdParmas)

	case "file": // 如果是文件路径
		inExt := filepath.Ext(filepath.Base(inSource)) // 获取文件后缀

		// 输入文件类型分类
		switch inExt {

		case ".html":
			return cmdrun.HTMLToPdfE(inSource, &cmdParmas)

		case ".md":

			// 输出文件类型分类
			switch cmdParmas.Type {

			case "tag":
				return cmdrun.MdToTagE(inSource, &cmdParmas)

			case "html":
				return cmdrun.MdToHTMLE(inSource, &cmdParmas)

			case "pdf":
				return cmdrun.MdToPdfE(inSource, &cmdParmas)

			default:
				return errors.New("未找到md输入匹配的输出文件类型")
			}

		default:
			return errors.New("未找到匹配的输入文件类型")
		}

	default:
		return errors.New("未找到匹配的输入类型")
	}
}
