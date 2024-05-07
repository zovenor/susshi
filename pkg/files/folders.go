package files

import (
	"os"
	"path"
)

func CreateFolder(filePath string) error {
	dirPath, _ := path.Split(filePath)
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
