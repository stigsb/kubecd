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
	"fmt"
	"github.com/kubecd/kubecd/pkg/model"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:       "list {env,release,cluster}",
	Short:     "list clusters, environments or releases",
	Long:      ``,
	Args:      matchAllArgs(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"env", "envs", "release", "releases", "cluster", "clusters"},
	RunE: func(cmd *cobra.Command, args []string) error {
		kcdConfig, err := model.NewConfigFromFile(environmentsFile)
		if err != nil {
			return err
		}
		switch args[0] {
		case "env", "envs":
			for _, env := range kcdConfig.Environments {
				fmt.Println(env.Name)
			}
		case "release", "releases":
			for _, env := range kcdConfig.Environments {
				for _, release := range env.AllReleases() {
					fmt.Printf("%s -r %s\n", env.Name, release.Name)
				}
			}
		case "cluster", "clusters":
			for _, cluster := range kcdConfig.AllClusters() {
				fmt.Println(cluster.Name)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
