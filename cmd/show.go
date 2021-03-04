package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

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
			case "theme-list":
				fileInfoList, err := ioutil.ReadDir(static.ThemeFolderFullName)
				if err != nil {
					log.Errorln(err)
				}
				themeList := []string{}
				for i := range fileInfoList {
					if fileInfoList[i].IsDir() {
						themeList = append(themeList, fileInfoList[i].Name())
					}
				}

				if len(themeList) == 0 {
					fmt.Println("暂无主题")
					return nil
				}

				for _, themeName := range themeList {
					fmt.Println(themeName)
				}

				return nil
			default:
				err := fmt.Errorf("无法识别的 show 命令输入 %s，请检查输入。", in)
				log.Errorln(err)
				return err
			}

			return nil
		},
	}

	cmd.AddCommand(getInstallThemeCmd())
	return cmd
}
