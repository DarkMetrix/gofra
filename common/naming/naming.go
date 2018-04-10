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
	Locations map[string]LocationConfig `mapstructure:locations`    //Locations, eg:"user_service":"local|127.0.0.1:8088"
}

type LocationConfig struct {
	IsTest bool `mapstructure:is_test`                   //Flag to indicate which location to use
	LocationReal string `mapstructure:location_real`     //Real location in production
	LocationTest string `mapstructure:location_test`     //Test location
}

func (config *LocationConfig) GetLocation() (string, string, error) {
	realLocation := ""

	if config.IsTest {
		realLocation = config.LocationTest
	} else {
		realLocation = config.LocationReal
	}

	parts := strings.Split(realLocation, "|")

	if len(parts) != 2 {
		return "", "", errors.New(fmt.Sprintf("Addr is not valid! addr:%v", realLocation))
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
			Locations: make(map[string]LocationConfig),
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
	err = viper.Unmarshal(config)

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
	name, addrAlias, err := location.GetLocation()

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
