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
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	
	"github.com/spf13/cobra"
	"github.com/kirsle/configdir"
)

func decodeFile(machine, fname string) (*pem.Block, error) {
	configPath := configdir.LocalConfig("dm")
	filename := filepath.Join(configPath, "machines", machine, fname)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	return block, nil
}

type Config struct {
	CA       *pem.Block
	Cert     *pem.Block
	Key      *pem.Block
	Settings struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
}

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export a machine configuration",
	Long: `export the machine configuration to a PEM file. this configuration can be imported later with the import command.

for example:
    bash$ dm export MACHINE_NAME > MACHINE_NAME.pem
`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := configdir.LocalConfig("dm")
		configPathMachines := filepath.Join(configPath, "machines")
		s := &Config{}
		var err error
		
		
		s.CA, err = decodeFile(args[0], "ca.pem")
		if err != nil {
			log.Fatal(err)
		}
		s.Cert, err = decodeFile(args[0], "cert.pem")
		if err != nil {
			log.Fatal(err)
		}
		s.Key, err = decodeFile(args[0], "key.pem")
		if err != nil {
			log.Fatal(err)
		}
		
		data, err := ioutil.ReadFile(filepath.Join(configPathMachines,
			args[0], "config.json"))
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(data, &s.Settings); err != nil {
			log.Fatal(err)
		}

		s.CA.Type = "CERTIFICATE AUTHORITY"
		if err := pem.Encode(os.Stdout, s.CA); err != nil {
			log.Fatal(err)
		}
		if err := pem.Encode(os.Stdout, s.Cert); err != nil {
			log.Fatal(err)
		}
		if err := pem.Encode(os.Stdout, s.Key); err != nil {
			log.Fatal(err)
		}

		data, err = json.Marshal(&s.Settings)
		if err != nil {
			log.Fatal(err)
		}
		cfg := &pem.Block{
			Type:  "SETTINGS",
			Bytes: data,
		}
		if err := pem.Encode(os.Stdout, cfg); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
