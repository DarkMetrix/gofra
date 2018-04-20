package template

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

	config "{{.WorkingPathRelative}}/src/config"
	application "{{.WorkingPathRelative}}/src/application"
)

func main() {
	// start
	fmt.Println("====== Server [{{.Project}}] Start ======")

	// init config
	conf := config.GetConfig()

	err := conf.Init("../conf/config.toml")

	if err != nil {
		fmt.Printf("Init config failed! error:%v\r\n", err.Error())
		os.Exit(-1)
	}

	// init application
	var application application.Application

	err = application.Init(conf)

	if err != nil {
		fmt.Printf("Application init failed! error:%v\r\n", err.Error())
		os.Exit(-2)
	}

	fmt.Printf("Listen on port [%v]\r\n", conf.Server.Addr)

	err = application.Run(conf.Server.Addr)

	if err != nil {
		fmt.Printf("Application run failed! error:%v\r\n", err.Error())
		os.Exit(-3)
	}
}
`
