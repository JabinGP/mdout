package theme

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/JabinGP/mdout/config"
	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/static"
	"github.com/JabinGP/mdout/tool"
	"github.com/JabinGP/mdout/ziper"
)

// DownloadTheme 从github仓库下载主题文件
func DownloadTheme(themeName string) error {
	path := filepath.FromSlash(static.ThemeFolderFullName + "/" + themeName)
	if tool.IsExists(path) {
		err := errors.New("主题包：" + path + " 已经存在，如需重新下载请删除后重试！")
		log.Errorln(err)
		return err
	}
	theme, err := GetTheme(themeName)
	if err != nil {
		return err
	}
	switch theme.DownloadType {
	case "zip":
		err = DownloadThemeByZip(theme, path)
	case "git":
		err = DownloadThemeByGit(theme, path)
	default:
		err = DownloadThemeByGit(theme, path)
	}
	if err != nil {
		log.Errorln(err)
		return err
	}

	log.Infoln("下载主题成功，保存在 " + path)
	return nil
}

// DownloadThemeByZip 从github仓库下载主题文件
func DownloadThemeByZip(theme model.Theme, path string) error {
	themeName := theme.Name
	themeZipURL := theme.ZipURL
	themeZipFileFullName := filepath.FromSlash(static.ThemeFolderFullName + "/" + themeName + ".zip")
	// 打开文件
	file, err := os.Create(themeZipFileFullName)
	if err != nil {
		return err
	}
	defer file.Close()
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

	// 从url获取资源
	urlFileResp, err := http.Get(themeZipURL)
	if err != nil {
		return err
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
	err = ziper.UnZip(themeZipFileFullName, path)
	if err != nil {
		log.Errorln("解压失败")
		log.Errorln(err)
		return err
	}
	log.Debugln("解压成功，解压到" + path)

	log.Infoln("成功下载 " + themeName + " 主题到 " + path)
	return nil
}

// GetTheme 获取主题信息
func GetTheme(themeName string) (model.Theme, error) {
	var targetTheme model.Theme
	themeList, err := DownloadThemeMap()
	if err != nil {
		return targetTheme, err
	}
	for _, theme := range themeList {
		if theme.Name == themeName {
			targetTheme = theme
			break
		}
	}
	if targetTheme.Name == "" {
		return targetTheme, errors.New("下载失败，未在主题列表中找到主题：" + themeName)
	}

	return targetTheme, nil
}

// DownloadThemeMap 从github获取所有主题的下载地址
func DownloadThemeMap() ([]model.Theme, error) {
	log.Debugln("正在获取主题列表...")
	var themeMap model.ThemeMap
	var mapURL = static.ThemeMapGithubURL
	var version = static.ThemeVersion
	mapURL = strings.Replace(mapURL, "{version}", version, 1)

	mapReso, err := http.Get(mapURL)
	if err != nil {
		return nil, err
	}
	defer mapReso.Body.Close()
	if _, err := toml.DecodeReader(mapReso.Body, &themeMap); err != nil {
		return nil, err
	}
	log.Debugln("从 " + mapURL + " 获取主题列表成功")
	return themeMap.ThemeList, nil
}

// DownloadThemeByGit 使用git clone命令下载主题
func DownloadThemeByGit(theme model.Theme, path string) error {
	runtime := config.Obj.Runtime
	gitPath := runtime.GitPath
	cloneCmd := exec.Command(gitPath, "clone", theme.RepoURL, path)

	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cloneCmd.StdoutPipe()
	if err != nil {
		return err
	}
	// 保证关闭输出流
	defer stdout.Close()
	defer func() {
		// 读取输出结果
		opBytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			return
		}
		opString := string(opBytes)
		if opString != "" {
			log.Infoln(opString)
		}
	}()
	// 运行命令
	log.Debugln("执行命令 " + gitPath + " clone " + theme.RepoURL + " " + path)
	if err := cloneCmd.Run(); err != nil {
		return err
	}
	return nil
}
