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

import (
	"application/core/model"

	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logutils"
	"go.mongodb.org/mongo-driver/bson"
)

// GetExample gets example by id
func (a Adapter) GetExample(orgID string, appID string, id string) (*model.Example, error) {
	filter := bson.M{"org_id": orgID, "app_id": appID, "_id": id}

	var data *model.Example
	err := a.db.examples.FindOneWithContext(a.context, filter, &data, nil)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, model.TypeExample, nil, err)
	}

	return data, nil
}

// InsertExample inserts a new example
func (a Adapter) InsertExample(example model.Example) error {
	_, err := a.db.examples.InsertOneWithContext(a.context, example)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionInsert, model.TypeExample, nil, err)
	}

	return nil
}

// UpdateExample updates an example
func (a Adapter) UpdateExample(example model.Example) error {
	filter := bson.M{"org_id": example.OrgID, "app_id": example.AppID, "_id": example.ID}
	update := bson.M{"$set": bson.M{"data": example.Data}}

	_, err := a.db.examples.UpdateOneWithContext(a.context, filter, update, nil)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionUpdate, model.TypeExample, nil, err)
	}
	return nil
}

// DeleteExample deletes an example
func (a Adapter) DeleteExample(orgID string, appID string, id string) error {
	filter := bson.M{"org_id": orgID, "app_id": appID, "_id": id}

	res, err := a.db.examples.DeleteOneWithContext(a.context, filter, nil)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionDelete, model.TypeExample, nil, err)
	}
	if res.DeletedCount != 1 {
		return errors.ErrorData(logutils.StatusMissing, model.TypeConfig, logutils.StringArgs(id))
	}

	return nil
}
