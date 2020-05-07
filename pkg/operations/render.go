package operations

import (
	"strings"

	"github.com/kubecd/kubecd/pkg/helm"
	"github.com/kubecd/kubecd/pkg/model"
)

type Render struct {
	*CommandBase
	Release *model.Release
}

func (o *Render) Prepare() error {
	if o.Release.Chart != nil {
		return o.prepareForHelmChart()
	}
	if o.Release.ResourceFiles != nil {
		return o.prepareForResourceFiles()
	}
	return nil
}

func (o *Render) prepareForHelmChart() error {
	templateCmds, err := helm.GenerateTemplateCommands(o.Release)
	if err != nil {
		return err
	}
	for _, cmd := range templateCmds {
		o.Commands = append(o.Commands, Command{Argv: cmd, DryRun: o.DryRun})
	}
	return nil
}

func (o *Render) prepareForResourceFiles() error {
	relFile := o.Release.FromFile
	for _, resourceFile := range o.Release.ResourceFiles {
		o.Commands = append(o.Commands, Command{[]string{"echo", "---"}, o.DryRun})
		o.Commands = append(o.Commands, Command{[]string{"echo", "#", "Source:", model.ResolvePathFromFile(resourceFile, relFile)}, o.DryRun})
		o.Commands = append(o.Commands, Command{[]string{"cat", model.ResolvePathFromFile(resourceFile, relFile)}, o.DryRun})
	}
	return nil
}

func (o Render) String() string {
	var builder strings.Builder
	builder.WriteString("Render(Environment=\"")
	builder.WriteString(o.Release.Environment.Name)
	builder.WriteString("\", Release=\"")
	builder.WriteString(o.Release.Name)
	builder.WriteString("\") {\n")
	builder.WriteString(o.CommandBase.String())
	builder.WriteString("}")
	return builder.String()
}

func NewRender(release *model.Release, dryRun bool) *Render {
	return &Render{
		CommandBase: NewCommandBase(dryRun),
		Release:     release,
	}
}

var _ Operation = &Render{}
