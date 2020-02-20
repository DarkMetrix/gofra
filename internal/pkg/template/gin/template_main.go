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

	config "{{.WorkingPathRelative}}/internal/pkg/config"
	application "{{.WorkingPathRelative}}/internal/app"
)

func main() {
	// start
	fmt.Println("====== Server [{{.Project}}] Start ======")

	// init config
	conf := config.GetConfig()

	err := conf.Init("../configs/config.toml")

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
