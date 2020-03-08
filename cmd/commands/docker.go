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
	"github.com/spf13/cobra"
	"os"

	dockerTemplate "github.com/DarkMetrix/gofra/internal/pkg/template/docker"
)

// dockerCmd represents the docker command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "docker operations [generate]",
	Long: `Gofra is a framework using gRPC/gin as the communication layer. 
docker command will help to generate docker file.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// generateCmd represents the docker generate command
var generateDockerCmd = &cobra.Command{
	Use:   "generate",
	Short: "Add generated docker file to project",
	Long: `Gofra is a framework using gRPC/gin as the communication layer. 
docker generate command will help to generate docker file.`,
	Run: func(cmd *cobra.Command, args []string) {
		generateDocker(override)
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)
	dockerCmd.AddCommand(generateDockerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	generateDockerCmd.PersistentFlags().BoolVar(&override, "override", false,"If override when file exists")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func generateDocker(override bool) error {
	fmt.Println("====== Gofra docker generate ======")

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

	//Generate docker file
	fmt.Printf("\r\nGenerating docker file ......")
	err = dockerTemplate.GenerateDockerFile(workingPath, templateInfo, override)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	return nil
}

