/*
Copyright © 2019 Zedge, Inc.

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
package main

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var yaml2jsonCmd = &cobra.Command{
	Use:   "yaml2json",
	Short: "YAML to JSON conversion utility (stdin/stdout)",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var rawObject interface{}
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return errors.Wrap(err, `reading stdin`)
		}
		err = yaml.Unmarshal(data, &rawObject)
		if err != nil {
			return errors.Wrap(err, `decoding JSON`)
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)
		if err = encoder.Encode(&rawObject); err != nil {
			return errors.Wrap(err, `encoding YAML`)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(yaml2jsonCmd)
}
