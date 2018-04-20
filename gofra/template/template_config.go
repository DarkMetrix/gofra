package template

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

//Global config
var globalConfig *Config

//Tracing info
type TracingInfo struct {
	Params []string "mapstructure:\"params\" json:\"params\""
}

//Monitor info
type MonitorInfo struct {
	Params []string "mapstructure:\"params\" json:\"params\""
}

//Server config
type ServerInfo struct {
	Addr string "mapstructure:\"addr\" json:\"addr\""
}

//Client config
type ClientInfo struct {
	Pool PoolInfo "mapstructure:\"pool\" json:\"pool\""
}

//Pool config
type PoolInfo struct {
	InitConns int "mapstructure:\"init_conns\" json:\"init_conns\""
	MaxConns int "mapstructure:\"max_conns\" json:\"max_conns\""
	IdleTime int "mapstructure:\"idle_time\" json:\"idle_time\""
}

//Config sturcture
type Config struct {
	Monitor MonitorInfo "mapstructure:\"monitor\" json:\"monitor\""
	Tracing TracingInfo "mapstructure:\"tracing\" json:\"tracing\""
	Server ServerInfo "mapstructure:\"server\" json:\"server\""
	Client ClientInfo "mapstructure:\"client\" json:\"client\""
}

//New Config
func NewConfig() *Config {
	return &Config{}
}

//Get singleton config
func GetConfig() *Config {
	if globalConfig == nil {
		globalConfig = NewConfig()
	}

	return globalConfig
}

//Init config from json file
func (config *Config) Init (path string) error {
	//Set viper setting
	viper.SetConfigType("toml")
	viper.SetConfigFile(path)
	viper.AddConfigPath("../conf/")

	//Read in config
	err := viper.ReadInConfig()

	if err != nil {
		return err
	}

	//Unmarshal config
	err = viper.Unmarshal(config)

	if err != nil {
		return err
	}

	return nil
}
`

type ConfigTomlInfo struct {
	Addr string
	InitConns int
	MaxConns int
	IdleTime int

	MonitorInitParams string
	TracingInitParams string
}

var ConfigTomlTemplate string = `
# Server configuration
#
# server.addr
#	Server's address to listen on.
# 	eg: 
#		localhost:58888
#		127.0.0.1:58888
#		eth0:58888
[server]
    addr="{{.Addr}}"

# Client configuration
#
# client.pool.init_conns
#	Client's initialized connections of the pool
#
# client.pool.max_conns
#	Client's max connections of the pool
#
# client.pool.idle_time
#	Client's idle time of connection to be recycled
[client.pool]
    init_conns={{.InitConns}}
    max_conns={{.MaxConns}}
    idle_time={{.IdleTime}}

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
#	Gofra take the zipkin as the default tracing system
#	so the params has 4 parts(all in a string array)
#		1.zipkin's url
#		2.debug flag 
#		3.server's address
#		4.the project's name
#	eg:
#		params=["http://127.0.0.1:9411/api/v1/spans", "false", "localhost:58888", "demo"]
[tracing]
    params=[{{.TracingInitParams}}]
`

