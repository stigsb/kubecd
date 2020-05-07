package operations

import (
	"strings"

	"github.com/kubecd/kubecd/pkg/helm"
	"github.com/kubecd/kubecd/pkg/model"
)

type Apply struct {
	*CommandBase
	Release *model.Release
	Debug   bool
}

func (o *Apply) Prepare() error {
	if o.Release.Chart != nil {
		return o.prepareForHelmChart()
	}
	if o.Release.ResourceFiles != nil {
		return o.prepareForResourceFiles()
	}
	return nil
}

func (o Apply) String() string {
	var builder strings.Builder
	builder.WriteString("Apply(Environment=\"")
	builder.WriteString(o.Release.Environment.Name)
	builder.WriteString("\",Release=\"")
	builder.WriteString(o.Release.Name)
	builder.WriteString("\") {\n")
	builder.WriteString(o.CommandBase.String())
	builder.WriteString("}")
	return builder.String()
}

func (o *Apply) prepareForHelmChart() error {
	cmd, err := helm.GenerateHelmApplyArgv(o.Release, o.DryRun, o.Debug)
	if err != nil {
		return err
	}
	// the generated helm command will have --dry-run, so we can run it in dry-run mode
	o.Commands = append(o.Commands, Command{cmd, false})
	return nil
}

func (o *Apply) prepareForResourceFiles() error {
	relFile := o.Release.FromFile
	absFiles := make([]string, len(o.Release.ResourceFiles))
	for i, path := range o.Release.ResourceFiles {
		absFiles[i] = model.ResolvePathFromFile(path, relFile)
	}
	cmd := []string{"kubectl", "--context", model.KubeContextName(o.Release.Environment.Name), "apply"}
	if o.DryRun {
		cmd = append(cmd, "--dry-run")
	}
	for _, file := range o.Release.ResourceFiles {
		cmd = append(cmd, "-f", file)
	}
	dryRun := false // the generated kubectl command will have --dry-run, so we can run it in dry-run mode
	o.Commands = append(o.Commands, Command{cmd, dryRun})
	return nil
}

func NewApply(release *model.Release, dryRun, debug bool) *Apply {
	return &Apply{
		CommandBase: NewCommandBase(dryRun),
		Release:     release,
		Debug:       debug,
	}
}

var _ Operation = &Apply{}
