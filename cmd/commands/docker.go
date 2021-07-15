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
	"strings"

	"github.com/DarkMetrix/gofra/internal/pkg/directory"
	"github.com/DarkMetrix/gofra/internal/pkg/generate"
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/utils/gomod"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

// dockerCmd represents the docker command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "docker operations [generate]",
	Long: `Gofra is a framework using gRPC as the communication layer. 
docker command will help to generate docker file.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// generateCmd represents the docker generate command
var generateDockerCmd = &cobra.Command{
	Use:   "generate",
	Short: "Add generated docker file to project",
	Long: `Gofra is a framework using gRPC as the communication layer. 
docker generate command will help to generate docker file.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("====== gofra docker generate ======")

		// validate format of labels
		if err := checkLabels(labels); err != nil {
			log.Fatalf("checkLabels failed! error:%+v", err)
		}

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
			option.WithLabels(labels),
		}

		// init docker file
		layout := directory.NewGRPCLayout(opts...)
		if err := generate.InitDockerFile(layout, opts...); err != nil {
			log.Fatalf("generate.InitDockerFile failed: %+v", err)
		}
	},
}

// checkLabels checks if label format is valid
func checkLabels(labels []string) error {
	for _, label := range labels {
		kv := strings.Split(label, "=")
		if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
			return xerrors.Errorf("label not invalid! label:%v", label)
		}
	}
	return nil
}

var (
	project string
	labels  []string
)

func init() {
	rootCmd.AddCommand(dockerCmd)
	dockerCmd.AddCommand(generateDockerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	generateDockerCmd.PersistentFlags().StringVar(&outputPath,
		"output-path", filepath.Join("."), "output path, default is '.'")
	generateDockerCmd.PersistentFlags().StringVar(&project,
		"project", "", "project name, it will be used as the ENTRYPOINT, e.g.: ENTRYPOINT ./application/bin/xxx"+
			"if project is not specified, will try to look up from go.mod from output path.")
	generateDockerCmd.PersistentFlags().BoolVar(&override,
		"override", false, "If override when file exists")
	generateDockerCmd.PersistentFlags().StringSliceVar(&labels,
		"labels", []string{}, "Labels set to docker file, e.g.:labels=app=bar")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
