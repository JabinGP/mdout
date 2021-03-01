package static

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/JabinGP/mdout/tool"
)

var (
	Name    string
	Version string

	ConfigFolderName     string
	ConfigURL            string
	ConfigFileName       string
	ConfigFolderFullName string
	ConfigFileFullName   string

	ThemeFolderName     string
	ThemeFolderFullName string

	LogLevel          string
	LogFolderName     string
	LogFolderFullName string
	LogFileName       string
	LogFileFullName   string
)

func init() {
	Name = "mdout"
	Version = "v0.6"

	ConfigFolderName = "mdout"
	ConfigURL = fmt.Sprintf("https://raw.githubusercontent.com/JabinGP/mdout-repo/master/config/%s/conf.toml", Version)
	ConfigFileName = "conf.toml"
	ConfigFolderFullName = tool.GetHomeDir() +
		"/" + ConfigFolderName
	ConfigFolderFullName = filepath.FromSlash(ConfigFolderFullName)
	ConfigFileFullName = ConfigFolderFullName +
		"/" + ConfigFileName
	ConfigFileFullName = filepath.FromSlash(ConfigFileFullName)

	ThemeFolderName = "theme"
	ThemeFolderFullName = ConfigFolderFullName + "/" + ThemeFolderName
	ThemeFolderFullName = filepath.FromSlash(ThemeFolderFullName)

	LogFolderName = "log"
	LogFolderFullName = ConfigFolderFullName + "/" + LogFolderName
	LogFolderFullName = filepath.FromSlash(LogFolderFullName)
	LogFileName = fmt.Sprintf("%d-%d.log", time.Now().Year(), time.Now().Month())
	LogFileFullName = LogFolderFullName + "/" + LogFileName
	LogFileFullName = filepath.FromSlash(LogFileFullName)
}
