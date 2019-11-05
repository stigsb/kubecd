package watch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChartNameAndVersion(t *testing.T) {
	cm := ChartMessage{Name: "evtail-server-v1.1.19.tgz"}
	n, v := cm.GetChartNameAndVersion()
	assert.Equal(t, "evtail-server", n)
	assert.Equal(t, "v1.1.19", v)
}
