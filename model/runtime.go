package model

// Runtime 只能在配置文件中设置的参数
type Runtime struct {
	LogRuntime
	EditorRuntime
	MDRuntime
}

// EditorRuntime 编辑器相关的参数
type EditorRuntime struct {
	// 打开配置文件的编辑器命令
	EditorPath string
	// 打开配置文件的编辑器参数
	EditorParmas []string
}

// LogRuntime 日志参数
type LogRuntime struct {
	// 控制台输出日志等级
	StdoutLogLevel string
	// 文件记录日志等级
	FileLogLevel string
}

// MDRuntime 与 markdown 解析时相关的参数
type MDRuntime struct {
	// 参考：https://gitlab.com/golang-commonmark/markdown
	// 是否允许 HTML 标签
	EnableHTMLTag bool
	// 是否开启 XHTMLOutput
	EnableXHTMLOutput bool
}
