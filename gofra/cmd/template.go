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
	"github.com/spf13/cobra"
	"os"

	gofraTemplate "github.com/DarkMetrix/gofra/gofra/template"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Template operations [init]",
	Long: `Gofra is a framework using gRPC as the communication layer.\r\ntemplate command will help to generate template.json.default file.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// templateCmd represents the template command
var initTemplateCmd = &cobra.Command{
	Use:   "init",
	Short: "Template initialization, a template.json.default file will be generated",
	Long: `Gofra is a framework using gRPC as the communication layer.\r\ninit command will help to initialize template.json.default file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("====== Gofra service add ======")

		//Check path
		fmt.Printf("\r\nGet Working Path ......")
		workingPath, err := os.Getwd()

		if err != nil {
			fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
			return
		} else {
			fmt.Printf(" success! \r\nWorking path:%v\r\n", workingPath)
		}

		//Generate template.json
		fmt.Printf("\r\nGenerating template.json ...... \r\n")

		var jsonInfo gofraTemplate.JsonInfo
		fmt.Print("Author Name:")
		fmt.Scanln(&jsonInfo.Author)
		fmt.Print("Project Name:")
		fmt.Scanln(&jsonInfo.Project)
		fmt.Print("Project Address:")
		fmt.Scanln(&jsonInfo.Addr)

		err = gofraTemplate.GenerateTemplateJsonFile(workingPath, override, jsonInfo)

		if err != nil {
			fmt.Printf(" Generating template.json failed! \r\nerror:%v\r\n", err.Error())
			return
		} else {
			fmt.Printf(" Generating template.json success! \r\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(initTemplateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templateCmd.PersistentFlags().String("foo", "", "A help for foo")
	initTemplateCmd.PersistentFlags().BoolVar(&override, "override", false,"If override when file exists")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
