package cmd

import (
	"fmt"

	"github.com/JabinGP/mdout/theme"
	"github.com/spf13/cobra"
)

var installParams = struct {
	URL string
}{}

func getInstallCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "install",
		Short: "下载资源",
		Long:  "指定路径下载配置文件或主题包",
	}

	cmd.AddCommand(getInstallThemeCmd())
	return cmd
}

func getInstallThemeCmd() *cobra.Command {
	var url, name string
	var cmd = &cobra.Command{
		Use:   "theme",
		Short: "下载主题",
		Long:  "到指定的地址下载主题",
		RunE: func(cmd *cobra.Command, args []string) error {
			if url == "" || name == "" {
				return fmt.Errorf("url = %s，name = %s，指定的 url 或者 name 为空！", url, name)
			}
			return theme.DownloadTheme(url, name)
		},
	}

	cmd.Flags().StringVarP(&url, "url", "u", "", "主题文件zip包的地址")
	cmd.Flags().StringVarP(&name, "name", "n", "", "主题文件保存的文件夹名")
	return cmd
}
