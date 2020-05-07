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
	"github.com/kubecd/kubecd/pkg/operations"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/kubecd/kubecd/pkg/image"
	"github.com/kubecd/kubecd/pkg/model"
	"github.com/kubecd/kubecd/pkg/updates"
)

var (
	observePatch    bool
	observeReleases []string
	observeCluster  string
	observeImage    string
	observeChart    string
	observeVerify   bool
)

var observeCmd = &cobra.Command{
	Use:   "observe [ENV]",
	Short: "observe a new version of an image or chart",
	Long:  ``,
	Args:  matchAllArgs(clusterFlagOrEnvArg(&observeCluster), imageOrChart(&observeImage, &observeChart)),
	RunE: func(cmd *cobra.Command, args []string) error {
		kcdConfig, err := model.NewConfigFromFile(environmentsFile)
		if err != nil {
			return err
		}
		ops, err := buildObserveOperations(kcdConfig, args)
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

func buildObserveOperations(kcdConfig *model.KubeCDConfig, args []string) ([]operations.Operation, error) {
	//environments, err := environmentsFromArgs(kcdConfig, observeCluster, args)
	//if err != nil {
	//	return nil, err
	//}
	ops := make([]operations.Operation, 0)
	if observeVerify {
		ops = append(ops, &operations.VerifyImage{Image: observeImage})
	}
	releaseFilters := makeObserveReleaseFilters(args)
	imageIndex, err := updates.ImageReleaseIndex(kcdConfig, releaseFilters...)
	if err != nil {
		return nil, err
	}
	newImage := image.NewDockerImageRef(observeImage)
	imageTags := updates.BuildTagIndexFromNewImageRef(newImage, imageIndex)
	allUpdates := make([]updates.ImageUpdate, 0)
	for _, release := range imageIndex[newImage.WithoutTag()] {
		imageUpdates, err := updates.FindImageUpdatesForRelease(release, imageTags)
		if err != nil {
			return nil, err
		}
		allUpdates = append(allUpdates, imageUpdates...)
	}
	if len(allUpdates) == 0 {
		fmt.Printf("No matching release found for image %s.\n", observeImage)
		noOps := make([]operations.Operation, 0)
		return noOps, nil
	}
	verb := "May"
	if observePatch {
		verb = "Will"
	}
	for _, update := range allUpdates {
		fmt.Printf("%s update env %q release %q image %q tag %q -> %q\n", verb, update.Release.Environment.Name, update.Release.Name, update.ImageRepo, update.OldTag, update.NewTag)
		if observePatch {
			ops = append(ops, operations.NewPatchReleaseFile(update))
		}
	}
	return ops, nil
}

func init() {
	rootCmd.AddCommand(observeCmd)
	observeCmd.Flags().StringVarP(&observeImage, "image", "i", "", "a new image, including tag")
	observeCmd.Flags().StringSliceVarP(&observeReleases, "releases", "r", []string{}, "limit the update to or more specific releases")
	observeCmd.Flags().StringVarP(&observeCluster, "cluster", "c", "", "look for updates in all environments in CLUSTER")
	observeCmd.Flags().StringVar(&observeChart, "chart", "", "a new chart version, format: REPO/CHART:VERSION")
	observeCmd.Flags().BoolVar(&observePatch, "patch", false, "patch release files with updated tags")
	observeCmd.Flags().BoolVar(&observeVerify, "verify", false, "verify that image:tag exists")
}

func makeObserveReleaseFilters(args []string) []updates.ReleaseFilterFunc {
	filters := make([]updates.ReleaseFilterFunc, 0)
	if observeCluster != "" {
		filters = append(filters, updates.ClusterReleaseFilter(observeCluster))
	}
	if len(observeReleases) > 0 {
		filters = append(filters, updates.ReleaseFilter(observeReleases))
	}
	if len(args) == 1 {
		filters = append(filters, updates.EnvironmentReleaseFilter(args[0]))
	}
	if observeImage != "" {
		filters = append(filters, updates.ImageReleaseFilter(observeImage))
	}
	return filters
}

func patchReleasesFilesMaybe(imageUpdates []updates.ImageUpdate, patch bool) error {
	verb := "May"
	if patch {
		verb = "Will"
	}
	for _, update := range imageUpdates {
		fmt.Printf("%s update env %q release %q image %q tag %q -> %q\n", verb, update.Release.Environment.Name, update.Release.Name, update.ImageRepo, update.OldTag, update.NewTag)
	}
	if patch {
		for file, fileUpdates := range updates.GroupImageUpdatesByReleasesFile(imageUpdates) {
			fmt.Printf("Patching file: %s\n", file)
			if err := updates.PatchReleasesFiles(file, fileUpdates); err != nil {
				return err
			}
		}
	}
	return nil
}

var observeChartVersionRegex = regexp.MustCompile(`^([^/:]+/[^/:]+):(.+)$`)

func observeChartVersion(kcdConfig *model.KubeCDConfig) error {
	m := observeChartVersionRegex.FindStringSubmatch(observeChart)
	if len(m) != 4 {
		return fmt.Errorf(`--chart format must be "REPO/CHART:VERSION"", got %q`, observeChart)
	}
	return NotYetImplementedError("observe --chart")
}

func imageOrChart(image, chart *string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if (*chart != "" && *image != "") || (*chart == "" && *image == "") {
			return fmt.Errorf("specify exactly one of --image and --chart")
		}
		return nil
	}
}
