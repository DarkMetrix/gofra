package monitor

import (
	"github.com/alexcesaro/statsd"
)

//Statsd client
var client *statsd.Client = nil

//Init statsd client, if addr is empty, using default setting
func InitStatsd(addr string) {
	//If addr is empty, use default addr setting which is ":8125" in udp
	var err error
	if len(addr) == 0 {
		client, err = statsd.New()
	} else {
		client, err = statsd.New(statsd.Address(addr))
	}

	if err != nil {
		panic(err)
	}
}

//Get statsd client
func GetStatsd() *statsd.Client{
	if client == nil {
		InitStatsd("")
	}

	return client
}
