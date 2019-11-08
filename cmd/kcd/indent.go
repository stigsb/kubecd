/*
 * Copyright 2018-2020 Zedge, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package main

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/kubecd/kubecd/pkg/updates"
)

var indentCmd = &cobra.Command{
	Use:   "indent file [file...]",
	Short: "canonically indent YAML files",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, file := range args {
			var doc yaml.Node
			data, err := ioutil.ReadFile(file)
			if err != nil {
				return errors.Wrapf(err, `error reading %q`, file)
			}
			err = yaml.Unmarshal(data, &doc)
			if err != nil {
				return errors.Wrapf(err, `error decoding yaml in %q`, file)
			}
			if err = updates.WriteIndentedYamlToFile(file, &doc); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(indentCmd)
}
