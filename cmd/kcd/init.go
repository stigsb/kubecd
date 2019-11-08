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
	"github.com/spf13/cobra"

	"github.com/zedge/kubecd/pkg/model"
	"github.com/zedge/kubecd/pkg/operations"
)

var initCluster string
var initContextsOnly bool
var initDryRun bool
var initGitlabMode bool

var initCmd = &cobra.Command{
	Use:   "init [ENV]",
	Short: "Initialize credentials and contexts",
	Long:  ``,
	Args:  clusterFlagOrEnvArg(&initCluster),
	RunE: func(cmd *cobra.Command, args []string) error {
		kcdConfig, err := model.NewConfigFromFile(environmentsFile)
		if err != nil {
			return err
		}
		ops, err := buildInitOperations(kcdConfig, initCluster, initDryRun, initGitlabMode, args)
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
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&initCluster, "cluster", "", "Initialize contexts for all environments in a cluster")
	initCmd.Flags().BoolVar(&initContextsOnly, "contexts-only", false, "initialize contexts only, assuming that cluster credentials are set up")
	initCmd.Flags().BoolVarP(&initDryRun, "dry-run", "n", false, "print commands instead of running them")
	initCmd.Flags().BoolVar(&initGitlabMode, "gitlab", false, "grab kube config from GitLab environment")
}

func buildInitOperations(kcdConfig *model.KubeCDConfig, cluster string, dryRun, gitLab bool, args []string) ([]operations.Operation, error) {
	envsToInit, err := environmentsFromArgs(kcdConfig, cluster, args)
	if err != nil {
		return nil, err
	}
	ops := make([]operations.Operation, 0)
	op := operations.NewHelmInit(kcdConfig, dryRun)
	if err := op.Prepare(); err != nil {
		return nil, err
	}
	ops = append(ops, op)
	for _, env := range envsToInit {
		op := operations.NewEnvInit(env, dryRun, gitLab)
		if err := op.Prepare(); err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}
	return ops, nil
}
