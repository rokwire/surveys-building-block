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

package core

import (
	"application/core/model"
)

// appSystem contains system implementations
type appSystem struct {
	app *Application
}

// GetConfig gets the configs for the provided id
func (a appSystem) GetConfig(id string) (*model.Config, error) {
	return a.app.storage.GetConfig(id)
}

// SaveConfig saves the provided configs
func (a appSystem) SaveConfig(configs model.Config) error {
	return a.app.storage.SaveConfig(configs)
}

// DeleteConfig deletes the configs for the provided id
func (a appSystem) DeleteConfig(id string) error {
	return a.app.storage.DeleteConfig(id)
}

// newAppSystem creates new appSystem
func newAppSystem(app *Application) appSystem {
	return appSystem{app: app}
}
