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

var applyReleases []string
var applyCluster string
var applyInit bool
var applyGitlab bool
var applyDryRun bool
var applyDebug bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply changes to Kubernetes",
	Long:  ``,
	Args:  clusterFlagOrEnvArg(&applyCluster),
	RunE: func(cmd *cobra.Command, args []string) error {
		kcdConfig, err := model.NewConfigFromFile(environmentsFile)
		if err != nil {
			return err
		}
		ops, err := buildApplyOperations(kcdConfig, args)
		if err != nil {
			return err
		}
		for _, op := range ops {
			//fmt.Println(op.String())
			if err := op.Execute(); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().BoolVarP(&applyDryRun, "dry-run", "n", false, "dry run mode, only print commands")
	applyCmd.Flags().BoolVar(&applyDebug, "debug", false, "run helm with --debug")
	applyCmd.Flags().StringSliceVarP(&applyReleases, "releases", "r", []string{}, "apply only these releases")
	applyCmd.Flags().StringVarP(&applyCluster, "cluster", "c", "", "apply all environments in CLUSTER")
	applyCmd.Flags().BoolVar(&applyInit, "init", false, "initialize credentials and contexts")
	applyCmd.Flags().BoolVar(&applyGitlab, "gitlab", false, "initialize in gitlab mode")
}

func buildApplyOperations(kcdConfig *model.KubeCDConfig, args []string) ([]operations.Operation, error) {
	environments, err := environmentsFromArgs(kcdConfig, applyCluster, args)
	if err != nil {
		return nil, err
	}
	ops := make([]operations.Operation, 0)
	if applyInit {
		initOps, err := buildInitOperations(kcdConfig, applyCluster, applyDryRun, applyGitlab, args)
		if err != nil {
			return nil, err
		}
		ops = append(ops, initOps...)
	}
	for _, env := range environments {
		releases := make([]*model.Release, 0)
		if len(applyReleases) > 0 {
			for _, relName := range applyReleases {
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
			op := operations.NewApply(release, applyDryRun, applyDebug)
			if err := op.Prepare(); err != nil {
				return nil, err
			}
			ops = append(ops, op)
		}
	}
	return ops, nil
}
