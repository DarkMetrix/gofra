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
	"github.com/spf13/cobra"
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
		//addService(servicePath, false, false)
	},
}

// updateServiceCmd represents the service update command
var updateServiceCmd = &cobra.Command{
	Use:   "update",
	Short: "Update service (*.proto) to project",
	Long: `Gofra is a framework using gRPC/gin as the communication layer.
service update command will help to manipulate .proto file to update service frame & handler.`,
	Run: func(cmd *cobra.Command, args []string) {
		//addService(servicePath, override, true)
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
	addServiceCmd.PersistentFlags().StringArrayVar(&protoFileIncludePath, "proto_include_path", []string{}, "proto files include path used by protoc's command '--proto_path'")

	updateServiceCmd.PersistentFlags().StringVar(&servicePath, "path", "", "A .proto file to generate codes")
	updateServiceCmd.PersistentFlags().StringVar(&protocPath, "protoc_path", "protoc", "protoc binary path, in case user has multi versions of protoc")
	updateServiceCmd.PersistentFlags().StringArrayVar(&protoFileIncludePath, "proto_path", []string{}, "proto files path include used by protoc's command '--proto_path'")
	updateServiceCmd.PersistentFlags().BoolVar(&override, "override", false, "If override when file exists")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
