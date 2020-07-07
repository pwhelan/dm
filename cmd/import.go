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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/kirsle/configdir"
)

func dumpFile(machine, file string, data []byte) error {
	configPath := configdir.LocalConfig("dm")
	configFile := filepath.Join(configPath, "machines", machine, file)
	return ioutil.WriteFile(configFile, data, 0600)
}

func dumpBlock(machine, file string, data *pem.Block) error {
	configPath := configdir.LocalConfig("dm")
	configFile := filepath.Join(configPath, "machines", machine, file)
	fd, err := os.Create(configFile)
	if err != nil {
		log.Fatal(err)
	}
	return pem.Encode(fd, data)
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import a machine configuration",
	Long: `import the machine configuration from a PEM file. this configuration can be imported from a previously exported configuration using the export command.

for example:
    bash$ dm import MACHINE_NAME < MACHINE_NAME.pem
`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := configdir.LocalConfig("dm")
		configPathMachines := filepath.Join(configPath, "machines")
		configPathMachine := filepath.Join(configPathMachines, args[0])
		var data []byte
		var block *pem.Block

		if _, err := os.Stat(configPath); err != nil {
			if err := os.Mkdir(configPath, 0700); err != nil {
				log.Fatal(err)
			}
		}

		if _, err := os.Stat(configPathMachines); err != nil {
			if err := os.Mkdir(configPathMachines, 0700); err != nil {
				log.Fatal(err)
			}
		}

		if err := os.Mkdir(configPathMachine, 0700); err != nil {
			log.Fatal(err)
		}

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
