package model

// Config 配置文件对应实体
type Config struct {
	Parmas
	Runtime
}

// Parmas 用户输入参数
type Parmas struct {
	// 文件输出路径
	Out string
	// 目标文件类型
	Type string
	// 主题
	Theme string
	// 打印页面格式
	PageFormat string
	// 打印页面方向
	PageOrientation string
	// 打印边距
	PageMargin string
	// 指定Chrome程序执行路径
	ExecPath string
}

// Runtime ...
type Runtime struct {
	// 打开配置文件的编辑器路径或命令
	EditorPath string
}
