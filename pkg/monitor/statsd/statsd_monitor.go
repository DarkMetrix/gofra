package statsd

import (
	"fmt"
	"errors"

	log "github.com/cihub/seelog"

	"github.com/alexcesaro/statsd"
)

// statsd client
var client *statsd.Client = nil
var addr string = "localhost:8125"
var project string = "Default"

// init statsd client, if addr is empty, using default setting
func InitStatsd(addr string, project string) error {
	// if addr is empty, use default addr setting which is ":8125" in udp
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

// get statsd client
func GetStatsd() *statsd.Client {
	if client == nil {
		err := InitStatsd(addr, project)

		if err != nil {
			return nil
		}
	}

	return client
}

// init
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

// increment
func Increment(bucket string) {
	statsd := GetStatsd()

	if statsd == nil {
		log.Tracef("monitor increment failed! bucket:%v", bucket)
		return
	}

	log.Tracef("monitor increment success! bucket:%v", bucket)

	statsd.Increment(bucket)
}

// count
func Count(bucket string, number interface{}) {
	statsd := GetStatsd()

	if statsd == nil {
		log.Tracef("monitor count failed! bucket:%v, count:%v", bucket, number)
		return
	}

	log.Tracef("monitor count success! bucket:%v, count:%v", bucket, number)

	statsd.Count(bucket, number)
}

// gauge
func Gauge(bucket string, value interface{}) {
	statsd := GetStatsd()

	if statsd == nil {
		log.Tracef("monitor gauge failed! bucket:%v, value:%v", bucket, value)
		return
	}

	log.Tracef("monitor gauge success! bucket:%v, value:%v", bucket, value)

	statsd.Gauge(bucket, value)
}

type MonitorTiming struct {
	statsd.Timing
}

// timing
func NewTiming() *MonitorTiming {
	statsd := GetStatsd()

	if statsd == nil {
		log.Tracef("monitor NewTiming failed!")
		return nil
	}

	log.Tracef("monitor NewTiming success!")

	timing := statsd.NewTiming()
	return &MonitorTiming{timing}
}

// sent timing
func (timing *MonitorTiming) SendTiming(bucket string) {
	log.Tracef("monitor SendTiming success! bucket:%v", bucket)

	timing.Send(bucket)
}
