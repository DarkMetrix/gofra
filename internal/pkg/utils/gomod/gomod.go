package gomod

import (
	"io/ioutil"

	"golang.org/x/mod/modfile"
	"golang.org/x/xerrors"
)

// GetGoModule get go module from go.mod
func GetGoModule(goModFilePath string) (string, error) {
	// read the go module file contents
	raw, err := ioutil.ReadFile(goModFilePath)
	if err != nil {
		return "", xerrors.Errorf("ioutil.ReadFile failed! error:%w", err)
	}

	goModule := modfile.ModulePath(raw)
	if goModule == "" {
		return "", xerrors.New("go module not found from go.mod")
	}

	return goModule, nil
}
