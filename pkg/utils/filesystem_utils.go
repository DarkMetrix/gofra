package utils

import (
	"errors"
	"io/ioutil"
	"os"
)

// CopyFile copies file from src to dest
func CopyFile(src, dest string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(dest, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// CheckPathExists checks if the path is exist or not
func CheckPathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
}

// CreatePath creates path directory and if override is specified the directory will be removed first
func CreatePath(path string, override bool) error {
	isExist, err := CheckPathExists(path)
	if err != nil {
		return err
	}

	if isExist && override {
		if err = os.RemoveAll(path); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(path, os.ModeDir); err != nil {
		return err
	}

	if err := os.Chmod(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// CreatePaths creates a batch of directories
func CreatePaths(override bool, paths ...string) error {
	for _, path := range paths {
		if err := CreatePath(path, override); err != nil {
			return err
		}
	}

	return nil
}

// GetGOPATH returns the environment variable GOPATH
func GetGOPATH() (string, error) {
	goPath := os.Getenv("GOPATH")

	if len(goPath) == 0 {
		return "", errors.New("GOPATH is not set!")
	} else {
		return goPath, nil
	}
}
