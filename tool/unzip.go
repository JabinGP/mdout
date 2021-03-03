package tool

import (
	"archive/zip"
	"io"
	"os"
	"strings"
)

// UnZipGithubArchive 解压缩 github 归档zip文件到指定目录。
// 并且忽略 github 归档自动生成的外层文件夹。
func UnZipGithubArchive(zipFile string, dest string) (err error) {
	//目标文件夹不存在则创建
	if _, err = os.Stat(dest); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dest, 0755)
		}
	}

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}

	defer reader.Close()

	// 第一个 file 是最外层文件夹，记录下来
	outerFolderName := reader.File[0].Name

	// 处理文件名，将所有文件的名字去掉开头的最外层文件夹的名字
	for _, file := range reader.File[1:] {
		file.Name = strings.Replace(file.Name, outerFolderName, "", 1)
	}

	// 跳过最外层文件夹，照旧处理
	for _, file := range reader.File[1:] {
		// fmt.Println(file.FileInfo().IsDir(), file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(dest+"/"+file.Name, 0755)
			continue
		} else {
			err = os.MkdirAll(getDir(dest+"/"+file.Name), 0755)
			if err != nil {
				return err
			}
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		filename := dest + "/" + file.Name

		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()

		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
	}
	return
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
