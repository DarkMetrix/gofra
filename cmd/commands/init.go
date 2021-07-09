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

	"github.com/DarkMetrix/gofra/internal/pkg/generate"
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gofra application using templates.json",
	Long:  `Gofra is a framework using gRPC/gin as the communication layer.\r\ninit command will create basic framework structure.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("====== gofra init ======")

		opts := []option.Option{
			option.WithOutputPath(outputPath),
			option.WithOverride(override),
			option.WithGoModule(goModule),
			option.WithGoVersion(goVersion),
			option.WithProtocPath(protocPath),
			option.WithProtoFileIncludePath(protoFileIncludePath),
		}

		// init directory structure
		log.Info("Initializing directory structure......")
		layout, err := generate.InitGRPCDirectoryStructure(opts...)
		if err != nil {
			log.Fatalf("generate.InitGRPCDirectoryStructure failed! error:%+v", err)
		}

		// init gRPC service
		log.Info("Initializing gRPC service......")
		if err := generate.InitGRPCService(layout, opts...); err != nil {
			log.Fatalf("generate.InitGRPCService failed! error:%+v", err)
		}
	},
}

var (
	outputPath           string
	protocPath           string
	protoFileIncludePath []string
	override             bool
	goModule             string
	goVersion            string
)

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	initCmd.PersistentFlags().StringVar(&outputPath, "output-path", filepath.Join("."), "output path, default is './'")
	initCmd.PersistentFlags().BoolVar(&override, "override", false, "If override when file exists")
	initCmd.PersistentFlags().StringVar(&protocPath, "protoc-path", "protoc", "protoc binary path, in case user has multi versions of protoc")
	initCmd.PersistentFlags().StringArrayVar(&protoFileIncludePath, "proto_include_path", []string{}, "proto files include path used by protoc's command '--proto_path'")
	initCmd.PersistentFlags().StringVar(&goModule, "go-module", "", "go module name, default is empty, e.g.:github.com/foo/bar")
	initCmd.PersistentFlags().StringVar(&goVersion, "go-version", "1.16", "go version, default is '1.16'")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
