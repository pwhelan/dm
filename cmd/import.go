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
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func dumpFile(machine, file string, data []byte) error {
	return ioutil.WriteFile(fmt.Sprintf("%s/%s/%s",
		"/home/pwhelan/.config/dm/machines", machine, file), data, 0600)
}

func dumpBlock(machine, file string, data *pem.Block) error {
	fd, err := os.Create(fmt.Sprintf("%s/%s/%s",
		"/home/pwhelan/.config/dm/machines", machine, file))
	if err != nil {
		log.Fatal(err)
	}
	return pem.Encode(fd, data)
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var data []byte
		var block *pem.Block

		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		for len(data) > 0 {
			block, data = pem.Decode(data)
			if block == nil {
				log.Fatal("failed to decode PEM block")
			}
			switch block.Type {
			case "SETTINGS":
				dumpFile(args[0], "config.json", block.Bytes)
			case "RSA PRIVATE KEY":
				dumpBlock(args[0], "key.pem", block)
			case "CERTIFICATE":
				dumpBlock(args[0], "cert.pem", block)
			case "CERTIFICATE AUTHORITY":
				block.Type = "CERTIFICATE"
				dumpBlock(args[0], "ca.pem", block)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
