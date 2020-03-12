package seelog

import (
	"fmt"
	"errors"

	log "github.com/cihub/seelog"
)

// default log setting
var defaultSetting string = `
	<seelog>
        <outputs formatid="main">
			<console/>
        </outputs>
        <formats>
            <format id="main" format="[%LEVEL][%DateT%Time][%File:%Line][%FuncShort] => %Msg%n"/>
        </formats>
    </seelog>
	`

var path string = ""
var project string = ""

// init seelog using config file path, if empty using default setting
func Init(args ... string) error {
	// if path is empty, use default setting
	var logger log.LoggerInterface
	var err error

	if len(args) < 2 {
		return errors.New(fmt.Sprintf("param invalid! args:%v", args))
	}

	path = args[0]
	project = args[1]

	err = log.RegisterCustomFormatter("Project", createProjectFormatter)

	if err != nil {
		log.Tracef(err.Error())
		return err
	}

	if len(path) == 0 {
		logger, err = log.LoggerFromConfigAsString(defaultSetting)
	} else {
		logger, err = log.LoggerFromConfigAsFile(path)
	}

	if err != nil {
		return err
	}

	err = log.ReplaceLogger(logger)

	if err != nil {
		return err
	}

	return nil
}

func createProjectFormatter(params string) log.FormatterFunc {
	return func(message string, level log.LogLevel, context log.LogContextInterface) interface{} {
		return project
	}
}
