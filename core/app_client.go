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
	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logutils"
)

// appClient contains client implementations
type appClient struct {
	app *Application
}

// Surveys
// GetSurvey returns the survey with the provided ID
func (a appClient) GetSurvey(id string, orgID string, appID string) (*model.Survey, error) {
	return a.app.shared.getSurvey(id, orgID, appID)
}

// GetSurvey returns surveys matching the provided query
func (a appClient) GetSurveys(orgID string, appID string, surveyIDs []string, surveyTypes []string, limit *int, offset *int) ([]model.Survey, error) {
	return a.app.shared.getSurveys(orgID, appID, surveyIDs, surveyTypes, limit, offset)
}

// CreateSurvey creates a new survey
func (a appClient) CreateSurvey(survey model.Survey) (*model.Survey, error) {
	return a.app.shared.createSurvey(survey)
}

// UpdateSurvey updates the provided survey
func (a appClient) UpdateSurvey(survey model.Survey) error {
	return a.app.shared.updateSurvey(survey, false)
}

// DeleteSurvey deletes the survey with the specified ID
func (a appClient) DeleteSurvey(id string, orgID string, appID string, userID string) error {
	return a.app.shared.deleteSurvey(id, orgID, appID, &userID)
}

// Survey Response
// GetSurveyResponse returns the survey response with the provided ID
func (a appClient) GetSurveyResponse(id string, orgID string, appID string, userID string) (*model.SurveyResponse, error) {
	return a.app.storage.GetSurveyResponse(id, orgID, appID, userID)
}

// GetSurveyResponses returns the survey responses matching the provided filters
func (a appClient) GetSurveyResponses(orgID string, appID string, userID string, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time, limit *int, offset *int) ([]model.SurveyResponse, error) {
	return a.app.storage.GetSurveyResponses(orgID, appID, userID, surveyIDs, surveyTypes, startDate, endDate, limit, offset)
}

// CreateSurveyResponse creates a new survey response
func (a appClient) CreateSurveyResponse(surveyResponse model.SurveyResponse) (*model.SurveyResponse, error) {
	surveyResponse.ID = uuid.NewString()
	surveyResponse.DateCreated = time.Now().UTC()
	surveyResponse.DateUpdated = nil
	return a.app.storage.CreateSurveyResponse(surveyResponse)
}

// UpdateSurveyResponse updates the provided survey response
func (a appClient) UpdateSurveyResponse(surveyResponse model.SurveyResponse) error {
	return a.app.storage.UpdateSurveyResponse(surveyResponse)
}

// DeleteSurveyResponse deletes the survey with the specified ID
func (a appClient) DeleteSurveyResponse(id string, orgID string, appID string, userID string) error {
	return a.app.storage.DeleteSurveyResponse(id, orgID, appID, userID)
}

// DeleteSurveyResponses deletes the survey responses matching the provided filters
func (a appClient) DeleteSurveyResponses(orgID string, appID string, userID string, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time) error {
	return a.app.storage.DeleteSurveyResponses(orgID, appID, userID, surveyIDs, surveyTypes, startDate, endDate)
}

// Survey Alerts
// CreateSurveyAlert creates a new survey alert
func (a appClient) CreateSurveyAlert(surveyAlert model.SurveyAlert) error {
	contacts, err := a.app.storage.GetAlertContactsByKey(surveyAlert.ContactKey, surveyAlert.OrgID, surveyAlert.AppID)
	if err != nil {
		return err
	}

	for i := 0; i < len(contacts); i++ {
		if contacts[i].Type == "email" {
			subject, ok := surveyAlert.Content["subject"].(string)
			if !ok {
				return errors.ErrorData(logutils.StatusMissing, "subject", nil)
			}
			body, ok := surveyAlert.Content["body"].(string)
			if !ok {
				return errors.ErrorData(logutils.StatusMissing, "body", nil)
			}
			a.app.notifications.SendMail(contacts[i].Address, subject, body)
		}
	}

	return nil
}

// newAppClient creates new appClient
func newAppClient(app *Application) appClient {
	return appClient{app: app}
}
