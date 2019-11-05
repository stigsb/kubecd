package exec

import (
	"fmt"
	osexec "os/exec"
	"strings"
	"sync"
	"time"
)

type Runner interface {
	Run(string, ...string) ([]byte, error)
}

type RealRunner struct{}

func (r RealRunner) Run(cmd string, args ...string) ([]byte, error) {
	return osexec.Command(cmd, args...).Output()
}

func NewCachedRunner(ttl time.Duration) *CachedRunner {
	return &CachedRunner{
		cache:   make(map[string]cachedEntry),
		ttl:     ttl,
	}
}

type cachedEntry struct {
	data       []byte
	insertTime time.Time
}

type CachedRunner struct {
	cache   map[string]cachedEntry
	cacheMu sync.Mutex
	ttl     time.Duration
}

func (r CachedRunner) Run(cmd string, args ...string) ([]byte, error) {
	r.cacheMu.Lock()
	defer r.cacheMu.Unlock()
	key := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	entry, found := r.cache[key]
	now := time.Now()
	if !found || now.After(entry.insertTime.Add(r.ttl)) {
		data, err := osexec.Command(cmd, args...).Output()
		if err != nil {
			return nil, err
		}
		entry = cachedEntry{data: data, insertTime: now}
		r.cache[key] = entry
	}
	return entry.data, nil
}
