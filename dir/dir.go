package dir

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// HomeDir 返回用户家目录
func HomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	if v := os.Getenv("HOME"); v != "" {
		return v
	}
	if u, err := user.Current(); err == nil {
		return u.HomeDir
	}
	return ""
}

// IsExists 判断所给路径文件或文件夹是否存在
func IsExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// GetType 分析输入的路径，返回类型结果
func GetType(path string) (string, error) {

	// 检查输入是否为url
	urlRegexp := `(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`
	match, err := regexp.MatchString(urlRegexp, path)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// 如果输入符合url正则表达式
	if match {
		return "url", nil
	}

	// 将剩余情况匹配为文件路径
	return "file", nil
}

// GetDirNameExt 输入文件路径，返回文件夹、文件名、文件拓展名（绝对路径）
func GetDirNameExt(path string) (string, string, string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}
	dir := filepath.Dir(absPath)
	ext := filepath.Ext(filepath.Base(absPath))
	name := strings.TrimSuffix(
		filepath.Base(absPath),
		filepath.Ext(filepath.Base(absPath)),
	)
	return dir, name, ext, nil
}

// IsDir 返回所给路径是否为一个文件夹
func IsDir(path string) (bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Println(err)
		return false, err
	}
	fi, err := os.Stat(absPath)
	if err != nil {
		return false, nil
	}
	return fi.IsDir(), nil
}
