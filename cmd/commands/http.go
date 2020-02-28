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

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/cobra"

	httpTemplate "github.com/DarkMetrix/gofra/internal/pkg/template/gin"
)

// serviceCmd represents the service command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Http operations [add]",
	Long: `Gofra is a framework using gRPC/gin as the communication layer.\r\nhttp command will help to generate http frame & handler.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// addServiceCmd represents the service command
var addHttpHandlerCmd = &cobra.Command{
	Use:   "add",
	Short: "Add http handler to project",
	Long: `Gofra is a framework using gRPC/gin as the communication layer.\r\nadd command will help to generate http frame & handler.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := validation.Validate(&uri, validation.Required, is.RequestURI)

		if err != nil {
			fmt.Printf("URI invalid! uri:%v, error:%v", uri, err.Error())
			return
		}

		err = validation.Validate(&method, validation.Required, validation.In("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"))

		if err != nil {
			fmt.Printf("Method invalid! method:%v, error:%v", method, err.Error())
			return
		}

		addHttpHandler(override)
	},
}

var uri string
var method string

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.AddCommand(addHttpHandlerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")
	addHttpHandlerCmd.PersistentFlags().StringVar(&uri, "uri", "","Http URI, e.g.:'/health'")
	addHttpHandlerCmd.PersistentFlags().StringVar(&method, "method", "GET","Http method, e.g.:'GET, POST, PUT, PATCH, DELETE and OPTIONS'")
	addHttpHandlerCmd.PersistentFlags().BoolVar(&override, "override", false,"If override when file exists")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addHttpHandler(override bool) error {
	fmt.Println("====== Gofra http add ======")

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

	//Check server type
	if templateInfo.Type != "http" {
		fmt.Printf(" failed! \r\nerror:Server type is not 'http'!\r\n")
		return err
	}

	//Generate http handler
	fmt.Printf("\r\nGenerating http handler code ......")
	err = httpTemplate.GenerateHttpHandler(workingPath, templateInfo, uri, method, override)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	return nil
}
