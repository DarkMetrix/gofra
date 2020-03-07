package gin

type MainInfo struct {
	Author string
	Time string
	Project string
	WorkingPathRelative string
	Addr string
}

var MainTemplate string = `
/**********************************
 * Author : {{.Author}}
 * Time : {{.Time}}
 **********************************/

package main

import (
	"fmt"
	"os"

	viper "github.com/spf13/viper"
	pflag "github.com/spf13/pflag"

	config "{{.WorkingPathRelative}}/internal/pkg/config"
	application "{{.WorkingPathRelative}}/internal/app"
)

func main() {
	// start
	fmt.Println("====== Server [{{.Project}}] Start ======")

	// init flags
	pflag.String("config.path", "../configs/config.toml", "Config file path, default '../configs/config.toml'")
	pflag.String("log.config.path", "../configs/log.config", "Log config file path, default '../configs/log.config'")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	// init config
	conf := config.GetConfig()

	err := conf.Init(viper.GetString("config.path"))

	if err != nil {
		fmt.Printf("Init config failed! error:%v\r\n", err.Error())
		os.Exit(-1)
	}

	// init application
	app := application.GetApplication()

	if app == nil {
		fmt.Printf("Application get failed!\r\n")
		os.Exit(-2)
	}

	err = app.Init(conf)

	if err != nil {
		fmt.Printf("Application init failed! error:%v\r\n", err.Error())
		os.Exit(-3)
	}

	fmt.Printf("Listen on port [%v]\r\n", conf.Server.HTTPAddr)

	// run application
	err = app.Run(conf.Server.HTTPAddr)

	if err != nil {
		fmt.Printf("Application run failed! error:%v\r\n", err.Error())
		os.Exit(-4)
	}
}
`
