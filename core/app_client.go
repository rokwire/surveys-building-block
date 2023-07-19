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
	"application/driven/calendar"
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
func (a appClient) GetSurveys(orgID string, appID string, surveyIDs []string, surveyTypes []string, calendarEventID string, limit *int, offset *int) ([]model.Survey, error) {
	return a.app.shared.getSurveys(orgID, appID, surveyIDs, surveyTypes, calendarEventID, limit, offset)
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

// GetUserSurveyResponses returns the survey responses matching the provided filters for a specific user
func (a appClient) GetUserSurveyResponses(orgID string, appID string, userID string, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time, limit *int, offset *int) ([]model.SurveyResponse, error) {
	return a.app.storage.GetSurveyResponses(&orgID, &appID, &userID, surveyIDs, surveyTypes, startDate, endDate, limit, offset)
}

// GetAllSurveyResponses returns the survey responses matching the provided filters
func (a appClient) GetAllSurveyResponses(orgID string, appID string, userID string, surveyID string, startDate *time.Time, endDate *time.Time, limit *int, offset *int) ([]model.SurveyResponse, error) {
	var allResponses []model.SurveyResponse
	var err error

	survey, err := a.app.shared.getSurvey(surveyID, orgID, appID)
	if err != nil {
		return nil, err
	}

	// Check if user is admin of calendar event if survey associated with one
	if len(survey.CalendarEventID) > 0 {
		user := calendar.User { AccountID: userID }// TODO: Add networkID
		eventUsers, err := a.app.calendar.GetEventUsers(survey.OrgID, survey.AppID, survey.CalendarEventID, []calendar.User{user}, nil, "admin", nil)
		if err != nil {
			return nil, err
		}
		if len(eventUsers) == 0 {// user is not admin
			return nil, errors.Newf("account is not admin of calendar event")
		}
	}
	
	allResponses, err = a.app.storage.GetSurveyResponses(&orgID, &appID, nil, []string {surveyID}, nil, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}

	// If survey is anonymous strip userIDs
	if survey.Anonymous {
		for i := range allResponses {
			allResponses[i].UserID = ""
		}
	}

	return allResponses, nil
}

// CreateSurveyResponse creates a new survey response
func (a appClient) CreateSurveyResponse(surveyResponse model.SurveyResponse) (*model.SurveyResponse, error) {
	surveyResponse.ID = uuid.NewString()
	surveyResponse.DateCreated = time.Now().UTC()
	surveyResponse.DateUpdated = nil

	if len(surveyResponse.Survey.CalendarEventID) > 0 {
		// check if user attended calendar event
		user := calendar.User { AccountID: surveyResponse.UserID }// TODO: Add networkID
		attended := true
		eventUsers, err := a.app.calendar.GetEventUsers(surveyResponse.OrgID, surveyResponse.AppID, surveyResponse.Survey.CalendarEventID, []calendar.User{user}, nil, "", &attended)
		if err != nil {
			return nil, err
		}
		if len(eventUsers) == 0 {// user has not attended
			return nil, errors.Newf("account is not attended calendar event")
		}
	}

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
