package operations

import (
	"strings"

	"github.com/zedge/kubecd/pkg/helm"
	"github.com/zedge/kubecd/pkg/model"
	"github.com/zedge/kubecd/pkg/provider"
)

type EnvInit struct {
	*Base
	Environment *model.Environment
	GitlabMode  bool
}

type HelmInit struct {
	*Base
	Config *model.KubeCDConfig
}

func NewEnvInit(env *model.Environment, dryRun, gitlabMode bool) *EnvInit {
	return &EnvInit{
		Base:        newBase(dryRun),
		Environment: env,
		GitlabMode:  gitlabMode,
	}
}

func NewHelmInit(kcdConfig *model.KubeCDConfig, dryRun bool) *HelmInit {
	return &HelmInit{newBase(dryRun), kcdConfig}
}

func (o EnvInit) Prepare() error {
	cluster := o.Environment.GetCluster()
	cp, err := provider.GetClusterProvider(cluster, o.GitlabMode)
	if err != nil {
		return err
	}
	cmds, err := cp.GetClusterInitCommands()
	if err != nil {
		return err
	}
	for _, cmd := range cmds {
		o.Commands = append(o.Commands, Command{cmd, o.DryRun})
	}
	for _, cmd := range provider.GetContextInitCommands(cp, o.Environment) {
		o.Commands = append(o.Commands, Command{cmd, o.DryRun})
	}
	return nil
}

func (o EnvInit) String() string {
	var builder strings.Builder
	builder.WriteString("EnvInit(Environment=\"")
	builder.WriteString(o.Environment.Name)
	builder.WriteString("\") {\n")
	builder.WriteString(o.Base.String())
	builder.WriteString("}")
	return builder.String()
}

var _ Operation = &EnvInit{}

func (o HelmInit) Prepare() error {
	for _, cmd := range helm.RepoSetupCommands(o.Config.HelmRepos) {
		o.Commands = append(o.Commands, Command{cmd, o.DryRun})
	}
	return nil
}

func (o HelmInit) String() string {
	var builder strings.Builder
	builder.WriteString("HelmInit {\n")
	builder.WriteString(o.Base.String())
	builder.WriteString("}")
	return builder.String()
}
