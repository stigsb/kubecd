package exec

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCachedRunner(t *testing.T) {
	cmd := "dd"
	args := []string{"if=/dev/urandom", "bs=32", "count=1"}
	t.Run("cached", func(t *testing.T) {
		cr := NewCachedRunner(10 * time.Minute)
		output1, err := cr.Run(cmd, args...)
		assert.NoError(t, err)
		output2, err := cr.Run(cmd, args...)
		assert.NoError(t, err)
		// The second Run() should get its output from the cache and produce idential results
		assert.Equal(t, output1, output2)
	})
	t.Run("expired", func(t *testing.T) {
		cr := NewCachedRunner(10 * time.Nanosecond)
		output1, err := cr.Run(cmd, args...)
		assert.NoError(t, err)
		time.Sleep(time.Microsecond) // longer than the ttl (10ns)
		output2, err := cr.Run(cmd, args...)
		assert.NoError(t, err)
		// The cache entry should have expired when then second Run() is called
		assert.NotEqual(t, output1, output2)
	})
}
