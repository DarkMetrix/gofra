package src

import (
	"fmt"
	"os"

	config "github.com/DarkMetrix/gofra/test/server/src/config"
	application "github.com/DarkMetrix/gofra/test/server/src/application"
)

func main() {
	// start
	fmt.Println("====== Test grpc server ======")

	// init config
	conf := config.GetConfig()

	err := conf.Init("../conf/config.json")

	if err != nil {
		fmt.Println("Init config failed!")
		os.Exit(-1)
	}

	// init application
	var application application.Application

	err = application.Init(conf)

	if err != nil {
		fmt.Printf("Application init failed! error:%v\r\n", err.Error())
		os.Exit(-2)
	}

	err = application.Run("tcp", ":58888")

	if err != nil {
		fmt.Printf("Application init failed! error:%v\r\n", err.Error())
	}
}
