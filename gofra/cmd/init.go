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
	"os"
	"errors"
	"strings"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"

	gofraUtils "github.com/DarkMetrix/gofra/gofra/utils"
	gofraTemplate "github.com/DarkMetrix/gofra/gofra/template"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gofra application",
	Long: `Gofra is a framework using gRPC as the communication layer.
		   init command will create basic framework structure.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("====== Gofra init ======")

		//Check path
		fmt.Printf("\r\nChecking Path ......")
		goPath, workingPath, err := CheckPath()

		if err != nil {
			fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
			return
		} else {
			fmt.Printf(" success! \r\nGOPATH:%v\r\nWorking path:%v\r\n", goPath, workingPath)
		}

		//Read template
		fmt.Printf("\r\nReading template ......")
		if len(templatePath) == 0 {
			fmt.Printf(" failed! \r\nerror:Template file path is empty!\r\n")
			return
		}

		templateInfo, jsonBuffer, err := ReadTemplate(templatePath)

		if err != nil {
			fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
			return
		} else {
			fmt.Printf(" success! \r\nTemplate:\r\n%v\r\n", jsonBuffer)
		}

		//Init directory structure
		fmt.Printf("\r\nInitializing directory structure ......")
		err = InitDirectoryStructure(workingPath, templateInfo)

		if err != nil {
			fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
			return
		} else {
			fmt.Printf(" success!\r\n")
		}

		//Init all files
		fmt.Printf("\r\nInitializing all files ......")
		err = InitAllFiles(workingPath, templateInfo)

		if err != nil {
			fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
			return
		} else {
			fmt.Printf(" success!\r\n")
		}

		//Print application directory structure
		fmt.Printf("\r\nApplication '%v' directory structure", templateInfo.Project)
		filepath.Walk(workingPath, func(path string, info os.FileInfo, err error) error {
			relPath, err := filepath.Rel(workingPath, path)

			if err != nil {
				return err
			}

			fmt.Println(relPath)
			return nil
		})
	},
}

var templatePath string
var destPath string

type TemplateInfo struct {
	Project string
	Version string
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	initCmd.PersistentFlags().StringVar(&templatePath, "template_path", "./template.json", "A template file in json to tell how to generate codes")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//Check if current working directory is under gopath
func CheckPath() (string, string, error) {
	goPath := os.Getenv("GOPATH")

	if len(goPath) == 0 {
		return "", "", errors.New("GOPATH is not set!")
	}

	workingPath, err := os.Getwd()

	if err != nil {
		return "", "", err
	}

	isMatch := strings.HasPrefix(workingPath, filepath.Join(goPath, "src"))

	if !isMatch {
		return "", "", errors.New("The working path is not under $GOPATH/src!")
	}

	return goPath, workingPath, nil
}

//Read template json file to ge information about how to generate the application
func ReadTemplate(templatePath string) (*TemplateInfo, string, error) {
	templateFile, err := os.Open(templatePath)

	if err != nil {
		return nil, "", err
	}

	defer templateFile.Close()

	data, err := ioutil.ReadAll(templateFile)

	if err != nil {
		return nil, "", err
	}

	var info *TemplateInfo = new(TemplateInfo)

	err = json.Unmarshal(data, info)

	if err != nil {
		return nil, "", err
	}

	return info, string(data), nil
}

//Init application directory structure
func InitDirectoryStructure(workingPath string, info *TemplateInfo) error {
	binPath := filepath.Join(workingPath, "bin")
	confPath := filepath.Join(workingPath, "conf")
	logPath := filepath.Join(workingPath, "log")
	srcPath := filepath.Join(workingPath, "src")

	applicationPath := filepath.Join(workingPath, "src", "application")
	commonPath := filepath.Join(workingPath, "src", "common")
	configPath := filepath.Join(workingPath, "src",  "config")
	handlerPath := filepath.Join(workingPath, "src", "handler")

	//Create root directories
	err := gofraUtils.CreatePaths(binPath, confPath, logPath, srcPath)

	if err != nil {
		return err
	}

	//Create src sub directories
	err = gofraUtils.CreatePaths(applicationPath, commonPath, configPath, handlerPath)

	if err != nil {
		return err
	}

	return nil
}

//Init all go file with template
func InitAllFiles(workingPath string, info *TemplateInfo) error {
	err := GenerateCommonFile(workingPath, info)

	if err != nil {
		return err
	}

	err = GenerateConfigFile(workingPath, info)

	if err != nil {
		return err
	}

	err = GenerateApplicationFile(workingPath, info)

	if err != nil {
		return err
	}

	err = GenerateMainFile(workingPath, info)

	if err != nil {
		return err
	}

	err = GenerateHandler(workingPath, info)

	if err != nil {
		return err
	}

	return nil
}

//Generate common.go
func GenerateCommonFile(workingPath string, info *TemplateInfo) error {
	filePath := filepath.Join(workingPath, "src", "common", "common.go")

	//Check file is exist or not
	isExist, err := gofraUtils.CheckPathExists(filePath)

	if err != nil {
		return err
	}

	if isExist {
		filePathRel, err := filepath.Rel(workingPath, filePath)

		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
	}

	//Parse template
	commonTemplate, err := template.New("common").Parse(gofraTemplate.CommonTemplate)

	if err != nil {
		return err
	}

	commonInfo := &gofraTemplate.CommonInfo{
		Project: info.Project,
		Version: info.Version,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = commonTemplate.Execute(file, commonInfo)

	if err != nil {
		return err
	}

	return nil
}

//Generate config.go
func GenerateConfigFile(workingPath string, info *TemplateInfo) error {
	return nil
}

//Generate application.go
func GenerateApplicationFile(workingPath string, info *TemplateInfo) error {
	return nil
}

//Generate main.go
func GenerateMainFile(workingPath string, info *TemplateInfo) error {
	return nil
}

//Generate handler
func GenerateHandler(workingPath string, info *TemplateInfo) error {
	return nil
}
