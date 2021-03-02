package tool

import (
	"io/ioutil"
)

func SaveFile(bytes []byte, fileName string) error {
	err := ioutil.WriteFile(fileName, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
