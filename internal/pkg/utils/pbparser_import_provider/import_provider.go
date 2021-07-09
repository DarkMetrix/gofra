package pbparser_import_provider

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type GofraImportModuleProvider struct {
	ModuleSearchPaths []string
}

func (p *GofraImportModuleProvider) Provide(module string) (io.Reader, error) {
	for _, path := range p.ModuleSearchPaths {
		modulePath := path + string(filepath.Separator) + module

		// read the module file contents & create a reader...
		raw, err := ioutil.ReadFile(modulePath)
		if err != nil {
			continue
		}

		r := strings.NewReader(string(raw[:]))
		return r, nil
	}
	return nil, errors.New(fmt.Sprintf("No module file found in these paths:%v, module:%v", p.ModuleSearchPaths, module))
}
