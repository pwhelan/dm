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
	"io/ioutil"
	"log"
	"os"

	"github.com/pwhelan/dm/machine"

	"github.com/spf13/cobra"
)

func displaymachine(file os.FileInfo) {
	var machine machine.Machine
	if err := machine.Read(file); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\t%s\n", machine.Name, machine.URL)
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list machines.",
	Long: `List all of the currently configured machines.`,
	Run: func(cmd *cobra.Command, args []string) {
		files, err := ioutil.ReadDir("/home/pwhelan/.config/dm/machines")
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if file.IsDir() == true {
				displaymachine(file)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
