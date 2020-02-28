package gin

type ConfigInfo struct {
	Author string
	Time string
}

var ConfigTemplate string = `
/**********************************
 * Author : {{.Author}}
 * Time : {{.Time}}
 **********************************/

package config

import (
	"github.com/spf13/viper"
)

// global config
var globalConfig *Config

// pprof
type PprofInfo struct {
	Addr string "mapstructure:\"addr\" json:\"addr\""
	Active uint32 "mapstructure:\"active\" json:\"active\""
}

// performance info
type PerformanceInfo struct {
	Active uint32 "mapstructure:\"active\" json:\"active\""
	Type string "mapstructure:\"type\" json:\"type\""
}

// tracing info
type TracingInfo struct {
	Params []string "mapstructure:\"params\" json:\"params\""
}

// monitor info
type MonitorInfo struct {
	Params []string "mapstructure:\"params\" json:\"params\""
}

// server config
type ServerInfo struct {
	HTTPAddr string "mapstructure:\"http_addr\" json:\"http_addr\""
	GinDebug int "mapstructure:\"gin_debug\" json:\"gin_debug\""
}

// client config
type ClientInfo struct {
}

// pool config
type PoolInfo struct {
	InitConns int "mapstructure:\"init_conns\" json:\"init_conns\""
	MaxConns int "mapstructure:\"max_conns\" json:\"max_conns\""
	IdleTime int "mapstructure:\"idle_time\" json:\"idle_time\""
}

// config structure
type Config struct {
	Pprof PprofInfo "mapstructure:\"pprof\" json:\"pprof\""
	Performance PerformanceInfo "mapstructure:\"performance\" json:\"performance\""
	Monitor MonitorInfo "mapstructure:\"monitor\" json:\"monitor\""
	Tracing TracingInfo "mapstructure:\"tracing\" json:\"tracing\""
	Server ServerInfo "mapstructure:\"server\" json:\"server\""
	Client ClientInfo "mapstructure:\"client\" json:\"client\""
}

// new Config
func newConfig() *Config {
	return &Config{}
}

// get singleton config
func GetConfig() *Config {
	if globalConfig == nil {
		globalConfig = newConfig()
	}

	return globalConfig
}

// init config from json file
func (config *Config) Init (path string) error {
	// set viper setting
	viper.SetConfigType("toml")
	viper.SetConfigFile(path)
	viper.AddConfigPath("../configs/")

	// read in config
	err := viper.ReadInConfig()

	if err != nil {
		return err
	}

	// unmarshal config
	err = viper.Unmarshal(config)

	if err != nil {
		return err
	}

	return nil
}
`

type ConfigTomlInfo struct {
	Addr string

	MonitorInitParams string
	TracingInitParams string
}

var ConfigTomlTemplate string = `
# Server configuration
#
# server.http_addr
#	Server's address to listen on.
# 	eg: 
#		localhost:58888
#		127.0.0.1:58888
#		eth0:58888
#
# server.gin_debug
#	gin's debug flag, 0=release mode, other=debug mode
[server]
    http_addr="{{.Addr}}"
	gin_debug=1

# Client configuration
# [client]

# pprof configuration
#
# pprof.active
#	Is pprof active or not, 0 = not active, otherwise active
#
# pprof.addr
#	Http address to listen on for getting profile information
#	eg:
#		wget http://localhost:50000/debug/pprof/profile
[pprof]
	active=0
	addr="localhost:50000"

# performance configuration
# 	Performance will collect memory, gc, goroutine information periodically
#	and output to log or statsd
#
# performance.active
#	Is performance monitor active or not, 0 = not active, otherwise active
#
# performance.type
#   Type to output the performance monitor information, available option [log, statsd]
[performance]
	active=0
	type="log"

# Monitor configuration
#
# monitor.params
#	Monitor's params to init
#	Gofra take the statsd as the default monitor system
#	so the params has 2 parts(all in a string array)
#		1.statsd's UDP address
#		2.the project's name
#	eg:
#		params=["127.0.0.1:8125", "demo"]
[monitor]
    params=[{{.MonitorInitParams}}]

# Tracing configuration
#
# tracing.params
#	tracing's params to init
#	Gofra take the jaeger as the default tracing system
#	so the params has 2 parts(all in a string array)
#		1.jaeger's agent udp address
#		2.the project's name
#	eg:
#		params=["127.0.0.1:6831", "demo"]
[tracing]
    params=[{{.TracingInitParams}}]
`

