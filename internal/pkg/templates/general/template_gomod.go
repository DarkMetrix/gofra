package general

import (
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"golang.org/x/xerrors"

	"github.com/DarkMetrix/gofra/internal/pkg/templates"
)

// GoModuleInfo represents go module information
type GoModuleInfo struct {
	Opts      *option.Options
	GoModule  string
	GoVersion string
}

// NewGoModuleInfo returns a new GoModuleInfo pointer
func NewGoModuleInfo(opts ...option.Option) *GoModuleInfo {
	// init options
	newOpts := option.NewOptions(opts...)
	return &GoModuleInfo{
		Opts:      newOpts,
		GoModule:  newOpts.GoModule,
		GoVersion: newOpts.GoVersion,
	}
}

// RenderFile render template and output to file
func (info *GoModuleInfo) RenderFile(outputPath string) error {
	if err := templates.RenderToFile(outputPath, info.Opts.Override, info.Opts.IgnoreExist,
		"template-go-module", GoModuleTemplate, info); err != nil {
		return xerrors.Errorf("RenderToFile failed! error:%w", err)
	}
	return nil
}

var GoModuleTemplate string = `module {{.GoModule}} 

go {{.GoVersion}}
`
