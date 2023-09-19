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
	"application/core/interfaces"
	"application/core/model"
	"application/driven/calendar"
	"time"

	"github.com/google/uuid"
	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logutils"
)

// appShared contains shared implementations
type appShared struct {
	app *Application
}

func (a appShared) getSurvey(id string, orgID string, appID string) (*model.Survey, error) {
	return a.app.storage.GetSurvey(id, orgID, appID)
}

func (a appShared) getSurveys(orgID string, appID string, creatorID *string, surveyIDs []string, surveyTypes []string, calendarEventID string, limit *int, offset *int) ([]model.Survey, error) {
	return a.app.storage.GetSurveys(orgID, appID, creatorID, surveyIDs, surveyTypes, calendarEventID, limit, offset)
}

func (a appShared) createSurvey(survey model.Survey, externalIDs map[string]string) (*model.Survey, error) {
	survey.ID = uuid.NewString()
	survey.DateCreated = time.Now().UTC()
	survey.DateUpdated = nil

	if len(survey.CalendarEventID) > 0 {
		// check if user is admin of calendar event

		// Get external ID field
		envConfig, err := a.app.GetEnvConfigs()
		if err != nil {
			return nil, errors.WrapErrorAction(logutils.ActionGet, model.TypeConfig, logutils.StringArgs(model.ConfigTypeEnv), err)
		}
		externalID := externalIDs[envConfig.ExternalID]

		user := calendar.User{AccountID: survey.CreatorID, ExternalID: externalID}
		eventUsers, err := a.app.calendar.GetEventUsers(survey.OrgID, survey.AppID, survey.CalendarEventID, []calendar.User{user}, nil, calendar.EventRoleAdmin, nil)
		if err != nil {
			return nil, err
		}
		for _, eventUser := range eventUsers {
			if (eventUser.User.ExternalID == externalID || eventUser.User.AccountID == survey.CreatorID) && eventUser.Role == calendar.EventRoleAdmin {
				return a.app.storage.CreateSurvey(survey)
			}
		}
		return nil, errors.Newf("account not an admin of calendar event")
		// if len(eventUsers) == 0 { // user is not admin
		// 	return nil, errors.Newf("")
		// }
	}

	return a.app.storage.CreateSurvey(survey)
}

func (a appShared) updateSurvey(survey model.Survey, userID string, externalID string, admin bool) error {
	admin, err := a.isAdminUser(survey, userID, externalID, admin)
	if err != nil {
		return errors.WrapErrorAction("checking", "event admin", &logutils.FieldArgs{"user_id": userID, "external_id": externalID}, err)
	}

	return a.app.storage.UpdateSurvey(survey, admin)
}

func (a appShared) deleteSurvey(id string, orgID string, appID string, userID string, externalID string, admin bool) error {
	transaction := func(storage interfaces.Storage) error {
		//1. find survey
		survey, err := storage.GetSurvey(id, orgID, appID)
		if err != nil {
			return errors.WrapErrorAction(logutils.ActionGet, model.TypeSurvey, nil, err)
		}
		if survey == nil {
			return errors.ErrorData(logutils.StatusMissing, model.TypeSurvey, &logutils.FieldArgs{"id": id, "app_id": appID, "org_id": orgID})
		}

		//2. check if user is event admin
		admin, err = a.isAdminUser(*survey, userID, externalID, admin)
		if err != nil {
			return errors.WrapErrorAction("checking", "event admin", &logutils.FieldArgs{"user_id": userID, "external_id": externalID}, err)
		}

		//3. delete survey
		err = storage.DeleteSurvey(survey.ID, survey.OrgID, survey.AppID, userID, admin)
		if err != nil {
			return errors.WrapErrorAction(logutils.ActionDelete, model.TypeSurvey, nil, err)
		}

		return nil
	}

	return a.app.storage.PerformTransaction(transaction)
}

func (a appShared) isAdminUser(survey model.Survey, userID string, externalID string, admin bool) (bool, error) {
	user := calendar.User{AccountID: userID, ExternalID: externalID}

	// if survey has an associated event and the current user is not already an admin user, check if user is an event admin
	if !admin && survey.CalendarEventID != "" {
		// the calendar BB event users API only uses the event ID for now
		eventUsers, err := a.app.calendar.GetEventUsers(survey.OrgID, survey.AppID, survey.CalendarEventID, []calendar.User{user}, nil, "", nil)
		if err != nil {
			return false, errors.WrapErrorAction(logutils.ActionGet, "event users", &logutils.FieldArgs{"event_id": survey.CalendarEventID}, err)
		}
		for _, eventUser := range eventUsers {
			// if the current user is an event admin, then act as an admin update
			if (eventUser.User.AccountID == userID || eventUser.User.ExternalID == externalID) && eventUser.Role == calendar.EventRoleAdmin {
				return true, nil
			}
		}
	}

	return admin, nil
}

// newAppShared creates new appShared
func newAppShared(app *Application) appShared {
	return appShared{app: app}
}
