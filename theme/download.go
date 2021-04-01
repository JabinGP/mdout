package theme

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/static"
	"github.com/JabinGP/mdout/tool"
)

// DownloadTheme 从指定 url 下载主题压缩包并解压
func DownloadTheme(themeZipURL string, themeName string, skipTlsVerify bool) error {
	path := filepath.FromSlash(static.ThemeFolderFullName + "/" + themeName)
	if tool.IsExists(path) {
		err := errors.New("主题包：" + path + " 已经存在，如需重新下载请删除后重试！")
		log.Errorln(err)
		return err
	}

	themeZipFileFullName := filepath.FromSlash(static.ThemeFolderFullName + "/" + themeName + ".zip")
	// 打开文件
	file, err := os.Create(themeZipFileFullName)
	if err != nil {
		return err
	}
	defer func() {
		if tool.IsExists(themeZipFileFullName) {
			log.Debugln("开始删除临时压缩文件 " + themeZipFileFullName)
			err := os.Remove(themeZipFileFullName)
			if err != nil {
				log.Errorln("删除文件 " + themeZipFileFullName + "失败，但不会影响使用，可以手动删除！")
				log.Errorln(err)
			}
		}
		log.Debugln("删除 " + themeZipFileFullName + "成功")
	}()
	defer file.Close()

	var urlFileResp *http.Response
	if !skipTlsVerify {
		// 从url获取资源
		urlFileResp, err = http.Get(themeZipURL)
	} else {
		// 跳过https请求时对服务器证书的认证。解决构建镜像时安装主题校验github证书失败问题
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{
			Transport: tr,
		}
		urlFileResp, err = client.Get(themeZipURL)
	}

	if err != nil {
		return err
	}
	if urlFileResp.StatusCode != 200 {
		return fmt.Errorf("从 %s 下载主题包失败，HTTP 状态码 %d。", themeZipURL, urlFileResp.StatusCode)
	}
	defer urlFileResp.Body.Close()

	// 将获取资源复制到文件
	_, err = io.Copy(file, urlFileResp.Body)
	if err != nil {
		return err
	}
	// 复制完成，关闭文件
	file.Close()
	urlFileResp.Body.Close()

	log.Debugln("下载主题压缩包成功，保存到" + themeZipFileFullName)
	log.Debugln("开始解压")

	// 解压缩主题文件
	err = tool.UnZipGithubArchive(themeZipFileFullName, path)
	if err != nil {
		log.Errorln("解压失败")
		log.Errorln(err)
		return err
	}
	log.Debugln("解压成功，解压到" + path)

	log.Infoln("成功下载主题到 " + path)
	return nil
}
