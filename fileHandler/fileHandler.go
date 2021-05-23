package fileHandler

import (
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func CheckCreateFile(path string) (string, error) {

	path = expandPath(path)
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return "", err
		}
		defer file.Close()
	} else if err != nil {
		return "", err
	}

	fileInfo, err = os.Stat(path)
	if err != nil {
		return "", err
	}
	mode := fileInfo.Mode()
	if !mode.IsRegular() {
		return "", errors.New("File already exists. File is not a regular file.")
	}
	return path, nil
}

func ReadFromFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	jsonData, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func WriteToFile(jsonData []byte, path string) error {
	var f *os.File
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(jsonData)
	return err
}

// ExpandPath is helper function to expand file location
func expandPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if strings.HasPrefix(path, "~/") {
		user, err := user.Current()
		if err != nil {
			return path
		}
		return filepath.Join(user.HomeDir, path[2:])
	}
	abspath, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abspath
}
