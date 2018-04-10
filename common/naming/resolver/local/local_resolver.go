package local

import (
	"fmt"
	"time"
	"errors"
	"strings"

	commonUtils "github.com/DarkMetrix/gofra/common/utils"
)

type LocalNamingResolver struct {
}

func (resolver *LocalNamingResolver) Lookup(addrAlias string) (string, error) {
	addrs := strings.Split(addrAlias, ",")

	if len(addrs) < 1 {
		return "", errors.New(fmt.Sprintf("addrs is empty! addrAlias:%v", addrAlias))
	}

	index := time.Now().UnixNano() % int64(len(addrs))

	return commonUtils.GetRealAddrByNetwork(addrs[index]), nil
}
