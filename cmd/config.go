package cmd

import (
	"io/ioutil"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/JabinGP/mdout/config"
	"github.com/JabinGP/mdout/tool"
)

func getConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "修改配置",
		Long:  "通过编辑器修改配置文件，默认打开vscode",
		RunE:  configRunE,
	}
}

func configRunE(cmd *cobra.Command, args []string) error {
	runtime := config.Obj.Runtime

	if runtime.EditorPath == "" {
		log.Debugln("未设置编辑器，将使用vscode打开配置文件")
		runtime.EditorPath = "code"
	}
	configFullName := tool.GetHomeDir() + "/mdout/conf.toml"
	execCmd := exec.Command(runtime.EditorPath, configFullName)
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		log.Errorln(err)
		return err
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := execCmd.Start(); err != nil {
		log.Errorln(err)
		return err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Errorln(err)
		return err
	}
	opString := string(opBytes)
	if opString != "" {
		log.Infoln(opString)
	}
	return nil
}
