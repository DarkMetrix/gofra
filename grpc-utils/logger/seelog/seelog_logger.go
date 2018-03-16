package seelog

import (
	log "github.com/cihub/seelog"
)

//Default log setting
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

//Init seelog using config file path, if empty using default setting
func Init(args ... string) {
	//If path is empty, use default setting
	var logger log.LoggerInterface
	var err error

	if len(args) < 1 {
		panic("Init args length < 1")
	}

	path := args[0]

	if len(path) == 0 {
		logger, err = log.LoggerFromConfigAsString(defaultSetting)
	} else {
		logger, err = log.LoggerFromConfigAsFile(path)
	}

	if err != nil {
		panic(err)
	}

	err = log.ReplaceLogger(logger)

	if err != nil {
		panic(err)
	}
}
