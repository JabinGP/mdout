package static

import (
	"path/filepath"

	"github.com/JabinGP/mdout/tool"
)

var (
	Name                 string
	ConfigFolderName     string
	ConfigFileViperName  string
	ConfigFileType       string
	ConfigFileName       string
	ConfigFolderFullName string
	ConfigFileFullName   string
	ThemeFolderName      string
	ThemeFolderFullName  string
	ConfigVersion        string
	ThemeVersion         string
	ConfigGithubURL      string
	ThemeMapGithubURL    string
)

func init() {
	Name = "mdout"
	ConfigFolderName = "mdout"
	ConfigFileViperName = "conf"
	ConfigFileType = "toml"
	ConfigFileName = ConfigFileViperName + "." + ConfigFileType
	ConfigFolderFullName = tool.GetHomeDir() +
		"/" + ConfigFolderName
	ConfigFileFullName = ConfigFolderFullName +
		"/" + ConfigFileName
	ThemeFolderName = "theme"
	ThemeFolderFullName = ConfigFolderFullName + "/" + ThemeFolderName
	ConfigVersion = "v1"
	ThemeVersion = "v1"
	ConfigGithubURL = "https://raw.githubusercontent.com/JabinGP/mdout-repo/master/config/{version}/conf.toml"
	ThemeMapGithubURL = "https://raw.githubusercontent.com/JabinGP/mdout-repo/master/theme/{version}/map.toml"
	// 将 path 中的 '/' 转换为系统相关的路径分隔符
	ConfigFolderFullName = filepath.FromSlash(ConfigFolderFullName)
	ConfigFileFullName = filepath.FromSlash(ConfigFileFullName)
	ThemeFolderFullName = filepath.FromSlash(ThemeFolderFullName)
}
