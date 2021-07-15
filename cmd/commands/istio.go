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
	"path"
	"path/filepath"

	"github.com/DarkMetrix/gofra/internal/pkg/directory"
	"github.com/DarkMetrix/gofra/internal/pkg/generate"
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/utils/gomod"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// istioCmd represents the istio command
var istioCmd = &cobra.Command{
	Use:   "istio",
	Short: "istio operations [virtual-service, destination-rule]",
	Long: `Gofra is a framework using gRPC as the communication layer. 
istio command will help to generate istio virtual service and destination rule yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// virtualServiceCmd represents the istio virtual-service command
var virtualServiceCmd = &cobra.Command{
	Use:   "virtual-service",
	Short: "Add generated istio virtual-service.yml to project",
	Long: `Gofra is a framework using gRPC as the communication layer. 
istio virtual-service command will help to generate istio virtual service yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("====== gofra istio virtual-service ======")

		// get project name
		if project == "" {
			// get project name from go.mod
			goModule, err := gomod.GetGoModule(filepath.Join(outputPath, "go.mod"))
			if err != nil {
				log.Fatalf("utils.GetGoModule failed! error:%v", err)
			}
			project = path.Base(goModule)
		}

		opts := []option.Option{
			option.WithOutputPath(outputPath),
			option.WithOverride(override),
			option.WithProject(project),
			option.WithNamespace(namespace),
			option.WithVersion(version),
			option.WithPort(port),
		}

		// init istio virtual_service.yaml
		layout := directory.NewGRPCLayout(opts...)
		if err := generate.InitIstioVirtualService(layout, opts...); err != nil {
			log.Fatalf("generate.InitIstioVirtualService failed: %+v", err)
		}
	},
}

// destinationRuleCmd represents the istio destination-rule command
var destinationRuleCmd = &cobra.Command{
	Use:   "destination-rule",
	Short: "Add generated istio destination-rule.yml to project",
	Long: `Gofra is a framework using gRPC as the communication layer. 
istio destination-rule command will help to generate istio destination rule yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("====== gofra istio destination-rule ======")

		// get project name
		if project == "" {
			// get project name from go.mod
			goModule, err := gomod.GetGoModule(filepath.Join(outputPath, "go.mod"))
			if err != nil {
				log.Fatalf("utils.GetGoModule failed! error:%v", err)
			}
			project = path.Base(goModule)
		}

		opts := []option.Option{
			option.WithOutputPath(outputPath),
			option.WithOverride(override),
			option.WithProject(project),
			option.WithNamespace(namespace),
			option.WithVersion(version),
		}

		// init istio destination_rule.yaml
		layout := directory.NewGRPCLayout(opts...)
		if err := generate.InitIstioDestinationRule(layout, opts...); err != nil {
			log.Fatalf("generate.InitIstioDestinationRule failed: %+v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(istioCmd)
	istioCmd.AddCommand(virtualServiceCmd)
	istioCmd.AddCommand(destinationRuleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	virtualServiceCmd.PersistentFlags().StringVar(&outputPath,
		"output-path", filepath.Join("."), "output path, default is '.'")
	virtualServiceCmd.PersistentFlags().StringVar(&project,
		"project", "", "project name it will be used as the metadata.name in deployment, "+
			"if project is not specified, will try to look up from go.mod from output path.")
	virtualServiceCmd.PersistentFlags().BoolVar(&override, "override", false, "If override when file exists")
	virtualServiceCmd.PersistentFlags().StringVar(&namespace,
		"namespace", "default", "Kubernetes namespace, default is 'default'")
	virtualServiceCmd.PersistentFlags().StringVar(&version,
		"version", "v1", "Kubernetes version, default is 'v1'")
	virtualServiceCmd.PersistentFlags().StringVar(&port, "port", "6666", "Kubernetes namespace, default is '6666'")

	destinationRuleCmd.PersistentFlags().StringVar(&outputPath,
		"output-path", filepath.Join("."), "output path, default is '.'")
	destinationRuleCmd.PersistentFlags().StringVar(&project,
		"project", "", "project name it will be used as the metadata.name in deployment, "+
			"if project is not specified, will try to look up from go.mod from output path.")
	destinationRuleCmd.PersistentFlags().BoolVar(&override, "override", false, "If override when file exists")
	destinationRuleCmd.PersistentFlags().StringVar(&namespace,
		"namespace", "default", "Kubernetes namespace, default is 'default'")
	destinationRuleCmd.PersistentFlags().StringVar(&version,
		"version", "v1", "Kubernetes version, default is 'v1'")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
