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
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func decodeFile(machine, fname string) (*pem.Block, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s",
		"/home/pwhelan/.config/dm/machines", machine, fname))
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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
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
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s",
			"/home/pwhelan/.config/dm/machines", args[0], "config.json"))
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
