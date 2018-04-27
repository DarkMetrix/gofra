package local

import (
	"fmt"
	"time"
	"errors"
	"strings"

	log "github.com/cihub/seelog"

	commonUtils "github.com/DarkMetrix/gofra/common/utils"
)

type LocalNamingResolver struct {
}

func (resolver *LocalNamingResolver) Lookup(location string) (string, error) {
	addrs := strings.Split(location, ",")

	if len(addrs) < 1 {
		log.Tracef("lookup failed! error:addrs is empty, location:%v", location)
		return "", errors.New(fmt.Sprintf("addrs is empty, location:%v", location))
	}

	index := time.Now().UnixNano() % int64(len(addrs))

	addr := commonUtils.GetRealAddrByNetwork(addrs[index])

	log.Tracef("lookup success! location:%v, addr:%v", location, addr)

	return addr, nil
}
