package updates

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zedge/kubecd/pkg/helm"
	"github.com/zedge/kubecd/pkg/image"
	"github.com/zedge/kubecd/pkg/model"
)

type ImageUpdate struct {
	OldTag    string
	NewTag    string
	Release   *model.Release
	TagValue  string
	ImageRepo string
	Reason    string
}

type ChartUpdate struct {
	Release    *model.Release
	OldVersion string
	NewVersion string
	Reason     string
}

func FindImageUpdatesForRelease(release *model.Release, tagIndex TagIndex) ([]ImageUpdate, error) {
	updates := make([]ImageUpdate, 0)
	if release.Triggers == nil {
		return updates, nil
	}
	for _, trigger := range release.Triggers {
		if trigger.Image == nil || trigger.Image.Track == "" {
			fmt.Println("no trigger")
			continue
		}
		values, err := helm.GetResolvedValues(release)
		if err != nil {
			return nil, fmt.Errorf(`while looking for updates for release %q: %v`, release.Name, err)
		}
		imageRef := helm.GetImageRefFromImageTrigger(trigger.Image, values)
		if imageRef == nil {
			continue
		}
		imageTags := tagIndex.GetTags(imageRef)
		if imageTags == nil {
			//fmt.Printf("env:%s release:%s imageTags == nil, tagIndex=%#v, imageRef=%#v\n", release.Environment.Name, release.Name, tagIndex, *imageRef)
			continue
		}
		var currentTag image.TimestampedTag
		foundTag := false
		for _, tag := range imageTags {
			if imageRef.Tag == tag.Tag {
				currentTag = tag
				foundTag = true
			}
		}
		if !foundTag {
			fmt.Printf("did not find %s in %#v\n", imageRef.Tag, imageTags)
			continue
		}
		newestTag := image.GetNewestMatchingTag(currentTag, imageTags, trigger.Image.Track)
		if newestTag.Tag != currentTag.Tag {
			updates = append(updates, ImageUpdate{
				OldTag:    currentTag.Tag,
				NewTag:    newestTag.Tag,
				Release:   release,
				TagValue:  trigger.Image.TagValueString(),
				ImageRepo: imageRef.WithoutTag(),
				Reason:    "FIXME",
			})
		}
	}
	return updates, nil
}

type ReleaseFilterFunc func(*model.Release) bool

func ImageReleaseIndex(kcdConfig *model.KubeCDConfig, filters ...ReleaseFilterFunc) (map[string][]*model.Release, error) {
	result := make(map[string][]*model.Release)
releaseLoop:
	for _, release := range kcdConfig.AllReleases() {
		for _, filter := range filters {
			if !filter(release) {
				continue releaseLoop
			}
		}
		//fmt.Printf("evaluating release %q\n", release.Name)
		values, err := helm.GetResolvedValues(release)
		if err != nil {
			return nil, errors.Wrapf(err, "resolving values for env %q release %q", release.Environment.Name, release.Name)
		}
		//fmt.Printf("release %q triggers: %#v\n", release.Name, release.Triggers)
		for _, t := range release.Triggers {
			if t.Image == nil {
				//fmt.Printf("release %q has no trigger\n", release.Name)
				continue
			}
			repo := helm.GetImageRefFromImageTrigger(t.Image, values).WithoutTag()
			//fmt.Printf("release %q repo: %q\n", release.Name, repo)
			if _, found := result[repo]; !found {
				result[repo] = make([]*model.Release, 0)
			}
			result[repo] = append(result[repo], release)
		}
	}
	return result, nil
}
