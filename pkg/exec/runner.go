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

package exec

import (
	"context"
	"fmt"
	osexec "os/exec"
	"strings"
	"sync"
	"time"
)

type Runner interface {
	Run(string, ...string) ([]byte, error)
	RunContext(context.Context, string, ...string) ([]byte, error)
}

type RealRunner struct{}

func (r RealRunner) RunContext(ctx context.Context, cmd string, args ...string) ([]byte, error) {
	return osexec.CommandContext(ctx, cmd, args...).Output()
}

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
	return r.RunContext(context.Background(), cmd, args...)
}

func (r CachedRunner) RunContext(ctx context.Context, cmd string, args ...string) ([]byte, error) {
	r.cacheMu.Lock()
	defer r.cacheMu.Unlock()
	key := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	entry, found := r.cache[key]
	now := time.Now()
	if !found || now.After(entry.insertTime.Add(r.ttl)) {
		data, err := osexec.CommandContext(ctx, cmd, args...).Output()
		if err != nil {
			return nil, err
		}
		entry = cachedEntry{data: data, insertTime: now}
		r.cache[key] = entry
	}
	return entry.data, nil
}
