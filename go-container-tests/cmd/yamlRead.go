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
	"io/ioutil"

	"github.com/lucassabreu/go-container/def"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// yamlReadCmd represents the yamlRead command
var yamlReadCmd = &cobra.Command{
	Use:   "yamlRead",
	Short: "Reads a YAML to the def.Config struct",
	Run: func(cmd *cobra.Command, args []string) {
		var s def.Container

		yaml.Unmarshal([]byte(`services: { Test: {factory: banana}}`), &s)

		fmt.Printf("%#v\n", len(s.Services))

		for _, filename := range args {
			importYaml(filename)
		}
	},
}

func importYaml(filename string) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var c struct {
		Container def.Container
	}
	err = yaml.Unmarshal(dat, &c)
	if err != nil {
		panic(err)
	}

	fmt.Println(*c.Container.Packages[1].Alias)
}

func init() {
	RootCmd.AddCommand(yamlReadCmd)
}
