// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	gofraTemplate "github.com/DarkMetrix/gofra/internal/pkg/template"
	httpTemplate "github.com/DarkMetrix/gofra/internal/pkg/template/gin"
	grpcTemplate "github.com/DarkMetrix/gofra/internal/pkg/template/grpc"
	commonUtils "github.com/DarkMetrix/gofra/pkg/utils"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gofra application using template.json",
	Long: `Gofra is a framework using gRPC/gin as the communication layer.\r\ninit command will create basic framework structure.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("====== Gofra init ======")

		//Check path
		fmt.Printf("\r\nChecking Path ......")
		goPath, workingPath, err := CheckPath()

		if err != nil {
			fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
			return
		} else {
			fmt.Printf(" success! \r\nGOPATH:%v\r\nWorking path:%v\r\n", goPath, workingPath)
		}

		//Read template
		fmt.Printf("\r\nReading template ......")
		if len(templatePath) == 0 {
			fmt.Printf(" failed! \r\nerror:Template file path is empty!\r\n")
			return
		}

		templateInfo, jsonBuffer, err := ReadTemplate(templatePath)

		if err != nil {
			fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
			return
		} else {
			fmt.Printf(" success! \r\nTemplate:\r\n%v\r\n", jsonBuffer)
		}

		//Init directory structure
		fmt.Printf("\r\nInitializing directory structure ......")

		switch templateInfo.Type {
		case "grpc":
			//Init directory structure for grpc
			err = InitGrpcDirectoryStructure(workingPath, templateInfo)

			if err != nil {
				fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
				return
			} else {
				fmt.Printf(" success!\r\n")
			}

			//Init all files
			fmt.Printf("\r\nInitializing all files ......")
			err = InitGrpcAllFiles(workingPath, goPath, templateInfo)

			if err != nil {
				fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
				return
			} else {
				fmt.Printf(" success!\r\n")
			}

		case "http":
			//Init directory structure for http
			err = InitHttpDirectoryStructure(workingPath, templateInfo)

			if err != nil {
				fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
				return
			} else {
				fmt.Printf(" success!\r\n")
			}

			//Init all files
			fmt.Printf("\r\nInitializing all files ......")
			err = InitHttpAllFiles(workingPath, goPath, templateInfo)

			if err != nil {
				fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
				return
			} else {
				fmt.Printf(" success!\r\n")
			}

		default:
			fmt.Printf(" failed! \r\nerror:Invalid server type\r\n")
		}

		//Print application directory structure
		fmt.Printf("\r\nApplication '%v' directory structure", templateInfo.Project)
		filepath.Walk(workingPath, func(path string, info os.FileInfo, err error) error {
			relPath, err := filepath.Rel(workingPath, path)

			if err != nil {
				return err
			}

			fmt.Println(relPath)
			return nil
		})
	},
}

var templatePath string
var protocPath string
var protoFileIncludePath []string
var override bool

//Server config
type ServerInfo struct {
	Addr string `json:"addr"`
}

//Client config
type ClientInfo struct {
	Pool PoolInfo `json:"pool"`
}

//Pool config
type PoolInfo struct {
	InitConns int `json:"init_conns"`
	MaxConns int `json:"max_conns"`
	IdleTime int `json:"idle_time"`
}

//Monitor package
type MonitorPackageInfo struct {
	Package string `json:"package"`
	InitParam string `json:"init_param"`
}

//Interceptor package
type InterceptorPackageInfo struct {
	Package string `json:"package"`
}

//Template info
type TemplateInfo struct {
	Author string `json:"author"`
	Project string `json:"project"`
	Version string `json:"version"`
	Server ServerInfo `json:"server"`
	Client ClientInfo `json:"client"`
	MonitorPackage MonitorPackageInfo `json:"monitor_package"`
	InterceptorPackage InterceptorPackageInfo `json:"interceptor_package"`
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	initCmd.PersistentFlags().StringVar(&templatePath, "path", "./template.json", "A template file in json to tell how to generate codes")
	initCmd.PersistentFlags().StringVar(&protocPath, "protoc_path", "protoc", "protoc binary path, in case user has multi versions of protoc")
	initCmd.PersistentFlags().BoolVar(&override, "override", false,"If override when file exists")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//Check if current working directory is under gopath
func CheckPath() (string, string, error) {
	goPath := os.Getenv("GOPATH")

	if len(goPath) == 0 {
		return "", "", errors.New("GOPATH is not set!")
	}

	workingPath, err := os.Getwd()

	if err != nil {
		return "", "", err
	}

	isMatch := strings.HasPrefix(workingPath, filepath.Join(goPath, "src"))

	if !isMatch {
		return "", "", errors.New("The working path is not under $GOPATH/src!")
	}

	return goPath, workingPath, nil
}

//Read template json file to ge information about how to generate the application
func ReadTemplate(templatePath string) (*gofraTemplate.TemplateInfo, string, error) {
	data, err := ioutil.ReadFile(templatePath)

	if err != nil {
		return nil, "", err
	}

	var info *gofraTemplate.TemplateInfo = new(gofraTemplate.TemplateInfo)

	err = json.Unmarshal(data, info)

	if err != nil {
		return nil, "", err
	}

	if len(info.Type) == 0 {
		info.Type = "grpc"
	}

	return info, string(data), nil
}

//Init application directory structure for grpc
func InitGrpcDirectoryStructure(workingPath string, info *gofraTemplate.TemplateInfo) error {
	buildPath := filepath.Join(workingPath, "build")
	confPath := filepath.Join(workingPath, "configs")
	logPath := filepath.Join(workingPath, "log")
	cmdPath := filepath.Join(workingPath, "cmd")
	apiPath := filepath.Join(workingPath, "api")
	internalPath := filepath.Join(workingPath, "internal")
	testPath := filepath.Join(workingPath, "test")

	internalAppPath := filepath.Join(workingPath, "internal", "app")
	internalPkgPath := filepath.Join(workingPath, "internal", "pkg")
	internalCommonPath := filepath.Join(workingPath, "internal", "pkg", "common")
	internalConfigPath := filepath.Join(workingPath, "internal", "pkg", "config")
	internalHandlerPath := filepath.Join(workingPath, "internal", "grpc_handler")

	apiProtobufPath := filepath.Join(workingPath, "api", "protobuf_spec")
	apiProtobufHealthCheckPath := filepath.Join(workingPath, "api", "protobuf_spec", "health_check")

	//Create root directories
	err := commonUtils.CreatePaths(override, buildPath, confPath, logPath, cmdPath, apiPath, internalPath, testPath)

	if err != nil {
		return err
	}

	//Create internal sub directories
	err = commonUtils.CreatePaths(override, internalAppPath, internalPkgPath, internalCommonPath, internalConfigPath, internalHandlerPath)

	if err != nil {
		return err
	}

	//Create api sub directories
	err = commonUtils.CreatePaths(override, apiProtobufPath, apiProtobufHealthCheckPath)

	if err != nil {
		return err
	}

	return nil
}

//Init application directory structure for http
func InitHttpDirectoryStructure(workingPath string, info *gofraTemplate.TemplateInfo) error {
	buildPath := filepath.Join(workingPath, "build")
	confPath := filepath.Join(workingPath, "configs")
	logPath := filepath.Join(workingPath, "log")
	cmdPath := filepath.Join(workingPath, "cmd")
	apiPath := filepath.Join(workingPath, "api")
	internalPath := filepath.Join(workingPath, "internal")
	testPath := filepath.Join(workingPath, "test")

	internalAppPath := filepath.Join(workingPath, "internal", "app")
	internalPkgPath := filepath.Join(workingPath, "internal", "pkg")
	internalCommonPath := filepath.Join(workingPath, "internal", "pkg", "common")
	internalConfigPath := filepath.Join(workingPath, "internal", "pkg", "config")
	internalHandlerPath := filepath.Join(workingPath, "internal", "http_handler")

	apiHttpPath:= filepath.Join(workingPath, "api", "http_spec")
	apiHttpHealthCheckPath := filepath.Join(workingPath, "api", "http_spec", "health_check")

	//Create root directories
	err := commonUtils.CreatePaths(override, buildPath, confPath, logPath, cmdPath, apiPath, internalPath, testPath)

	if err != nil {
		return err
	}

	//Create internal sub directories
	err = commonUtils.CreatePaths(override, internalAppPath, internalPkgPath, internalCommonPath, internalConfigPath, internalHandlerPath)

	if err != nil {
		return err
	}

	//Create api sub directories
	err = commonUtils.CreatePaths(override, apiHttpPath, apiHttpHealthCheckPath)

	if err != nil {
		return err
	}

	return nil
}

//Init all go file with template for grpc
func InitGrpcAllFiles(workingPath, goPath string, info *gofraTemplate.TemplateInfo) error {
	//Set protoc binary path
	if len(protocPath) != 0 {
		grpcTemplate.ProtocPath = protocPath
	}

	err := grpcTemplate.GenerateCommonFile(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateConfigFile(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateConfigTomlFile(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateConfigLogFile(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateApplicationFile(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateMainFile(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateHealthCheckHandler(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateTestClient(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	return nil
}

//Init all go file with template for http
func InitHttpAllFiles(workingPath, goPath string, info *gofraTemplate.TemplateInfo) error {
	err := httpTemplate.GenerateCommonFile(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateConfigFile(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateConfigTomlFile(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateConfigLogFile(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateApplicationFile(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateMainFile(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateHealthCheckHttpHandler(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateTestClient(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	return nil
}
