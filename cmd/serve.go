package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func getServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "实时预览",
		Long:  "对文件启动实时预览服务器，监听文件更改自动刷新",
		RunE:  serveRunE,
	}
}
func serveRunE(cmd *cobra.Command, args []string) error {
	httpServe()
	return nil
}

func httpServe() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(rw, "hello World")
	})
	http.ListenAndServe(":8888", nil)
}
