// Copyright 2022 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

type storageListener struct {
	adapter *Adapter
	DefaultListenerImpl
}

func (s *storageListener) OnConfigsUpdated() {
	s.adapter.cacheConfigs()
}

// Listener represents storage listener
type Listener interface {
	OnConfigsUpdated()
	OnExamplesUpdated()
}

// DefaultListenerImpl default listener implementation
type DefaultListenerImpl struct{}

// OnConfigsUpdated notifies that the configs collection has been updated
func (d *DefaultListenerImpl) OnConfigsUpdated() {}

// OnExamplesUpdated notifies that the examples collection has been updated
func (d *DefaultListenerImpl) OnExamplesUpdated() {}
