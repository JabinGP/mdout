package model

import "html/template"

// TemplateData 替换主题文件 index.html 中的变量
type TemplateData struct {
	MarkdownBody template.HTML
	ThemePath    template.URL
}
