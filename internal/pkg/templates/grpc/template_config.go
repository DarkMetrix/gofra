package grpc

import (
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"golang.org/x/xerrors"

	"github.com/DarkMetrix/gofra/internal/pkg/templates"
)

// ConfigInfo represents the Config information
type ConfigInfo struct {
	Opts *option.Options
}

// NewConfigInfo returns a new ConfigInfo pointer
func NewConfigInfo(opts ...option.Option) *ConfigInfo {
	// init options
	newOpts := option.NewOptions(opts...)
	return &ConfigInfo{Opts: newOpts}
}

// RenderFile render template and output to file
func (config *ConfigInfo) RenderFile(outputPath string) error {
	if err := templates.RenderToFile(outputPath, config.Opts.Override, config.Opts.IgnoreExist,
		"template-config", ConfigTemplate, config); err != nil {
		return xerrors.Errorf("RenderToFile failed! error:%w", err)
	}
	return nil
}

// ConfigTemplate defines service configuration in YAML
var ConfigTemplate string = `package config

import (
	"encoding/json"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

var globalConfig *Config

// Config definition
type Config struct {
	Server ServerInfo "mapstructure:\"server\" json:\"server\""
	Observability ObservabilityInfo "mapstructure:\"observability\" json:\"observability\""
	Pprof PprofInfo "mapstructure:\"pprof\" json:\"pprof\""
	Performance PerformanceInfo "mapstructure:\"performance\" json:\"performance\""
}

// ServerInfo definition
type ServerInfo struct {
	Addr string "mapstructure:\"addr\" json:\"addr\""
}

// ObservabilityInfo definition
type ObservabilityInfo struct {
	Log LogInfo "mapstructure:\"log\" json:\"log\""
}

// LogInfo definition
type LogInfo struct {}

// MetricsInfo definition
type MetricsInfo struct {}

// PprofInfo definition
type PprofInfo struct {
	Enable bool "mapstructure:\"enable\" json:\"enable\""
	Addr string "mapstructure:\"addr\" json:\"addr\""
}

// PerformanceInfo definition
type PerformanceInfo struct {
	Enable bool "mapstructure:\"enable\" json:\"enable\""
	Type string "mapstructure:\"type\" json:\"type\""
}

// newConfig returns a new config pointer
func newConfig() *Config {
	return &Config{}
}

// GetConfig returns the global config pointer
func GetConfig() *Config {
	if globalConfig == nil {
		globalConfig = newConfig()
	}
	return globalConfig
}

// String returns config formatted in JSON
func (config *Config) String() string {
	jsonBuf, _ := json.Marshal(config)
	return cast.ToString(jsonBuf)
}

// Init initializes the config from config file
func (config *Config) Init (path string) error {
	// set viper setting
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	viper.AddConfigPath("../configs/")

	// read in config
	if err := viper.ReadInConfig(); err != nil {
		return xerrors.Errorf("viper.ReadInConfig failed! error:%w", err)
	}

	// unmarshal config
	if err := viper.Unmarshal(config); err != nil {
		return xerrors.Errorf("viper.Unmarshal failed! error:%w", err)
	}
	return nil
}
`

// ConfigYAMLInfo represents Config YAML information
type ConfigYAMLInfo struct {
	Opts *option.Options
	Addr string
}

// NewConfigYAMLInfo returns a new ConfigYAMLInfo pointer
func NewConfigYAMLInfo(opts ...option.Option) *ConfigYAMLInfo {
	// init options
	newOpts := option.NewOptions(opts...)
	return &ConfigYAMLInfo{Addr: newOpts.Addr, Opts: newOpts}
}

// RenderFile render template and output to file
func (yaml *ConfigYAMLInfo) RenderFile(outputPath string) error {
	if err := templates.RenderToFile(outputPath, yaml.Opts.Override, yaml.Opts.IgnoreExist,
		"template-config-yaml", ConfigYAMLTemplate, yaml); err != nil {
		return xerrors.Errorf("RenderToFile failed! error:%w", err)
	}
	return nil
}

var ConfigYAMLTemplate string = `# Server configuration
#
# server.addr
#	Server's address to listen on.
# 	eg: 
#		localhost:58888
#		127.0.0.1:58888
#		eth0:58888
server:
  addr: "{{.Addr}}"

# Observability configuration
observability:
  log:
  metrics:

# Client configuration
# [client]

# pprof configuration
#
# pprof.active
#	Is pprof enabled or not
# pprof.addr
#	Http address to listen on for getting profile information
#	eg:
#		wget http://localhost:50000/debug/pprof/profile
pprof:
  enable: false
  addr: "localhost:50000"

# performance configuration
# 	Performance will collect memory, gc, goroutine information periodically
#	and output to log or statsd
#
# performance.enable
#	Is performance monitor enabled or not
# performance.type
#   Type to output the performance monitor information, available option [log, statsd]
performance:
  enable: false
  type: "log"
`
