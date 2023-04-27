package cmd

import (
	"fmt"
	"regexp"
	"time"

	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/parser"
	"github.com/JabinGP/mdout/requester"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func getHttpCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "http",
		Short: "接口使用",
		Long: "接口使用",
		Args: cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			setRuntimeLoggerLevel()

			gin.SetMode(gin.ReleaseMode)
			log.Infoln("开始监听9090端口")
			router := gin.Default()
			router.MaxMultipartMemory = 100 << 20
			router.POST("/md", func(c *gin.Context) {
				str := c.PostForm("markdown")
				reg := regexp.MustCompile(`:::.*`)
				str = reg.ReplaceAllString(str, ``)
				// 构建请求
				req := requester.NewRequestFromData([]byte(str), cmdParams)
				// log.Infoln(req)
				// 执行请求
				err := parser.Parse(req)
				if err != nil {
					log.Errorln(err.Error())
					return
				}

				c.Header("content-disposition", fmt.Sprintf("attachment; filename=api_%d.pdf", rune(time.Now().Unix())))
				c.Data(200, "application/pdf", req.Data.([]byte))
			})
			router.POST("/md/file", func(c *gin.Context) {
				// 单文件
				file, err := c.FormFile("file")
				if err != nil {
					log.Errorln("file error")
				}
				fo, err := file.Open()
				if err != nil {
					log.Errorln("file error")
				}
				defer fo.Close()
				buf := make([]byte, file.Size)
				fo.Read(buf)
				str := string(buf)
				reg := regexp.MustCompile(`:::.*`)
				str = reg.ReplaceAllString(str, ``)
				// 构建请求
				req := requester.NewRequestFromData([]byte(str), cmdParams)
				// log.Infoln(req)
				// 执行请求
				err = parser.Parse(req)
				if err != nil {
					log.Errorln(err.Error())
					return
				}

				c.Header("content-disposition", `attachment; filename=api_doc_`+string(rune(time.Now().Unix()))+`.pdf`)
				c.Data(200, "application/pdf", req.Data.([]byte))
			})
			router.Run(":9090")
			return nil
		},
	}
	return cmd
}
