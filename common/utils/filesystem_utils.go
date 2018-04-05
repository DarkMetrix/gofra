package utils

import (
	"os"
	"io/ioutil"
)

func CopyFile(src, dest string) error {
	data, err := ioutil.ReadFile(src)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, data, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

func CheckPathExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func CreatePath(path string, override bool) error {
	isExist, err := CheckPathExists(path)

	if err != nil {
		return err
	}

	if isExist && override {
		err = os.RemoveAll(path)

		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(path, os.ModeDir)

	if err != nil {
		return err
	}

	err = os.Chmod(path, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

func CreatePaths(override bool, paths... string) error {
	for _, path := range paths {
		err := CreatePath(path, override)

		if err != nil {
			return err
		}
	}

	return nil
}
