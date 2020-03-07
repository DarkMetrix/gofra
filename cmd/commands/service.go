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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	grpcTemplate "github.com/DarkMetrix/gofra/internal/pkg/template/grpc"
	commonUtils "github.com/DarkMetrix/gofra/pkg/utils"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Service operations [add, update]",
	Long: `Gofra is a framework using gRPC/gin as the communication layer.
service command will help to manipulate .proto file to generate service frame & handler.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// addServiceCmd represents the service add command
var addServiceCmd = &cobra.Command{
	Use:   "add",
	Short: "Add service (*.proto) to project",
	Long: `Gofra is a framework using gRPC/gin as the communication layer.
service add command will help to manipulate .proto file to add service frame & handler.`,
	Run: func(cmd *cobra.Command, args []string) {
		addService(servicePath, false, false)
	},
}

// updateServiceCmd represents the service update command
var updateServiceCmd = &cobra.Command{
	Use:   "update",
	Short: "Update service (*.proto) to project",
	Long: `Gofra is a framework using gRPC/gin as the communication layer.
service update command will help to manipulate .proto file to update service frame & handler.`,
	Run: func(cmd *cobra.Command, args []string) {
		addService(servicePath, override, true)
	},
}

var servicePath string
var update bool

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(addServiceCmd)
	serviceCmd.AddCommand(updateServiceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")
	addServiceCmd.PersistentFlags().StringVar(&servicePath, "path", "", "A .proto file to generate codes")
	addServiceCmd.PersistentFlags().StringVar(&protocPath, "protoc_path", "protoc", "protoc binary path, in case user has multi versions of protoc")
	addServiceCmd.PersistentFlags().StringArrayVar(&protoFileIncludePath, "proto_path", []string{}, "proto files include path used by protoc's command '--proto_path'")

	updateServiceCmd.PersistentFlags().StringVar(&servicePath, "path", "", "A .proto file to generate codes")
	updateServiceCmd.PersistentFlags().StringVar(&protocPath, "protoc_path", "protoc", "protoc binary path, in case user has multi versions of protoc")
	updateServiceCmd.PersistentFlags().StringArrayVar(&protoFileIncludePath, "proto_path", []string{}, "proto files path include used by protoc's command '--proto_path'")
	updateServiceCmd.PersistentFlags().BoolVar(&override, "override", false,"If override when file exists")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func addService(path string, override, update bool) error {
	fmt.Println("====== Gofra service add ======")

	//Set protoc binary path
	if len(protocPath) != 0 {
		grpcTemplate.ProtocPath = protocPath
	}

	//Check path
	fmt.Printf("\r\nChecking Path ......")
	workingPath, err := os.Getwd()

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\nWorking path:%v\r\n", workingPath)
	}

	//Read template
	fmt.Printf("\r\nReading template ......")
	if len(templatePath) == 0 {
		fmt.Printf(" failed! \r\nerror:Template file path is empty!\r\n")
		return err
	}

	templateInfo, _, err := ReadTemplate(templatePath)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	//Check server type
	if templateInfo.Type != "grpc" {
		fmt.Printf(" failed! \r\nerror:Server type is not 'grpc'!\r\n")
		return err
	}

	//Mkdir
	fmt.Printf("\r\nMake dir ......")
	filename := filepath.Base(path)
	protoPath := filepath.Join(workingPath, "api", "protobuf_spec", strings.TrimSuffix(filename, ".proto"))

	commonUtils.CreatePath(protoPath, override)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	//Copy proto file path & generate .pb.go file
	fmt.Printf("\r\nCopy proto file ......")
	protoFilePath := filepath.Join(workingPath, "api", "protobuf_spec", strings.TrimSuffix(filename, ".proto"), filename)
	protoFilePathRelative := filepath.Join(".", "api", "protobuf_spec", strings.TrimSuffix(filename, ".proto"), filename)

	err = commonUtils.CopyFile(path, protoFilePath)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	//Execute protoc to generate .pb.go file
	fmt.Printf("\r\nCompile proto file ......")

	args := []string{}

	goPath, err := commonUtils.GetGOPATH()

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	}

	protoFileIncludePath = append(protoFileIncludePath, workingPath, goPath + "/src")

	for _, path := range protoFileIncludePath {
		arg := fmt.Sprintf("--proto_path=%v", path)
		args = append(args, arg)
	}

	args = append(args,"--go_out=plugins=grpc:.")
	args = append(args, protoFilePathRelative)

	shellCmd := exec.Command(grpcTemplate.ProtocPath, args...)

	output, err := shellCmd.CombinedOutput()

	if err != nil {
		fmt.Printf(" failed! cmd:%v %v \r\nerror:%v\r\noutput:%v\r\n", grpcTemplate.ProtocPath, args, err.Error(), string(output))
		return err
	} else {
		fmt.Printf(" success! cmd:%v %v\r\n", grpcTemplate.ProtocPath, args)
	}

	//Generate service
	fmt.Printf("\r\nGenerating service code ......")
	err = grpcTemplate.GenerateService(workingPath, protoFileIncludePath, templateInfo, protoFilePathRelative, override, update)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	return nil
}
