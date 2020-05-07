package operations

import (
	"fmt"
	"github.com/kubecd/kubecd/pkg/image"
	"github.com/kubecd/kubecd/pkg/updates"
	"strings"
)

type VerifyImage struct {
	Image string
}

func (o VerifyImage) Execute() error {
	panic("implement me")
}

func (o VerifyImage) String() string {
	var builder strings.Builder
	builder.WriteString(`VerifyImage(Image="`)
	builder.WriteString(o.Image)
	builder.WriteString(`")`)
	return builder.String()
}

func (o VerifyImage) Prepare() error {
	imageRef := image.NewDockerImageRef(o.Image)
	existingTags, err := image.GetTagsForDockerImage(o.Image)
	if err != nil {
		return err
	}
	for _, tsTag := range existingTags {
		if tsTag.Tag == imageRef.Tag {
			return nil
		}
	}
	return fmt.Errorf(`tag %q not found for imageRepo %q`, imageRef.Tag, imageRef.WithoutTag())
}

var _ Operation = &VerifyImage{}

func NewPatchReleaseFile(update updates.ImageUpdate) *PatchReleaseFile {
	return &PatchReleaseFile{
		CommandBase: NewCommandBase(false),
		ImageUpdate: update,
	}
}

type PatchReleaseFile struct {
	*CommandBase
	ImageUpdate updates.ImageUpdate
}

func (o PatchReleaseFile) Prepare() error {
	return nil
}

var _ Operation = &PatchReleaseFile{}
