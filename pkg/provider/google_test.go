/*
 * Copyright 2018-2020 Zedge, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kubecd/kubecd/pkg/exec"
	"github.com/kubecd/kubecd/pkg/model"
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
