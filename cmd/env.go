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

	"github.com/pwhelan/dm/machine"

	"github.com/spf13/cobra"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Generate env config for a machine.",
	Long: `WIP ...`,
	Run: func(cmd *cobra.Command, args []string) {
		var mach machine.Machine
		file, err := os.Stat(fmt.Sprintf("%s/%s",
			"/home/pwhelan/.config/dm/machines", args[0]))
		if err != nil {
			log.Fatal(err)
		}
		if err := mach.Read(file); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("set -gx DOCKER_TLS_VERIFY \"1\";\n")
		fmt.Printf("set -gx DOCKER_HOST \"%s\";\n", mach.URL)
		fmt.Printf("set -gx DOCKER_CERT_PATH \"%s/%s\";\n", 
			"/home/pwhelan/.config/dm/machines", file.Name())
		fmt.Printf("set -gx DOCKER_MACHINE_NAME \"%s\";\n", 
			mach.Name)
	},
}

func init() {
	rootCmd.AddCommand(envCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
