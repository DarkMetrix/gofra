package naming

import (
	"fmt"
	"errors"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type NamingResovler interface {
	GetAddr(addrAlias string) (string, error)
}

type NamingConfig struct {
	Locations map[string]string `mapstructure:locations`    //Locations, eg:"user_service":"local|127.0.0.1:8088"
}

func getLocation(location string) (string, string, error) {
	parts := strings.Split(location, "|")

	if len(parts) != 2 {
		return "", "", errors.New(fmt.Sprintf("Addr is not valid! addr:%v", location))
	}

	name := parts[0]
	addr := parts[1]

	return name, addr, nil
}

type Naming struct {
	Resolvers map[string]NamingResovler     //Naming resovlers
	Config NamingConfig                     //Naming config
}

var naming *Naming = nil
var rwMutex sync.RWMutex

//Init naming
func Init(args... string) {
	if len(args) < 1 {
		panic("Init args length < 1")
	}

	if naming != nil {
		return
	}

	namingConfigPath := args[0]

	naming = &Naming{
		Resolvers: make(map[string]NamingResovler),
		Config: NamingConfig{
			Locations: make(map[string]string),
		},
	}

	//Set viper setting
	viper.SetConfigType("json")
	viper.SetConfigFile(namingConfigPath)
	viper.AddConfigPath("../conf/")

	//Read in config
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	//Unmarshal config
	err = viper.Unmarshal(&naming.Config)

	if err != nil {
		panic(err)
	}
}

func AddResolver(name string, resovler NamingResovler) {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	naming.Resolvers[name] = resovler
}

func GetAddr(service string) (string, error) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	//Get service location
	location, ok := naming.Config.Locations[service]

	if !ok {
		return "", errors.New(fmt.Sprintf("GetAddr - Service not found! service:%v", service))
	}

	//Get service addr
	name, addrAlias, err := getLocation(location)

	if err != nil {
		return "", err
	}

	resolver, ok := naming.Resolvers[name]

	if !ok {
		return "", errors.New(fmt.Sprintf("GetAddr - Resolver not found! resolver:%v", name))
	}

	addr, err := resolver.GetAddr(addrAlias)

	if err != nil {
		return "", err
	}

	return addr, nil
}
