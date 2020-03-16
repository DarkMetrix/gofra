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

	istioTemplate "github.com/DarkMetrix/gofra/internal/pkg/template/istio"
	commonUtils "github.com/DarkMetrix/gofra/pkg/utils"
)

// istioCmd represents the istio command
var istioCmd = &cobra.Command{
	Use:   "istio",
	Short: "istio operations [virtual-service, destination-rule]",
	Long: `Gofra is a framework using gRPC/gin as the communication layer. 
istio command will help to generate istio virtual service and destination rule yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// virtualServiceCmd represents the istio virtual-service command
var virtualServiceCmd = &cobra.Command{
	Use:   "virtual-service",
	Short: "Add generated istio virtual-service.yml to project",
	Long: `Gofra is a framework using gRPC/gin as the communication layer. 
istio virtual-service command will help to generate istio virtual service yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		virtualServiceIstio(namespace, override)
	},
}

// destinationRuleCmd represents the istio destination-rule command
var destinationRuleCmd = &cobra.Command{
	Use:   "destination-rule",
	Short: "Add generated istio destination-rule.yml to project",
	Long: `Gofra is a framework using gRPC/gin as the communication layer. 
istio destination-rule command will help to generate istio destination rule yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		destinationRuleIstio(namespace, override)
	},
}

func init() {
	rootCmd.AddCommand(istioCmd)
	istioCmd.AddCommand(virtualServiceCmd)
	istioCmd.AddCommand(destinationRuleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	virtualServiceCmd.PersistentFlags().BoolVar(&override, "override", false,"If override when file exists")
	virtualServiceCmd.PersistentFlags().StringVar(&namespace, "namespace", "","Kubernetes namespace, default is ''")

	destinationRuleCmd.PersistentFlags().BoolVar(&override, "override", false,"If override when file exists")
	destinationRuleCmd.PersistentFlags().StringVar(&namespace, "namespace", "","Kubernetes namespace, default is ''")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func virtualServiceIstio(namespace string, override bool) error {
	fmt.Println("====== Gofra istio virtual-service ======")

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
	istioPath := filepath.Join(workingPath, "istio")

	commonUtils.CreatePath(istioPath, false)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	// generate virtual service yaml file
	fmt.Printf("\r\nGenerating istio virtual service yaml file ......")
	err = istioTemplate.GenerateIstioVirtaulServiceYAMLFile(workingPath, namespace, templateInfo, override)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	return nil
}

func destinationRuleIstio(namespace string, override bool) error {
	fmt.Println("====== Gofra istio destination-rule ======")

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
	istioPath := filepath.Join(workingPath, "istio")

	commonUtils.CreatePath(istioPath, false)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	// generate virtual service yaml file
	fmt.Printf("\r\nGenerating istio destination rule yaml file ......")
	err = istioTemplate.GenerateIstioDestinationRuleYAMLFile(workingPath, namespace, templateInfo, override)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	return nil
}

