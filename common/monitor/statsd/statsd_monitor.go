package statsd

import (
	"fmt"
	"errors"

	log "github.com/cihub/seelog"

	"github.com/alexcesaro/statsd"
)

//Statsd client
var client *statsd.Client = nil
var addr string = "localhost:8125"
var project string = "Default"

//Init statsd client, if addr is empty, using default setting
func InitStatsd(addr string, project string) error {
	//If addr is empty, use default addr setting which is ":8125" in udp
	var err error
	if len(addr) == 0 {
		client, err = statsd.New()
	} else {
		client, err = statsd.New(
			statsd.Address(addr),
			statsd.Tags("project", project),
			statsd.TagsFormat(statsd.InfluxDB),
			statsd.ErrorHandler(func(err error) {
				log.Warnf("Statsd error:%v", err.Error())
			}))
	}

	if err != nil {
		log.Tracef(fmt.Sprintf("init statsd failed! error:%v", err.Error()))

		client = nil
		return err
	}

	log.Tracef(fmt.Sprintf("init statsd success! addr:%v, project:%v", addr, project))

	return nil
}

//Get statsd client
func GetStatsd() *statsd.Client {
	if client == nil {
		err := InitStatsd(addr, project)

		if err != nil {
			return nil
		}
	}

	return client
}

//Init
func Init(args... string) error {
	if len(args) < 2 {
		return errors.New(fmt.Sprintf("param invalid! args:%v", args))
	}

	addr = args[0]
	project = args[1]

	err := InitStatsd(addr, project)

	if err != nil {
		return err
	}

	return nil
}

//Increment
func Increment(bucket string) {
	statsd := GetStatsd()

	if statsd == nil {
		log.Tracef("monitor increment failed! bucket:%v", bucket)
		return
	}

	log.Tracef("monitor increment failed! bucket:%v", bucket)

	statsd.Increment(bucket)
}
