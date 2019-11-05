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
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/zedge/kubecd/pkg/image"
	"github.com/zedge/kubecd/pkg/model"
	"github.com/zedge/kubecd/pkg/updates"
	"github.com/zedge/kubecd/pkg/watch"
)

var (
	watchGCPProject      string
	watchGCRSubscription string
	//var watchChartsSubscription string
	watchAckDeadline        time.Duration
	defaultWatchAckDeadline = 5 * time.Minute
	watchCluster            string
	watchReleases           []string
)

/*
	ctx := context.Background()
	projectID := "zedge-test"
	subID := "gcr-kubecd"
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	sub := client.SubscriptionName(subID)
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		fmt.Printf("Got message: %q\n", string(msg.Data))
	})
	if err != nil {
		panic(err)
	}

*/
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for changes in pubsub or webhooks",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		kcdConfig, err := model.NewConfigFromFile(environmentsFile)
		if err != nil {
			return err
		}
		releaseFilters := makeImageWatchReleaseFilters(args)
		watcher, err := watch.NewGCRWatcher(context.Background(), watchGCPProject, watchGCRSubscription)
		if err != nil {
			return fmt.Errorf("watch.NewGCRWatcher: %w", err)
		}
		watcher.AckDeadline = watchAckDeadline
		return watcher.Run(func(img *image.DockerImageRef) error {
			fmt.Printf("kcd observe [env] --image=%s\n", img.WithTag())
			return watchObserveImageTag(kcdConfig, releaseFilters, img)
		})
	},
}

func watchObserveImageTag(kcdConfig *model.KubeCDConfig, releaseFilters []updates.ReleaseFilterFunc, newImage *image.DockerImageRef) error {
	imageIndex, err := updates.ImageReleaseIndex(kcdConfig, releaseFilters...)
	if err != nil {
		return err
	}
	imageTags := updates.BuildTagIndexFromNewImageRef(newImage, imageIndex)
	allUpdates := make([]updates.ImageUpdate, 0)
	for _, release := range imageIndex[newImage.WithoutTag()] {
		imageUpdates, err := updates.FindImageUpdatesForRelease(release, imageTags)
		if err != nil {
			return err
		}
		allUpdates = append(allUpdates, imageUpdates...)
	}
	if len(allUpdates) == 0 {
		fmt.Printf("No matching release found for image %s.\n", newImage.WithTag())
		return nil
	}
	if err = patchReleasesFilesMaybe(allUpdates, true); err != nil {
		return err
	}

	return nil
}

func makeImageWatchReleaseFilters(args []string) []updates.ReleaseFilterFunc {
	filters := make([]updates.ReleaseFilterFunc, 0)
	if len(watchReleases) > 0 {
		filters = append(filters, updates.ReleaseFilter(observeReleases))
	}
	if len(args) == 1 {
		filters = append(filters, updates.EnvironmentReleaseFilter(args[0]))
	}
	if watchCluster != "" {
		filters = append(filters, updates.ClusterReleaseFilter(watchCluster))
	}
	return filters
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().DurationVar(&watchAckDeadline, "ack-deadline", defaultWatchAckDeadline, "Ack deadline for handling pubsub messages")
	watchCmd.Flags().StringVar(&watchGCPProject, "gcp-project", "", "Google Cloud Project (for pubsub watches)")
	watchCmd.Flags().StringVar(&watchGCRSubscription, "gcr-subscription", "", "Google Cloud PubSub subscription to watch for image updates from GCR")
	watchCmd.Flags().StringVarP(&watchCluster, "cluster", "c", "", "watch for all environments in CLUSTER")
	watchCmd.Flags().StringSliceVarP(&watchReleases, "releases", "r", []string{}, "limit updates to or more specific releases")
	//watchCmd.Flags().StringVar(&watchChartsSubscription, "charts-subscription", "", "Google Cloud PubSub subscription to watch for chart updates from chartmusem+GCS")
}
