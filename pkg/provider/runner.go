package provider

import (
	"time"

	"github.com/kubecd/kubecd/pkg/exec"
)

var cachedRunner exec.Runner = exec.NewCachedRunner(10 * time.Minute)

//var runner exec.Runner = exec.RealRunner{}