package utils

import "os"

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

func CreatePath(path string) error {
	err := os.MkdirAll(path, os.ModeDir)

	if err != nil {
		return err
	}

	err = os.Chmod(path, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

func CreatePaths(paths... string) error {
	for _, path := range paths {
		err := CreatePath(path)

		if err != nil {
			return err
		}
	}

	return nil
}
