package provider

import (
	"fmt"

	"github.com/zedge/kubecd/pkg/model"
)

type AksClusterProvider struct{ baseClusterProvider }

func (p *AksClusterProvider) GetClusterInitCommands() ([][]string, error) {
	panic("implement me")
}

func (p *AksClusterProvider) GetClusterName() string {
	panic("implement me")
}

func (p *AksClusterProvider) GetUserName() string {
	panic("implement me")
}

func (p *AksClusterProvider) GetNamespace(env *model.Environment) string {
	panic("implement me")
}

// LookupValueFrom returns a value, whether it was found and an error
func (p *AksClusterProvider) LookupValueFrom(valueRef *model.ChartValueRef) (string, bool, error) {
	return "", false, fmt.Errorf("not yet implemented for AKS")
}
