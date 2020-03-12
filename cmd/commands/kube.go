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
	"path/filepath"

	"github.com/spf13/cobra"

	kubeTemplate "github.com/DarkMetrix/gofra/internal/pkg/template/kubenetes"
	commonUtils "github.com/DarkMetrix/gofra/pkg/utils"
)

// kubeCmd represents the kube command
var kubeCmd = &cobra.Command{
	Use:   "kube",
	Short: "kubenetes operations [deployment, service]",
	Long: `Gofra is a framework using gRPC/gin as the communication layer. 
kube command will help to generate kubernetes deployment, service yaml file and configmap command shell script.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// deploymentKubeCmd represents the kube deployment command
var deploymentKubeCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Add generated kubernetes deployment.yml to project",
	Long: `Gofra is a framework using gRPC/gin as the communication layer. 
kube deployment command will help to generate kubernetes deployment file.`,
	Run: func(cmd *cobra.Command, args []string) {
		deploymentKube(override)
	},
}

// serviceKubeCmd represents the kube service command
var serviceKubeCmd = &cobra.Command{
	Use:   "service",
	Short: "Add generated kubernetes service.yml to project",
	Long: `Gofra is a framework using gRPC/gin as the communication layer. 
kube service command will help to generate kubernetes service yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceKube(override)
	},
}

// configmapKubeCmd represents the kube configmap command
var configmapKubeCmd = &cobra.Command{
	Use:   "configmap",
	Short: "Add generated kubernetes configmap.yml to project",
	Long: `Gofra is a framework using gRPC/gin as the communication layer. 
kube configmap command will help to generate kubernetes configmap shell script.
configmap.sh offers 'create', 'update', 'delete' and 'get' commands to help manage the configmap`,
	Run: func(cmd *cobra.Command, args []string) {
		configmapKube(override)
	},
}

func init() {
	rootCmd.AddCommand(kubeCmd)
	kubeCmd.AddCommand(deploymentKubeCmd)
	kubeCmd.AddCommand(serviceKubeCmd)
	kubeCmd.AddCommand(configmapKubeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	deploymentKubeCmd.PersistentFlags().BoolVar(&override, "override", false,"If override when file exists")

	serviceKubeCmd.PersistentFlags().BoolVar(&override, "override", false,"If override when file exists")

	configmapKubeCmd.PersistentFlags().BoolVar(&override, "override", false,"If override when file exists")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func deploymentKube(override bool) error {
	fmt.Println("====== Gofra kubernetes deployment ======")

	// check path
	fmt.Printf("\r\nChecking Path ......")
	workingPath, err := os.Getwd()

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\nWorking path:%v\r\n", workingPath)
	}

	// read template
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

	// mkdir
	fmt.Printf("\r\nMake dir ......")
	kubernetesPath := filepath.Join(workingPath, "kubernetes")

	commonUtils.CreatePath(kubernetesPath, false)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	// input image path
	var imagePath string
	fmt.Print("Image path:")
	fmt.Scanln(&imagePath)

	// generate deployment yaml file
	fmt.Printf("\r\nGenerating kubernetes deployment yaml file ......")
	err = kubeTemplate.GenerateKubeDeploymentYAMLFile(workingPath, imagePath, templateInfo, override)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	return nil
}

func serviceKube(override bool) error {
	fmt.Println("====== Gofra kubernetes service ======")

	// check path
	fmt.Printf("\r\nChecking Path ......")
	workingPath, err := os.Getwd()

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\nWorking path:%v\r\n", workingPath)
	}

	// read template
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

	// mkdir
	fmt.Printf("\r\nMake dir ......")
	kubernetesPath := filepath.Join(workingPath, "kubernetes")

	commonUtils.CreatePath(kubernetesPath, false)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	// generate service yaml file
	fmt.Printf("\r\nGenerating kubernetes service yaml file ......")
	err = kubeTemplate.GenerateKubeServiceYAMLFile(workingPath, templateInfo, override)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	return nil
}

func configmapKube(override bool) error {
	fmt.Println("====== Gofra kubernetes configmap ======")

	// check path
	fmt.Printf("\r\nChecking Path ......")
	workingPath, err := os.Getwd()

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\nWorking path:%v\r\n", workingPath)
	}

	// read template
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

	// mkdir
	fmt.Printf("\r\nMake dir ......")
	kubernetesPath := filepath.Join(workingPath, "kubernetes")

	commonUtils.CreatePath(kubernetesPath, false)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	// generate service yaml file
	fmt.Printf("\r\nGenerating kubernetes service yaml file ......")
	err = kubeTemplate.GenerateKubeConfigmapYAMLFile(workingPath, templateInfo, override)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	return nil
}
