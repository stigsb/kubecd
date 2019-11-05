package watch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"

	"github.com/zedge/kubecd/pkg/image"
)

// {
//   "action":"INSERT",
//   "digest":"us.gcr.io/zedge-test/clean-event-property@sha256:ad4e8ab4a104804382af8f1c56de541f9c8a056e9e70187e6bd0fd1c98273495",
//   "tag":"us.gcr.io/zedge-test/clean-event-property:996bc29"
// }

type GCRMessage struct {
	Action string `json:"action"`
	Digest string `json:"digest"`
	Tag    string `json:"tag"`
}

type GCRWatcher struct {
	Project      string
	Subscription *pubsub.Subscription
	Client       *pubsub.Client
	AckDeadline  time.Duration
	Context      context.Context
}

func NewGCRWatcher(ctx context.Context, project, subscription string) (*GCRWatcher, error) {
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %w", err)
	}
	watcher := &GCRWatcher{
		Project:      project,
		Subscription: client.Subscription(subscription),
		Client:       client,
		AckDeadline:  5 * time.Minute,
		Context:      ctx,
	}
	watcher.Subscription.ReceiveSettings.Synchronous = true
	watcher.Subscription.ReceiveSettings.MaxOutstandingMessages = 1
	return watcher, nil
}

type ImageWatchCallback func(*image.DockerImageRef) error

func (w GCRWatcher) Run(cb ImageWatchCallback) error {
	if _, err := w.Subscription.Update(w.Context, pubsub.SubscriptionConfigToUpdate{AckDeadline: w.AckDeadline}); err != nil {
		return fmt.Errorf("sub.Update: %w", err)
	}
	return w.Subscription.Receive(w.Context, func(ctx context.Context, msg *pubsub.Message) {
		var imageUpdate GCRMessage
		log.Printf("Got GCR message: %v\n", string(msg.Data))
		if err := json.Unmarshal(msg.Data, &imageUpdate); err != nil {
			log.Printf("Could not decode GCR message: %v: %v\n", string(msg.Data), err)
			return
		}
		img := image.NewDockerImageRef(imageUpdate.Tag)
		if img.Tag == "" || img.Tag == "latest" {
			msg.Ack()
			return
		}
		if err := cb(img); err != nil {
			log.Printf("ImageWatchCallback failed: %v\n", err)
			return
		}
		//msg.Ack()
	})
}

//func (w GCRWatcher) Receive(ctx context.Context, msg *pubsub.Message) {
//	var gcrMessage GCRMessage
//	if err := json.Unmarshal(msg.Data, &gcrMessage); err != nil {
//		log.Printf("Could not decode GCR message: %v: %v\n", string(msg.Data), err)
//		return
//	}
//	log.Printf("GCRMessage: %#v\n", gcrMessage)
//	//msg.Ack()
//}
