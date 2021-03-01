package model

// Parmas 可供用户在命令行设置的参数
type Parmas struct {
	// 输入输出参数
	IOParmas
	// 打印页面参数
	PageParmas
	// 主题
	ThemeName string
	// 指定 Chrome 程序执行路径
	ExecPath string
	// 临时设置控制台日志等级为 "debug"
	Verbose bool
}

// IOParmas 输入输出参数
type IOParmas struct {
	// 文件输出路径
	OutPath string
	// 目标文件类型
	OutType string
}

// PageParmas 打印页面参数
type PageParmas struct {
	// 打印页面格式
	PageFormat string
	// 打印页面方向
	PageOrientation string
	// 打印边距
	PageMargin string
}
