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
	viper.SetConfigType("json")
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

type ConfigJsonInfo struct {
	Addr string
	InitConns int
	MaxConns int
	IdleTime int

	MonitorInitParams string
	TracingInitParams string
}

var ConfigJsonTemplate string = `
{
  "monitor":
  {
    "params"":[{{.MonitorInitParams}}]
  },
  "tracing":
  {
    "params"":[{{.TracingInitParams}}]
  },
  "server":
  {
    "addr":"{{.Addr}}"
  },
  "client":
  {
    "pool":
    {
      "init_conns":{{.InitConns}},
      "max_conns":{{.MaxConns}},
      "idle_time":{{.IdleTime}}
    }
  }
}
`

