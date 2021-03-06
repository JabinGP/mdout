package static

var TomlConfig = `[Params]
# 指定使用的主题名称
ThemeName = "github"

# 打印页面选项
PageFormat = "a4"
PageOrientation = "portrait"
PageMargin = "0.4"

# 自定义 Chrome 程序执行路径
ExecPath = ""

# 指定输出类型
OutType = "pdf"
# 指定输出路径
OutPath = ""
# 用于命令行临时指定控制台日志等级为 debug
Verbose = false

[Runtime]
# 使用 ATOM
# EditorPath = "atom"

# 使用 windows自自带记事本 示例
# EditorPath = "notepad"

# 使用 Sublime（在我电脑上的绝对路径，注意替换为自己电脑上的）
# EditorPath = "E:\\Sublime Text 3\\sublime_text.exe"
EditorPath = "code"
EditorParams = []

# 设置输出日志等级，可选列表
# "debug", "info", "error"

# 控制台输出日志等级
StdoutLogLevel = "info"
# 文件记录日志等级
FileLogLevel = "info"

# 参考：https://gitlab.com/golang-commonmark/markdown
# 是否允许 HTML 标签
EnableHTMLTag = true
# 是否开启 XHTMLOutput
EnableXHTMLOutput = true
`
