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

package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	gofraTemplate "github.com/DarkMetrix/gofra/gofra/template"
	gofraUtils "github.com/DarkMetrix/gofra/gofra/utils"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Service operations [add, update]",
	Long: `Gofra is a framework using gRPC as the communication layer.\r\nservice command will help to manipulate .proto file to generate service frame & handler.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// addServiceCmd represents the service command
var addServiceCmd = &cobra.Command{
	Use:   "add",
	Short: "Add service (*.proto) to project",
	Long: `Gofra is a framework using gRPC as the communication layer.\r\nservice add command will help to manipulate .proto file to add service frame & handler.`,
	Run: func(cmd *cobra.Command, args []string) {
		addService(servicePath, false, false)
	},
}

// updateServiceCmd represents the service command
var updateServiceCmd = &cobra.Command{
	Use:   "update",
	Short: "Update service (*.proto) to project",
	Long: `Gofra is a framework using gRPC as the communication layer.\r\nservice add command will help to manipulate .proto file to update service frame & handler.`,
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

	updateServiceCmd.PersistentFlags().StringVar(&servicePath, "path", "", "A .proto file to generate codes")
	updateServiceCmd.PersistentFlags().BoolVar(&override, "override", false,"If override when file exists")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func addService(path string, override, update bool) error {
	fmt.Println("====== Gofra service add ======")

	//Check path
	fmt.Printf("\r\nChecking Path ......")
	goPath, workingPath, err := CheckPath()

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\nGOPATH:%v\r\nWorking path:%v\r\n", goPath, workingPath)
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

	//Mkdir
	fmt.Printf("\r\nMake dir ......")
	filename := filepath.Base(path)
	protoPath := filepath.Join(workingPath, "src", "proto", strings.TrimSuffix(filename, ".proto"))

	gofraUtils.CreatePath(protoPath, override)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	//Copy proto file path & generate .pb.go file
	fmt.Printf("\r\nCopy proto file ......")
	protoFilePath := filepath.Join(workingPath, "src", "proto", strings.TrimSuffix(filename, ".proto"), filename)
	protoFilePathRelative := filepath.Join(".", "src", "proto", strings.TrimSuffix(filename, ".proto"), filename)

	err = gofraUtils.CopyFile(path, protoFilePath)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	//Execute protoc to generate .pb.go file
	fmt.Printf("\r\nCompile proto file ......")
	shellCmd := exec.Command("protoc", "--go_out=plugins=grpc:.", protoFilePathRelative)

	err = shellCmd.Run()

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	//Generate service
	fmt.Printf("\r\nGenerating service code ......")
	err = gofraTemplate.GenerateService(workingPath, goPath, templateInfo, path, override, update)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	return nil
}