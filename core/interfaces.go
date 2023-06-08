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
)

// Shared exposes shared APIs for other interface implementations
type Shared interface {
	// Surveys
	getSurvey(id string, orgID string, appID string) (*model.Survey, error)
	getSurveys(orgID string, appID string, surveyIDs []string, surveyTypes []string, limit *int, offset *int, groupID string) ([]model.Survey, error)
	getAllSurveyResponses(id string, orgID string, appID string, userToken string, userID string, groupID string, startDate *time.Time, endDate *time.Time, limit *int, offset *int) ([]model.SurveyResponse, error)
	createSurvey(survey model.Survey, user model.User) (*model.Survey, error)
	updateSurvey(survey model.Survey, userID string) error
	deleteSurvey(id string, orgID string, appID string, userID *string) error
}
