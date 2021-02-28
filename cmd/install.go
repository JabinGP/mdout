package cmd

import (
	"github.com/JabinGP/mdout/config"
	"github.com/JabinGP/mdout/static"
	"github.com/spf13/cobra"
)

var installParmas = struct {
	URL string
}{}

func getInstallCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "install",
		Short: "初始化配置文件",
		Long:  "到指定的地址下载配置文件",
		RunE:  installRunE,
	}

	initInstallCmdFlags(cmd)
	return cmd
}

func installRunE(cmd *cobra.Command, args []string) error {
	return config.DownloadConfig(installParmas.URL)
}

func initInstallCmdFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	// 添加 Flags
	flags.StringVarP(&installParmas.URL, "url", "u", static.ConfigURL, "toml配置文件地址")
}
