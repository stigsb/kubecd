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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zedge/kubecd/pkg/model"
	"github.com/zedge/kubecd/pkg/operations"
)

var (
	renderDryRun   bool
	renderReleases []string
	renderCluster  string
	renderGitlab   bool
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "show helm templates and plain kubernetes YAML resources",
	Long:  ``,
	Args:  clusterFlagOrEnvArg(&renderCluster),
	RunE: func(cmd *cobra.Command, args []string) error {
		kcdConfig, err := model.NewConfigFromFile(environmentsFile)
		if err != nil {
			return err
		}
		ops, err := buildRenderOperations(kcdConfig, args)
		if err != nil {
			return err
		}
		for _, op := range ops {
			fmt.Println(op.String())
			//if err := op.Execute(); err != nil {
			//	return err
			//}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().BoolVarP(&renderDryRun, "dry-run", "n", false, "dry run mode, only print commands")
	renderCmd.Flags().StringSliceVarP(&renderReleases, "releases", "r", []string{}, "generate template only these releases")
	renderCmd.Flags().StringVarP(&renderCluster, "cluster", "c", "", "template all environments in CLUSTER")
	renderCmd.Flags().BoolVar(&renderGitlab, "gitlab", false, "initialize in gitlab mode")
}

func buildRenderOperations(kcdConfig *model.KubeCDConfig, args []string) ([]operations.Operation, error) {
	environments, err := environmentsFromArgs(kcdConfig, applyCluster, args)
	if err != nil {
		return nil, err
	}
	ops := make([]operations.Operation, 0)
	for _, env := range environments {
		releases := make([]*model.Release, 0)
		if len(renderReleases) > 0 {
			for _, relName := range renderReleases {
				release := env.GetRelease(relName)
				if release == nil {
					return nil, fmt.Errorf("no release named %q in environment %q", relName, env.Name)
				}
				releases = append(releases, release)
			}
		} else {
			releases = append(releases, env.AllReleases()...)
		}
		for _, release := range releases {
			op := operations.NewRender(release, renderDryRun)
			if err := op.Prepare(); err != nil {
				return nil, err
			}
			ops = append(ops, op)
		}
	}
	return ops, nil

}
