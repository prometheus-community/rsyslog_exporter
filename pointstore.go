// Copyright 2024 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"sort"
	"sync"
)

var (
	errPointNotFound = errors.New("point does not exist")
)

type pointStore struct {
	pointMap map[string]*point
	lock     *sync.RWMutex
}

func newPointStore() *pointStore {
	return &pointStore{
		pointMap: make(map[string]*point),
		lock:     &sync.RWMutex{},
	}
}

func (ps *pointStore) keys() []string {
	ps.lock.Lock()
	keys := make([]string, 0)
	for k := range ps.pointMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ps.lock.Unlock()
	return keys
}

func (ps *pointStore) set(p *point) error {
	var err error
	ps.lock.Lock()
	ps.pointMap[p.key()] = p
	ps.lock.Unlock()
	return err
}

func (ps *pointStore) get(name string) (*point, error) {
	ps.lock.Lock()
	if p, ok := ps.pointMap[name]; ok {
		ps.lock.Unlock()
		return p, nil
	}
	ps.lock.Unlock()
	return &point{}, errPointNotFound
}
