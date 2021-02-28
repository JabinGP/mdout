package model

// Runtime 辅助打印功能的、只能在配置文件中设置的参数
type Runtime struct {
	// 打开配置文件的编辑器命令
	EditorPath string
	// 打开配置文件的编辑器参数
	EditorParmas []string
	// 控制台输出日志等级
	StdoutLogLevel string
	// 文件记录日志等级
	FileLogLevel string
}
