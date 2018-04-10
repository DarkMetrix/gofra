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

func (resolver *LocalNamingResolver) Lookup(addrAlias string) (string, error) {
	addrs := strings.Split(addrAlias, ",")

	if len(addrs) < 1 {
		log.Tracef("lookup failed! error:addrs is empty, addrAlias:%v", addrAlias)
		return "", errors.New(fmt.Sprintf("addrs is empty, addrAlias:%v", addrAlias))
	}

	index := time.Now().UnixNano() % int64(len(addrs))

	addr := commonUtils.GetRealAddrByNetwork(addrs[index])

	log.Tracef("lookup success! addrAlias:%v, addr:%v", addrAlias, addr)

	return addr, nil
}
