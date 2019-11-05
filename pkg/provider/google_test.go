package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zedge/kubecd/pkg/exec"
	"github.com/zedge/kubecd/pkg/model"
)

const testIpAddress = "1.2.3.4"

var _ ClusterProvider = &GkeClusterProvider{}

func TestResolveGceAddressValue(t *testing.T) {
	oldRunner := cachedRunner
	defer func() { cachedRunner = oldRunner }()
	cachedRunner = exec.TestRunner{Output: []byte(testIpAddress)}
	zone := "us-central1-a"
	cluster := model.Cluster{
		Name: "kcd-clustername",
		Provider: model.Provider{
			GKE: &model.GkeProvider{
				Project:     "test-project",
				Zone:        &zone,
				ClusterName: "gke-clustername",
			},
		},
	}
	address := &model.GceAddressValueRef{
		Name:     "my-address",
		IsGlobal: false,
	}
	p := GkeClusterProvider{baseClusterProvider{Cluster: &cluster}}
	valueRef := &model.ChartValueRef{GceResource: &model.GceValueRef{Address: address}}
	val, found, err := p.LookupValueFrom(valueRef)
	assert.True(t, found)
	assert.NoError(t, err)
	assert.Equal(t, testIpAddress, val)
}
