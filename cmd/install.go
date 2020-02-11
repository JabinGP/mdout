package cmd

import (
	"github.com/JabinGP/mdout/theme"
	"github.com/spf13/cobra"
)

func getInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "下载运行时必要的资源到本地",
		Long:  "下载配置文件，主题包并保存到用户家目录下的binmdout文件夹",
		RunE:  installRunE,
	}
}

// installRunE install命令执行时运行
func installRunE(cmd *cobra.Command, args []string) error {
	themeName := args[0]
	return theme.DownloadTheme(themeName)
}
