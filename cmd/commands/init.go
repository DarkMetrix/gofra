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
	"fmt"
	"os"
	"os/exec"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

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

		// check path
		fmt.Printf("\r\nChecking Path ......")
		workingPath, err := os.Getwd()

		if err != nil {
			fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
			return
		} else {
			fmt.Printf(" success! \r\nWorking path:%v\r\n", workingPath)
		}

		// read template
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

		// init directory structure
		fmt.Printf("\r\nInitializing directory structure ......")

		switch templateInfo.Type {
		case "grpc":
			// init directory structure for grpc
			err = InitGrpcDirectoryStructure(workingPath, templateInfo)

			if err != nil {
				fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
				return
			} else {
				fmt.Printf(" success!\r\n")
			}

			// init all files
			fmt.Printf("\r\nInitializing all files ......")
			err = InitGrpcAllFiles(workingPath, templateInfo)

			if err != nil {
				fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
				return
			} else {
				fmt.Printf(" success!\r\n")
			}

		case "http":
			// init directory structure for http
			err = InitHttpDirectoryStructure(workingPath, templateInfo)

			if err != nil {
				fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
				return
			} else {
				fmt.Printf(" success!\r\n")
			}

			// init all files
			fmt.Printf("\r\nInitializing all files ......")
			err = InitHttpAllFiles(workingPath, templateInfo)

			if err != nil {
				fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
				return
			} else {
				fmt.Printf(" success!\r\n")
			}

		default:
			fmt.Printf(" failed! \r\nerror:Invalid server type\r\n")
		}

		// init go module
		gomodPath := filepath.Join(workingPath, "go.mod")

		isExist, err := commonUtils.CheckPathExists(gomodPath)

		if err != nil {
			fmt.Printf("Check go.mod file existance failed! path:%v, error:%v", gomodPath, err.Error())
			return
		}

		if !isExist {
			shellCmd := exec.Command("go", "mod", "init", filepath.Base(workingPath))

			err = shellCmd.Run()

			if err != nil {
				fmt.Printf("go mod init %v failed! error:%v", filepath.Base(workingPath), err.Error())
			} else {
				fmt.Printf("go mod init %v success!", filepath.Base(workingPath))
			}
		} else {
			fmt.Printf("go.mod file already exist!")
		}

		// print application directory structure
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

// server config
type ServerInfo struct {
	Addr string `json:"addr"`
}

// client config
type ClientInfo struct {
	Pool PoolInfo `json:"pool"`
}

// pool config
type PoolInfo struct {
	InitConns int `json:"init_conns"`
	MaxConns int `json:"max_conns"`
	IdleTime int `json:"idle_time"`
}

// monitor package
type MonitorPackageInfo struct {
	Package string `json:"package"`
	InitParam string `json:"init_param"`
}

// interceptor package
type InterceptorPackageInfo struct {
	Package string `json:"package"`
}

// template info
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

// read template json file to ge information about how to generate the application
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

// init application directory structure for grpc
func InitGrpcDirectoryStructure(workingPath string, info *gofraTemplate.TemplateInfo) error {
	buildPath := filepath.Join(workingPath, "build")
	confPath := filepath.Join(workingPath, "configs")
	logPath := filepath.Join(workingPath, "log")
	cmdPath := filepath.Join(workingPath, "cmd")
	apiPath := filepath.Join(workingPath, "api")
	internalPath := filepath.Join(workingPath, "internal")

	internalAppPath := filepath.Join(workingPath, "internal", "app")
	internalPkgPath := filepath.Join(workingPath, "internal", "pkg")
	internalCommonPath := filepath.Join(workingPath, "internal", "pkg", "common")
	internalConfigPath := filepath.Join(workingPath, "internal", "pkg", "config")
	internalHandlerPath := filepath.Join(workingPath, "internal", "grpc_handler")

	apiProtobufPath := filepath.Join(workingPath, "api", "protobuf_spec")
	apiProtobufHealthCheckPath := filepath.Join(workingPath, "api", "protobuf_spec", "health_check")

	// create root directories
	err := commonUtils.CreatePaths(override, buildPath, confPath, logPath, cmdPath, apiPath, internalPath)

	if err != nil {
		return err
	}

	// create internal sub directories
	err = commonUtils.CreatePaths(override, internalAppPath, internalPkgPath, internalCommonPath, internalConfigPath, internalHandlerPath)

	if err != nil {
		return err
	}

	// create api sub directories
	err = commonUtils.CreatePaths(override, apiProtobufPath, apiProtobufHealthCheckPath)

	if err != nil {
		return err
	}

	return nil
}

// init application directory structure for http
func InitHttpDirectoryStructure(workingPath string, info *gofraTemplate.TemplateInfo) error {
	buildPath := filepath.Join(workingPath, "build")
	confPath := filepath.Join(workingPath, "configs")
	logPath := filepath.Join(workingPath, "log")
	cmdPath := filepath.Join(workingPath, "cmd")
	apiPath := filepath.Join(workingPath, "api")
	internalPath := filepath.Join(workingPath, "internal")

	internalAppPath := filepath.Join(workingPath, "internal", "app")
	internalPkgPath := filepath.Join(workingPath, "internal", "pkg")
	internalCommonPath := filepath.Join(workingPath, "internal", "pkg", "common")
	internalConfigPath := filepath.Join(workingPath, "internal", "pkg", "config")
	internalHandlerPath := filepath.Join(workingPath, "internal", "http_handler")

	apiHttpPath:= filepath.Join(workingPath, "api", "http_spec")

	// create root directories
	err := commonUtils.CreatePaths(override, buildPath, confPath, logPath, cmdPath, apiPath, internalPath)

	if err != nil {
		return err
	}

	// create internal sub directories
	err = commonUtils.CreatePaths(override, internalAppPath, internalPkgPath, internalCommonPath, internalConfigPath, internalHandlerPath)

	if err != nil {
		return err
	}

	// create api sub directories
	err = commonUtils.CreatePaths(override, apiHttpPath)

	if err != nil {
		return err
	}

	return nil
}

// init all go file with template for grpc
func InitGrpcAllFiles(workingPath string, info *gofraTemplate.TemplateInfo) error {
	// set protoc binary path
	if len(protocPath) != 0 {
		grpcTemplate.ProtocPath = protocPath
	}

	err := grpcTemplate.GenerateCommonFile(workingPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateConfigFile(workingPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateConfigTomlFile(workingPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateConfigLogFile(workingPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateApplicationFile(workingPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateMainFile(workingPath, info, override)

	if err != nil {
		return err
	}

	err = grpcTemplate.GenerateHealthCheckHandler(workingPath, info, override)

	if err != nil {
		return err
	}

	return nil
}

// init all go file with template for http
func InitHttpAllFiles(workingPath string, info *gofraTemplate.TemplateInfo) error {
	err := httpTemplate.GenerateCommonFile(workingPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateConfigFile(workingPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateConfigTomlFile(workingPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateConfigLogFile(workingPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateApplicationFile(workingPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateMainFile(workingPath, info, override)

	if err != nil {
		return err
	}

	err = httpTemplate.GenerateHealthCheckHttpHandler(workingPath, info, override)

	if err != nil {
		return err
	}

	return nil
}
