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
	"application/driven/storage"
)

// Default exposes client APIs for the driver adapters
type Default interface {
	GetVersion() string
}

// Client exposes client APIs for the driver adapters
type Client interface {
	GetExample(orgID string, appID string, id string) (*model.Example, error)
}

// Admin exposes administrative APIs for the driver adapters
type Admin interface {
	GetExample(orgID string, appID string, id string) (*model.Example, error)
	CreateExample(example model.Example) (*model.Example, error)
	UpdateExample(example model.Example) error
	DeleteExample(orgID string, appID string, id string) error
}

// BBs exposes Building Block APIs for the driver adapters
type BBs interface {
	GetExample(orgID string, appID string, id string) (*model.Example, error)
}

// TPS exposes third-party service APIs for the driver adapters
type TPS interface {
	GetExample(orgID string, appID string, id string) (*model.Example, error)
}

// System exposes system administrative APIs for the driver adapters
type System interface {
	GetConfig(id string) (*model.Config, error)
	SaveConfig(configs model.Config) error
	DeleteConfig(id string) error
}

// Shared exposes shared APIs for other interface implementations
type Shared interface {
	getExample(orgID string, appID string, id string) (*model.Example, error)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	RegisterStorageListener(storageListener storage.Listener)
	PerformTransaction(func(adapter storage.Adapter) error) error

	GetConfig(id string) (*model.Config, error)
	SaveConfig(configs model.Config) error
	DeleteConfig(id string) error

	GetExample(orgID string, appID string, id string) (*model.Example, error)
	InsertExample(example model.Example) error
	UpdateExample(example model.Example) error
	DeleteExample(orgID string, appID string, id string) error
}
