package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/JabinGP/mdout/config"
	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/static"
	"github.com/spf13/cobra"
)

func getShowCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "show",
		Short: "输出信息",
		Long:  "输出指定信息",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			in := args[0]
			switch in {
			case "config-folder":
				fmt.Println(static.ConfigFolderFullName)
			case "theme-folder":
				fmt.Println(static.ThemeFolderFullName)
			case "log-folder":
				fmt.Println(static.LogFolderFullName)
			case "config-info":
				confBytes, err := json.MarshalIndent(config.Obj, "", "\t")
				if err != nil {
					log.Errorln(err)
					return err
				}
				fmt.Printf(string(confBytes))
			}
			return nil
		},
	}

	cmd.AddCommand(getInstallThemeCmd())
	return cmd
}
