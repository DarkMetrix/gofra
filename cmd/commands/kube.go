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

// kubeCmd represents the kube command
var kubeCmd = &cobra.Command{
	Use:   "kube",
	Short: "kubenetes operations [deployment, service]",
	Long: `Gofra is a framework using gRPC as the communication layer. 
kube command will help to generate kubernetes deployment, service yaml file and configmap command shell script.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// deploymentKubeCmd represents the kube deployment command
var deploymentKubeCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Add generated kubernetes deployment.yml to project",
	Long: `Gofra is a framework using gRPC as the communication layer. 
kube deployment command will help to generate kubernetes deployment file.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("====== gofra kube deployment ======")

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
			option.WithImagePath(image),
			option.WithPort(port),
		}

		// init docker file
		layout := directory.NewGRPCLayout(opts...)
		if err := generate.InitKubeDeployment(layout, opts...); err != nil {
			log.Fatalf("generate.InitDockerFile failed: %+v", err)
		}
	},
}

// serviceKubeCmd represents the kube service command
var serviceKubeCmd = &cobra.Command{
	Use:   "service",
	Short: "Add generated kubernetes service.yml to project",
	Long: `Gofra is a framework using gRPC as the communication layer. 
kube service command will help to generate kubernetes service yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("====== gofra kube service ======")

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
			option.WithPort(port),
			option.WithTargetPort(port),
		}

		// init docker file
		layout := directory.NewGRPCLayout(opts...)
		if err := generate.InitKubeService(layout, opts...); err != nil {
			log.Fatalf("generate.InitDockerFile failed: %+v", err)
		}
	},
}

// configmapKubeCmd represents the kube configmap command
var configmapKubeCmd = &cobra.Command{
	Use:   "configmap",
	Short: "Add generated kubernetes configmap.yml to project",
	Long: `Gofra is a framework using gRPC as the communication layer. 
kube configmap command will help to generate kubernetes configmap shell script.
configmap.sh offers 'create', 'update', 'delete' and 'get' commands to help manage the configmap`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var (
	namespace string
	version   string
	image     string
	port      string
)

func init() {
	rootCmd.AddCommand(kubeCmd)
	kubeCmd.AddCommand(deploymentKubeCmd)
	kubeCmd.AddCommand(serviceKubeCmd)
	kubeCmd.AddCommand(configmapKubeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	deploymentKubeCmd.PersistentFlags().StringVar(&outputPath,
		"output-path", filepath.Join("."), "output path, default is '.'")
	deploymentKubeCmd.PersistentFlags().StringVar(&project,
		"project", "", "project name, it will used as the ENTRYPOINT, e.g.: ENTRYPOINT ./application/bin/xxx"+
			"if project is not specified, will try to look up from go.mod from output path.")
	deploymentKubeCmd.PersistentFlags().BoolVar(&override, "override", false, "If override when file exists")
	deploymentKubeCmd.PersistentFlags().StringVar(&namespace,
		"namespace", "default", "Kubernetes namespace, default is 'default'")
	deploymentKubeCmd.PersistentFlags().StringVar(&version,
		"version", "v0.0.1", "Kubernetes version, default is 'v0.1.1'")
	deploymentKubeCmd.PersistentFlags().StringVar(&image, "image", "", "Kubernetes namespace, default is ''")
	deploymentKubeCmd.PersistentFlags().StringVar(&port, "port", "6666", "Kubernetes namespace, default is '6666'")

	serviceKubeCmd.PersistentFlags().StringVar(&outputPath,
		"output-path", filepath.Join("."), "output path, default is '.'")
	serviceKubeCmd.PersistentFlags().StringVar(&project,
		"project", "", "project name, it will used as the ENTRYPOINT, e.g.: ENTRYPOINT ./application/bin/xxx"+
			"if project is not specified, will try to look up from go.mod from output path.")
	serviceKubeCmd.PersistentFlags().BoolVar(&override, "override", false, "If override when file exists")
	serviceKubeCmd.PersistentFlags().StringVar(&namespace,
		"namespace", "default", "Kubernetes namespace, default is 'default'")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
