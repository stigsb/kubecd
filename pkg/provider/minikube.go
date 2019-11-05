package provider

import (
	"fmt"

	"github.com/zedge/kubecd/pkg/model"
)

type MinikubeClusterProvider struct{ baseClusterProvider }

func (p *MinikubeClusterProvider) GetClusterInitCommands() ([][]string, error) {
	return [][]string{}, nil
}

func (p *MinikubeClusterProvider) GetClusterName() string {
	return "minikube"
}

func (p *MinikubeClusterProvider) GetUserName() string {
	return "minikube"
}

func (p *MinikubeClusterProvider) GetNamespace(env *model.Environment) string {
	return env.KubeNamespace
}

func (p *MinikubeClusterProvider) LookupValueFrom(valueRef *model.ChartValueRef) (string, bool, error) {
	return "", false, fmt.Errorf("not implemented for Minikube")
}
