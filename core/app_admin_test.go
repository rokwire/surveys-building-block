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

package core_test

import (
	"application/core"
	"application/core/mocks"
	"application/core/model"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

func Test_appAdmin_GetExample(t *testing.T) {
	app := buildTestAppGetExample(t)
	test_Shared_GetExample(t, app.Admin)
}

func Test_appAdmin_CreateExample(t *testing.T) {
	example := buildTestExample()
	example.ID = ""

	storageSuccess := mocks.NewStorage(t)
	storageSuccess.On("InsertExample", mock.AnythingOfType("model.Example")).Return(nil)
	appSuccess := buildTestApplication(storageSuccess)

	storageError := mocks.NewStorage(t)
	storageError.On("InsertExample", mock.AnythingOfType("model.Example")).Return(errors.New("An error occurred connecting to the database"))
	appError := buildTestApplication(storageError)

	type args struct {
		example model.Example
	}
	tests := []struct {
		name    string
		a       core.Admin
		args    args
		want    *model.Example
		wantErr bool
	}{
		{"success", appSuccess.Admin, args{example}, &example, false},
		{"error", appError.Admin, args{example}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.CreateExample(tt.args.example)
			if (err != nil) != tt.wantErr {
				t.Errorf("appAdmin.CreateExample() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if _, err := uuid.Parse(got.ID); err != nil {
					t.Errorf("appAdmin.CreateExample() invalid uuid (%s): %v", got.ID, err)
					return
				}
				tt.want.ID = got.ID
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appAdmin.CreateExample() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appAdmin_UpdateExample(t *testing.T) {
	example := buildTestExample()
	exampleWrongOrg := buildTestExample()
	exampleWrongOrg.OrgID = "org2"
	exampleWrongID := buildTestExample()
	exampleWrongID.ID = "id2"

	storage := mocks.NewStorage(t)
	storage.On("UpdateExample", example).Return(nil)
	storage.On("UpdateExample", exampleWrongOrg).Return(errors.New("no example found"))
	storage.On("UpdateExample", exampleWrongID).Return(errors.New("no example found"))
	app := buildTestApplication(storage)

	type args struct {
		example model.Example
	}
	tests := []struct {
		name    string
		a       core.Admin
		args    args
		wantErr bool
	}{
		{"found", app.Admin, args{example}, false},
		{"invalid org", app.Admin, args{exampleWrongOrg}, true},
		{"invalid id", app.Admin, args{exampleWrongID}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.UpdateExample(tt.args.example); (err != nil) != tt.wantErr {
				t.Errorf("appAdmin.UpdateExample() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_appAdmin_DeleteExample(t *testing.T) {
	example := buildTestExample()
	storage := mocks.NewStorage(t)
	storage.On("DeleteExample", example.OrgID, example.AppID, example.ID).Return(nil)
	storage.On("DeleteExample", "org2", example.AppID, example.ID).Return(errors.New("no example found"))
	storage.On("DeleteExample", example.OrgID, example.AppID, "id2").Return(errors.New("no example found"))
	app := buildTestApplication(storage)

	type args struct {
		orgID string
		appID string
		id    string
	}
	tests := []struct {
		name    string
		a       core.Admin
		args    args
		wantErr bool
	}{
		{"found", app.Admin, args{example.OrgID, example.AppID, example.ID}, false},
		{"invalid org", app.Admin, args{"org2", example.AppID, example.ID}, true},
		{"invalid id", app.Admin, args{example.OrgID, example.AppID, "id2"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.DeleteExample(tt.args.orgID, tt.args.appID, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("appAdmin.DeleteExample() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
