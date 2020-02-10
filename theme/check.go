package theme

import (
	"path/filepath"

	"github.com/JabinGP/mdout/static"
	"github.com/JabinGP/mdout/tool"
)

// CheckTheme 检测主题相关的文件夹是否存在
func CheckTheme(themeName string)bool{
	themeFolderFullName := filepath.FromSlash(static.ThemeFolderFullName+"/"+themeName)
	return tool.IsExists(themeFolderFullName)
}