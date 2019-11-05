package provider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zedge/kubecd/pkg/model"
)

var zoneToRegionRegexp = regexp.MustCompile(`-[a-z]$`)

type GkeClusterProvider struct { baseClusterProvider }

func (p *GkeClusterProvider) GetClusterInitCommands() ([][]string, error) {
	gcloudCommand := []string{
		"gcloud", "container", "clusters", "get-credentials", "--project", p.Provider.GKE.Project,
	}
	if p.Provider.GKE.Zone != nil {
		gcloudCommand = append(gcloudCommand, "--zone", *p.Provider.GKE.Zone)
	} else {
		gcloudCommand = append(gcloudCommand, "--region", *p.Provider.GKE.Region)
	}
	gcloudCommand = append(gcloudCommand, p.Provider.GKE.ClusterName)
	return [][]string{gcloudCommand}, nil
}

func (p *GkeClusterProvider) GetClusterName() string {
	// 'gke_{gke.project}_{zone_or_region}_{gke.clusterName}'
	return fmt.Sprintf("gke_%s_%s_%s", p.Provider.GKE.Project, regionOrZone(p.Provider.GKE), p.Provider.GKE.ClusterName)
}

func regionOrZone(gke *model.GkeProvider) string {
	if gke.Region != nil {
		return *gke.Region
	}
	return *gke.Zone
}

func (p *GkeClusterProvider) GetUserName() string {
	gke := p.Provider.GKE
	return fmt.Sprintf("gke_%s_%s_%s", gke.Project, regionOrZone(gke), gke.ClusterName)
}

func (p *GkeClusterProvider) GetNamespace(environment *model.Environment) string {
	return environment.KubeNamespace
}

// LookupValueFrom returns a value, whether it was found and an error
func (p *GkeClusterProvider) LookupValueFrom(valueRef *model.ChartValueRef) (string, bool, error) {
	if gceRes := valueRef.GceResource; gceRes != nil {
		if gceRes.Address != nil {
			addr, err := p.lookupAddress(valueRef.GceResource.Address)
			return addr, true, err
		}
	}
	return "", false, nil
}

func (p *GkeClusterProvider) lookupAddress(ref *model.GceAddressValueRef) (string, error) {
	gke := p.Provider.GKE
	argv := []string{"compute", "addresses", "describe", ref.Name, "--format", "value(address)", "--project", gke.Project}
	if ref.IsGlobal {
		argv = append(argv, "--global")
	} else {
		argv = append(argv, "--region")
		if gke.Zone != nil {
			argv = append(argv, zoneToRegionRegexp.ReplaceAllString(*gke.Zone, ""))
		} else {
			argv = append(argv, *gke.Region)
		}
	}
	out, err := cachedRunner.Run("gcloud", argv...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
