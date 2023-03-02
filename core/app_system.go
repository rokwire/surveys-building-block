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
	"time"

	"github.com/google/uuid"
)

// appSystem contains system implementations
type appSystem struct {
	app *Application
}

// GetConfig gets the configs for the provided id
func (a appSystem) GetConfig(id string) (*model.Config, error) {
	return a.app.storage.FindConfigByID(id)
}

// GetConfigs gets the configs for the provided type
func (a appSystem) GetConfigs(configType *string) ([]model.Config, error) {
	return a.app.storage.FindConfigs(configType)
}

// CreateConfig creates the provided config
func (a appSystem) CreateConfig(config model.Config) error {
	config.ID = uuid.NewString()
	config.DateCreated = time.Now()
	return a.app.storage.InsertConfig(config)
}

// DeleteConfig deletes the configs for the provided id
func (a appSystem) DeleteConfig(id string) error {
	return a.app.storage.DeleteConfig(id)
}

// newAppSystem creates new appSystem
func newAppSystem(app *Application) appSystem {
	return appSystem{app: app}
}
