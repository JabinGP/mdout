package tool

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/JabinGP/mdout/model"
)

// GetHomeDir 返回用户家目录
func GetHomeDir() string {
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
	absPath, err := Abs(path)
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
func IsDir(path string) bool {
	absPath, err := Abs(path)
	if err != nil {
		log.Println(err)
		return false
	}
	fi, err := os.Stat(absPath)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// GetOutFullNameFromIn ...
func GetOutFullNameFromIn(in string, option func(outDir *string, outName *string, outExt *string)) (string, error) {
	absIn, err := Abs(in)
	if err != nil {
		log.Println(err)
		return "", err
	}

	outDir, outName, outExt, err := GetDirNameExt(absIn)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if option != nil {
		option(&outDir, &outName, &outExt)
	}
	return Abs(outDir + "/" + outName + outExt)
}

// GetOutFullNameFromOut ...
func GetOutFullNameFromOut(out string, option func(outDir *string, outName *string, outExt *string)) (string, error) {
	absOut, err := Abs(out)
	if err != nil {
		log.Println(err)
		return "", err
	}

	outDir, outName, outExt, err := GetDirNameExt(absOut)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if option != nil {
		option(&outDir, &outName, &outExt)
	}

	return Abs(outDir + "/" + outName + outExt)
}

// GetOutFullName ...
func GetOutFullName(in string, parmas model.Parmas) (string, error) {
	var outFullName string

	absOut, err := Abs(parmas.Out)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if parmas.Out == "" { // 未指定输出位置
		outFullName, err = GetOutFullNameFromIn(in, func(outDir, outName, outExt *string) {
			*outExt = "." + parmas.Type
		})
	} else {
		if IsDir(absOut) { // 指定输出到文件夹
			outFullName, err = GetOutFullNameFromIn(in, func(outDir, outName, outExt *string) {
				*outDir = absOut
				*outExt = "." + parmas.Type
			})
		} else { // 指定输出到文件名
			outFullName, err = GetOutFullNameFromOut(parmas.Out, nil)
		}
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	return outFullName, nil
}

// Abs 将开头的~替换为家目录，再进行绝对值化
func Abs(path string) (string, error) {
	if path == "" {
		return filepath.Abs(path)
	}
	var tmpPath string
	if path[0] == '~' {
		tmpPath = GetHomeDir() + path[1:]
	} else {
		tmpPath = path
	}
	return filepath.Abs(tmpPath)
}
