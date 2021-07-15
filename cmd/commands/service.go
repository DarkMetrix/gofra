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
	"path/filepath"

	"github.com/DarkMetrix/gofra/internal/pkg/directory"
	"github.com/DarkMetrix/gofra/internal/pkg/generate"
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/utils/gomod"
	"github.com/DarkMetrix/gofra/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Service operations [add, update]",
	Long: `Gofra is a framework using gRPC as the communication layer.
service command will help to manipulate .proto file to generate service frame & handler.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// addServiceCmd represents the service add command
var addServiceCmd = &cobra.Command{
	Use:   "add",
	Short: "Add service (*.proto) to project",
	Long: `Gofra is a framework using gRPC as the communication layer.
service add command will help to manipulate .proto file to add service frame & handler.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("====== gofra service add ======")

		// get go module from go.mod
		goModule, err := gomod.GetGoModule(filepath.Join(outputPath, "go.mod"))
		if err != nil {
			log.Fatalf("utils.GetGoModule failed! error:%v", err)
		}

		opts := []option.Option{
			option.WithOutputPath(outputPath),
			option.WithGoModule(goModule),
			option.WithOverride(false),
			option.WithProtocPath(protocPath),
			option.WithProtoFileIncludePath(protoFileIncludePath),
		}

		// copy proto
		layout := directory.NewGRPCLayout(opts...)
		if err := utils.CreatePaths(false, layout.GetAPIProtobufPath(protoFilePath)); err != nil {
			log.Fatalf("utils.CreatePaths failed! error:%+v", err)
		}
		if err := utils.CopyFile(protoFilePath, layout.GetAPIProtobufFilePath(protoFilePath)); err != nil {
			log.Fatalf("utils.CopyFile failed! error:%+v", err)
		}

		// add service
		if err := generate.NewGRPCServiceGenerator().Add(
			layout.GetAPIProtobufFilePath(protoFilePath), layout, opts...); err != nil {
			log.Fatalf("generate.Add failed! error:%+v", err)
		}
	},
}

// updateServiceCmd represents the service update command
var updateServiceCmd = &cobra.Command{
	Use:   "update",
	Short: "Update service (*.proto) to project",
	Long: `Gofra is a framework using gRPC as the communication layer.
service update command will help to manipulate .proto file to update service frame & handler.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("====== gofra service update ======")

		// get go module from go.mod
		goModule, err := gomod.GetGoModule(filepath.Join(outputPath, "go.mod"))
		if err != nil {
			log.Fatalf("utils.GetGoModule failed! error:%v", err)
		}

		opts := []option.Option{
			option.WithOutputPath(outputPath),
			option.WithGoModule(goModule),
			option.WithOverride(false),
			option.WithProtocPath(protocPath),
			option.WithProtoFileIncludePath(protoFileIncludePath),
		}

		// copy proto
		layout := directory.NewGRPCLayout(opts...)
		if err := utils.CreatePaths(false, layout.GetAPIProtobufPath(protoFilePath)); err != nil {
			log.Fatalf("utils.CreatePaths failed! error:%+v", err)
		}
		if err := utils.CopyFile(protoFilePath, layout.GetAPIProtobufFilePath(protoFilePath)); err != nil {
			log.Fatalf("utils.CopyFile failed! error:%+v", err)
		}

		// update service
		if err := generate.NewGRPCServiceGenerator().Update(
			layout.GetAPIProtobufFilePath(protoFilePath), layout, opts...); err != nil {
			log.Fatalf("generate.Update failed! error:%+v", err)
		}
	},
}

var (
	protoFilePath string
)

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(addServiceCmd)
	serviceCmd.AddCommand(updateServiceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")
	addServiceCmd.PersistentFlags().StringVar(&outputPath,
		"output-path", filepath.Join("."), "output path, default is '.'")
	addServiceCmd.PersistentFlags().BoolVar(&override,
		"override", false, "If override when file exists")
	addServiceCmd.PersistentFlags().StringVar(&protoFilePath,
		"path", "", "A .proto file to generate codes")
	addServiceCmd.PersistentFlags().StringVar(&protocPath,
		"protoc-path", "protoc", "protoc binary path, in case user has multi versions of protoc")
	addServiceCmd.PersistentFlags().StringArrayVar(&protoFileIncludePath,
		"proto-include-path", []string{}, "proto files include path used by protoc's command '--proto_path'")

	updateServiceCmd.PersistentFlags().StringVar(&outputPath,
		"output-path", filepath.Join("."), "output path, default is '.'")
	updateServiceCmd.PersistentFlags().BoolVar(&override,
		"override", false, "If override when file exists")
	updateServiceCmd.PersistentFlags().StringVar(&protoFilePath,
		"path", "", "A .proto file to generate codes")
	updateServiceCmd.PersistentFlags().StringVar(&protocPath,
		"protoc-path", "protoc", "protoc binary path, in case user has multi versions of protoc")
	updateServiceCmd.PersistentFlags().StringArrayVar(&protoFileIncludePath,
		"proto-include-path", []string{}, "proto files path include used by protoc's command '--proto_path'")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
