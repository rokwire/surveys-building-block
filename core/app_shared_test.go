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
	"time"
)

type getExample interface {
	GetExample(string, string, string) (*model.Example, error)
}

func buildTestExample() model.Example {
	return model.Example{ID: "id1", OrgID: "org1", AppID: "app1", Data: "Example data", DateCreated: time.Now(), DateUpdated: nil}
}

func buildTestAppGetExample(t *testing.T) *core.Application {
	example := buildTestExample()
	storage := mocks.NewStorage(t)
	storage.On("GetExample", example.OrgID, example.AppID, example.ID).Return(&example, nil)
	storage.On("GetExample", "org2", example.AppID, example.ID).Return(nil, errors.New("no example found"))
	storage.On("GetExample", example.OrgID, example.AppID, "id2").Return(nil, errors.New("no example found"))
	return buildTestApplication(storage)
}

func test_Shared_GetExample(t *testing.T, impl getExample) {
	example := buildTestExample()

	type args struct {
		orgID string
		appID string
		id    string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Example
		wantErr bool
	}{
		{"found", args{example.OrgID, example.AppID, example.ID}, &example, false},
		{"invalid org", args{"org2", example.AppID, example.ID}, nil, true},
		{"invalid id", args{example.OrgID, example.AppID, "id2"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := impl.GetExample(tt.args.orgID, tt.args.appID, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("appTPS.GetExample() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appTPS.GetExample() = %v, want %v", got, tt.want)
			}
		})
	}
}
