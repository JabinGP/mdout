package static

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/JabinGP/mdout/tool"
)

var (
	Name string

	ConfigVersion        string
	ConfigGithubURL      string
	ConfigFolderName     string
	ConfigFileViperName  string
	ConfigFileType       string
	ConfigFileName       string
	ConfigFolderFullName string
	ConfigFileFullName   string

	ThemeVersion        string
	ThemeFolderName     string
	ThemeFolderFullName string
	ThemeMapGithubURL   string

	LogLevel          string
	LogFolderName     string
	LogFolderFullName string
	LogFileName       string
	LogFileFullName   string
)

func init() {
	Name = "mdout"

	ConfigVersion = "v1"
	ConfigGithubURL = "https://raw.githubusercontent.com/JabinGP/mdout-repo/master/config/{version}/conf.toml"
	ConfigFolderName = "mdout"
	ConfigFileViperName = "conf"
	ConfigFileType = "toml"
	ConfigFileName = ConfigFileViperName + "." + ConfigFileType
	ConfigFolderFullName = tool.GetHomeDir() +
		"/" + ConfigFolderName
	ConfigFileFullName = ConfigFolderFullName +
		"/" + ConfigFileName

	ThemeVersion = "v1"
	ThemeMapGithubURL = "https://raw.githubusercontent.com/JabinGP/mdout-repo/master/theme/{version}/map.toml"
	ThemeFolderName = "theme"
	ThemeFolderFullName = ConfigFolderFullName + "/" + ThemeFolderName

	LogFolderName = "log"
	LogFolderFullName = ConfigFolderFullName + "/" + LogFolderName
	LogFileName = fmt.Sprintf("%d-%d.log", time.Now().Year(), time.Now().Month())
	LogFileFullName = LogFolderFullName + "/" + LogFileName

	fromSlash()
}

func fromSlash() {
	// 将 path 中的 '/' 转换为系统相关的路径分隔符
	ConfigFolderFullName = filepath.FromSlash(ConfigFolderFullName)
	ConfigFileFullName = filepath.FromSlash(ConfigFileFullName)
	ThemeFolderFullName = filepath.FromSlash(ThemeFolderFullName)
	LogFolderFullName = filepath.FromSlash(LogFolderFullName)
	LogFileFullName = filepath.FromSlash(LogFileFullName)
}
