package model

// Parmas 用户输入参数
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
	// 指定Chrome程序执行路径
	ExecPath string
	// 临时设置日志等级为"debug"
	Verbose bool
}
