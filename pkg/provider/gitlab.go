package provider

import (
	"fmt"

	"github.com/zedge/kubecd/pkg/model"
)

type GitlabClusterProvider struct{ baseClusterProvider }

func (p *GitlabClusterProvider) GetClusterName() string {
	return "gitlab-deploy"
}

func (p *GitlabClusterProvider) GetUserName() string {
	return "gitlab-deploy"
}

func (p *GitlabClusterProvider) GetNamespace(env *model.Environment) string {
	return env.KubeNamespace
}

func (p *GitlabClusterProvider) GetClusterInitCommands() ([][]string, error) {
	return [][]string{}, nil
}

func (p *GitlabClusterProvider) LookupValueFrom(valueRef *model.ChartValueRef) (string, bool, error) {
	return "", false, fmt.Errorf("not implemented for Gitlab")
}
