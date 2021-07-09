package directory

import (
	"encoding/json"
	"path/filepath"

	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/pkg/utils"
	"golang.org/x/xerrors"
)

// Layout interface
type Layout interface {
	Save() error
	String() string
}

// ServiceLayout defines the layout interface of a typical service
type ServiceLayout interface {
	Layout

	// GetOutputPath returns root output directory path, e.g.: /...
	GetOutputPath() string
	// GetGoModuleFilePath returns the go.mod file path, e.g.: /.../go.mod
	GetGoModuleFilePath() string
	// GetCommandBasePath returns the cmd root path, e.g.: /.../cmd
	GetCommandBasePath() string
	// GetAPIBasePath returns the API root directory path, e.g.: /.../api
	GetAPIBasePath() string
	// GetInternalBasePath returns the internal root path, e.g.: /.../internal
	GetInternalBasePath() string
	// GetServiceBasePath returns the cmd root path, e.g.: /.../internal/service
	GetServiceBasePath() string
	// GetConfigYAMLBasePath returns the internal root path, e.g.: /.../configs
	GetConfigYAMLBasePath() string
	// GetBuildBasePath returns the cmd root path, e.g.: /.../build
	GetBuildBasePath() string
	// GetDeployBasePath returns the internal root path, e.g.: /.../deploy
	GetDeployBasePath() string
}

// BasicServiceLayout defines the layout of a general service adapting to the standards of https://github.com/golang-standards/project-layout
type BasicServiceLayout struct {
	Options *option.Options

	OutputPath string

	APIBasePath string

	CommandBasePath    string
	InternalBasePath   string
	InternalConfigPath string
	ServiceBashPath    string
	ConfigBasePath     string
	BuildBasePath      string
	DeployBasePath     string
}

// NewBasicServiceLayout returns the basic service directory layout pointer
func NewBasicServiceLayout(opts ...option.Option) *BasicServiceLayout {
	// init options
	options := option.NewOptions(opts...)
	outputPath := options.OutputPath
	return &BasicServiceLayout{
		Options:          options,
		OutputPath:       outputPath,
		APIBasePath:      filepath.Join(outputPath, "api"),
		CommandBasePath:  filepath.Join(outputPath, "cmd"),
		InternalBasePath: filepath.Join(outputPath, "internal"),
		ServiceBashPath:  filepath.Join(outputPath, "internal", "service"),
		ConfigBasePath:   filepath.Join(outputPath, "configs"),
		BuildBasePath:    filepath.Join(outputPath, "build"),
		DeployBasePath:   filepath.Join(outputPath, "deploy"),
	}
}

// String returns the JSON format of directory information
func (layout *BasicServiceLayout) String() string {
	jsonBuf, err := json.Marshal(layout)
	if err != nil {
		return ""
	}
	return string(jsonBuf)
}

// Save creates all the generated directory structure
func (layout *BasicServiceLayout) Save() error {
	// create output directory in case it doesn't exist
	if err := utils.CreatePaths(layout.Options.Override, layout.OutputPath); err != nil {
		return xerrors.Errorf("utils.CreatePaths failed! output path:%v, error:%w", layout.OutputPath, err)
	}

	// create directories
	if err := utils.CreatePaths(
		layout.Options.Override,
		layout.APIBasePath,
		layout.CommandBasePath,
		layout.InternalBasePath,
		layout.ServiceBashPath,
		layout.ConfigBasePath,
		layout.BuildBasePath,
		layout.DeployBasePath); err != nil {
		return xerrors.Errorf("utils.CreatePaths failed! layout:%v, error:%v", layout.String(), err)
	}
	return nil
}

// GetOutputPath returns the output path
func (layout *BasicServiceLayout) GetOutputPath() string {
	return layout.OutputPath
}

// GetGoModuleFilePath returns the go.mod file path
func (layout *BasicServiceLayout) GetGoModuleFilePath() string {
	return filepath.Join(layout.OutputPath, "go.mod")
}

// GetCommandBasePath returns the command base path
func (layout *BasicServiceLayout) GetCommandBasePath() string {
	return layout.CommandBasePath
}

// GetAPIBasePath returns the API base path
func (layout *BasicServiceLayout) GetAPIBasePath() string {
	return layout.APIBasePath
}

// GetInternalBasePath returns the internal base path
func (layout *BasicServiceLayout) GetInternalBasePath() string {
	return layout.InternalBasePath
}

// GetServiceBasePath returns the service base path
func (layout *BasicServiceLayout) GetServiceBasePath() string {
	return layout.ServiceBashPath
}

// GetConfigBasePath returns the config base path
func (layout *BasicServiceLayout) GetConfigYAMLBasePath() string {
	return layout.ConfigBasePath
}

// GetBuildBasePath returns the build base path
func (layout *BasicServiceLayout) GetBuildBasePath() string {
	return layout.BuildBasePath
}

// GetDeployBasePath returns the deploy base path
func (layout *BasicServiceLayout) GetDeployBasePath() string {
	return layout.DeployBasePath
}
