package provider

import (
	"testing"

	"github.com/zedge/kubecd/pkg/exec"
)

// TestHelperProcess is required boilerplate (one per package) for using exec.TestRunner
func TestHelperProcess(t *testing.T) {
	exec.InsideHelperProcess()
}

