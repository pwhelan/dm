/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/pwhelan/dm/machine"

	"github.com/spf13/cobra"
	"github.com/kirsle/configdir"
)

var unset bool

func detectShell() string {
	return path.Base(os.Getenv("SHELL"))
}

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "generate env config for a machine.",
	Long:  `generate the environment variables to configure the use of a remote docker instance.

an example:
    bash$ eval $(dm env MACHINE_NAME)

this would set the variables to access the docker machine MACHINE_NAME.`,
	Run: func(cmd *cobra.Command, args []string) {
		if (unset == true) {
			switch detectShell() {
			case "fish":
				fmt.Println("set -e DOCKER_TLS_VERIFY;")
				fmt.Println("set -e DOCKER_HOST;")
				fmt.Println("set -e DOCKER_CERT_PATH;")
				fmt.Println("set -e DOCKER_MACHINE_NAME;")
			case "bash":
				fmt.Println("unset DOCKER_TLS_VERIFY;")
				fmt.Println("unset DOCKER_HOST;")
				fmt.Println("unset DOCKER_CERT_PATH;")
				fmt.Println("unset DOCKER_MACHINE_NAME;")
			default:
				log.Fatal("unable to detect shell")
			}
			return
		}
		
		var mach machine.Machine
		configPath := configdir.LocalConfig("dm")
		configDir := filepath.Join(configPath, "machines", args[0])
		file, err := os.Stat(configDir)
		if err != nil {
			log.Fatal(err)
		}
		if err := mach.Read(file); err != nil {
			log.Fatal(err)
		}
		
		switch detectShell() {
		case "fish":
			fmt.Printf("set -gx DOCKER_TLS_VERIFY \"1\";\n")
			fmt.Printf("set -gx DOCKER_HOST \"%s\";\n", mach.URL)
			fmt.Printf("set -gx DOCKER_CERT_PATH \"%s\";\n", configDir)
			fmt.Printf("set -gx DOCKER_MACHINE_NAME \"%s\";\n",
				mach.Name)
		case "bash":
			fmt.Printf("export DOCKER_TLS_VERIFY=\"1\";\n")
			fmt.Printf("export DOCKER_HOST=\"%s\";\n", mach.URL)
			fmt.Printf("export DOCKER_CERT_PATH=\"%s\";\n",
				configDir)
			fmt.Printf("export DOCKER_MACHINE_NAME=\"%s\";\n",
				mach.Name)
		default:
			log.Fatal("unable to detect shell")
		}
	},
}

func init() {
	envCmd.Flags().BoolVarP(&unset, "unset", "u", false, "unset environment variables for docker")
	rootCmd.AddCommand(envCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
