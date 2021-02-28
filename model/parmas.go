package model

// Parmas 与转换直接相关的、可供用户在命令行设置的参数
type Parmas struct {
	// 文件输出路径
	OutPath string
	// 目标文件类型
	OutType string
	// 主题
	ThemeName string
	// 打印页面格式
	PageFormat string
	// 打印页面方向
	PageOrientation string
	// 打印边距
	PageMargin string
	// 指定 Chrome 程序执行路径
	ExecPath string
	// 临时设置控制台日志等级为 "debug"
	Verbose bool
}
