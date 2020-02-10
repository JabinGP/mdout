package theme

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/JabinGP/mdout/config"
	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/static"
	"github.com/JabinGP/mdout/tool"
	"github.com/JabinGP/mdout/ziper"
	"github.com/spf13/viper"
)

// DownloadTheme 从github仓库下载主题文件
func DownloadTheme(themeName string) error {
	theme, err := GetTheme(themeName)
	if err != nil {
		log.Println("获取主题信息失败！")
		return err
	}
	path := filepath.FromSlash(static.ThemeFolderFullName + "/" + themeName)
	switch theme.DownloadType {
	case "zip":
		err = DownloadThemeByZip(theme, path)
	case "git":
		err = DownloadThemeByGit(theme, path)
	default:
		err = DownloadThemeByGit(theme, path)
	}
	if err != nil {
		log.Println("获取主题失败！")
		return err
	}

	log.Println("获取主题成功，保存在 " + path)
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
		log.Println(err)
		return err
	}
	defer file.Close()
	defer func() {
		log.Println("开始删除临时压缩文件 " + themeZipFileFullName)

		if tool.IsExists(themeZipFileFullName) {
			log.Println("文件存在，正在删除 " + themeZipFileFullName)
			err := os.Remove(themeZipFileFullName)
			if err != nil {
				log.Println("删除文件 " + themeZipFileFullName + "失败，但不会影响使用，可以手动删除")
				log.Println(err)
			}
		}
		log.Println("删除 " + themeZipFileFullName + "成功")
	}()

	// 从url获取资源
	urlFileResp, err := http.Get(themeZipURL)
	if err != nil {
		log.Println(err)
		return err
	}
	defer urlFileResp.Body.Close()

	// 将获取资源复制到文件
	_, err = io.Copy(file, urlFileResp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	// 复制完成，关闭文件
	file.Close()
	urlFileResp.Body.Close()

	log.Println("下载主题压缩包成功，保存到" + themeZipFileFullName)
	log.Println("开始解压")

	// 解压缩主题文件
	err = ziper.UnZip(themeZipFileFullName, path)
	if err != nil {
		log.Println("解压失败")
		log.Println(err)
		return err
	}
	log.Println("解压成功，解压到" + path)

	log.Println("成功下载 " + themeName + " 主题到 " + path)
	return nil
}

// GetTheme 获取主题信息
func GetTheme(themeName string) (model.Theme, error) {
	var targetTheme model.Theme
	themeList, err := DownloadThemeMap()
	if err != nil {
		log.Println("获取主题列表失败!")
		return targetTheme, err
	}
	for _, theme := range themeList {
		if theme.Name == themeName {
			targetTheme = theme
			break
		}
	}
	if targetTheme.Name == "" {
		log.Println("未找到对应名称的主题！")
		return targetTheme, errors.New("未找到对应名称的主题！")
	}

	return targetTheme, nil
}

// DownloadThemeMap 从github获取所有主题的下载地址
func DownloadThemeMap() ([]model.Theme, error) {
	log.Println("正在获取主题列表...")
	var mapURL = static.ThemeMapGithubURL
	var version = static.ThemeVersion
	mapURL = strings.Replace(mapURL, "{version}", version, 1)

	mapReso, err := http.Get(mapURL)
	if err != nil {
		log.Println("对 " + mapURL + " 发起http请求失败！")
		return nil, err
	}
	defer mapReso.Body.Close()
	v := viper.New()
	v.SetConfigType("toml")
	err = v.ReadConfig(mapReso.Body)
	if err != nil {
		log.Println("从 " + mapURL + " 读取http响应body失败！")
		return nil, err
	}
	var mapList []model.Theme
	v.UnmarshalKey("Theme", &mapList)
	log.Println("从 " + mapURL + " 获取主题列表成功")
	return mapList, nil
}

// DownloadThemeByGit 使用git clone命令下载主题
func DownloadThemeByGit(theme model.Theme, path string) error {
	//themeFolderFullName := filepath.FromSlash(static.ThemeFolderFullName + "/" + theme.Name + "-test")
	cloneCmd := exec.Command(config.Obj.GitPath, "clone", theme.RepoURL, path)
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cloneCmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return err
	}
	// 保证关闭输出流
	defer stdout.Close()

	// 运行命令
	log.Println("执行命令 " + config.Obj.GitPath + " clone " + theme.RepoURL + " " + path)
	if err := cloneCmd.Start(); err != nil {
		log.Println(err)
		return err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println(err)
		return err
	}
	opString := string(opBytes)
	if opString != "" {
		log.Println(opString)
	}
	return nil
}
