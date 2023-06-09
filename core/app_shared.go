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
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rokwire/logging-library-go/v2/logutils"
	"github.com/rokwire/logging-library-go/v2/errors"
)

// appShared contains shared implementations
type appShared struct {
	app *Application
}

func (a appShared) getSurvey(id string, orgID string, appID string) (*model.Survey, error) {
	return a.app.storage.GetSurvey(id, orgID, appID)
}

func (a appShared) getSurveys(orgID string, appID string, surveyIDs []string, surveyTypes []string, limit *int, offset *int, groupID string) ([]model.Survey, error) {
	return a.app.storage.GetSurveys(orgID, appID, surveyIDs, surveyTypes, limit, offset, groupID)
}

func (a appShared) getAllSurveyResponses(id string, orgID string, appID string, userToken string, userID string, groupIDs []string, startDate *time.Time, endDate *time.Time, limit *int, offset *int) ([]model.SurveyResponse, error) {
	var allResponses []model.SurveyResponse

	for _, groupID := range groupIDs {
		group, err := a.app.groups.GetGroupDetails(userToken, groupID)
		if err != nil {
			return nil, err
		}

		survey, err := a.app.storage.GetSurvey(id, orgID, appID)
		if err != nil {
			return nil, err
		}

		if (group.IsCurrentUserAdmin(userID) || group.CreatorID == userID) && !survey.Sensitive {
			responses, err := a.app.storage.GetAllSurveyResponses(&orgID, &appID, &id, startDate, endDate, limit, offset)
			if err != nil {
				return nil, errors.WrapErrorAction(logutils.ActionGet, logutils.TypePermission, nil, fmt.Errorf("cannot get responses"))
			}
			for _, response := range responses {
				allResponses = append(allResponses, response)
			}
		}
	}

	return allResponses, nil
}

func (a appShared) createSurvey(survey model.Survey, userName string, token string) (*model.Survey, error) {
	survey.ID = uuid.NewString()
	survey.DateCreated = time.Now().UTC()
	survey.DateUpdated = nil

	if len(survey.UserIDs) > 0 {
		a.sendNotificationsToUserList(survey, survey.CreatorID, userName)
	} else if len(survey.GroupIDs) > 0 {
		for _, groupID := range survey.GroupIDs {
			a.sendNotificationsToGroup(survey, survey.CreatorID, userName, token, groupID)
		}
	}

	return a.app.storage.CreateSurvey(survey)
}

func (a appShared) updateSurvey(survey model.Survey, userID string, userToken string) error {
	oldSurvey, err := a.app.storage.GetSurvey(survey.ID, survey.OrgID, survey.AppID)
	if err != nil {
		errors.WrapErrorAction(logutils.ActionDelete, logutils.TypeResult, nil, fmt.Errorf("cannot find survey"))
	}

	if oldSurvey.CreatorID == userID {
		return a.app.storage.UpdateSurvey(survey)
	}

	for _, groupID := range oldSurvey.GroupIDs {
		group, err := a.app.groups.GetGroupDetails(userToken, groupID)
		if err != nil {
			// TODO something
		}
		if err == nil && group.IsCurrentUserAdmin(userID) {
			return a.app.storage.UpdateSurvey(survey)
		}
	}

	return errors.WrapErrorAction(logutils.ActionUpdate, logutils.TypePermission, nil, fmt.Errorf("cannot edit survey"))
}

func (a appShared) deleteSurvey(id string, orgID string, appID string, userID string, userToken string) error {

	oldSurvey, err := a.app.storage.GetSurvey(id, orgID, appID)
	if err != nil {
		errors.WrapErrorAction(logutils.ActionDelete, logutils.TypeResult, nil, fmt.Errorf("cannot find survey"))
	}

	if oldSurvey.CreatorID == userID {
		return a.app.storage.DeleteSurvey(id, orgID, appID, &userID)
	}

	for _, groupID := range oldSurvey.GroupIDs {
		group, err := a.app.groups.GetGroupDetails(userToken, groupID)
		if err != nil {
			// TODO something
		}
		if err == nil && group.IsCurrentUserAdmin(userID) {
			return a.app.storage.DeleteSurvey(id, orgID, appID, &userID)
		}
	}

	// TODO return error
	return errors.WrapErrorAction(logutils.ActionDelete, logutils.TypePermission, nil, fmt.Errorf("cannot delete survey"))
}

func (a appShared) sendNotificationsToUserList(survey model.Survey, userID string, userName string) {
	messageRecipients := make([]model.NotificationMessageRecipient, len(survey.UserIDs))
	for i, userID := range survey.UserIDs {
		messageRecipients[i] = model.NotificationMessageRecipient{
			UserID: userID,
			// TODO: check mute
			Mute: false,
		}
	}

	a.app.notifications.SendNotification(model.NotificationMessage{
		OrgID:      survey.OrgID,
		AppID:      survey.AppID,
		Recipients: messageRecipients,
		Sender: &model.Sender{
			Type: "user",
			User: &model.UserRef{
				UserID: userID,
				Name:   userName,
			},
		},
		Topic:   "survey",
		Subject: "Illinois",
		Body:    fmt.Sprintf("Survey '%s' has been created", survey.Title),
		Data: map[string]string{
			"type":        survey.Type,
			"operation":   "survey_created",
			"entity_type": "survey",
			"entity_id":   survey.ID,
			"entity_name": survey.Title,
		},
	})
}

func (a appShared) sendNotificationsToGroup(survey model.Survey, userID string, userName string, userToken string, groupID string) {
	members, err := a.app.groups.GetGroupMembers(userToken, groupID)
	if err != nil {
		return
	}
	messageRecipients := make([]model.UserRef, len(*members))
	for i, member := range *members {
		messageRecipients[i] = model.UserRef{
			UserID: member.ClientID,
			Name:   member.Name,
		}
	}

	a.app.groups.SendGroupNotification(groupID, model.GroupNotification{
		Members: messageRecipients,
		Sender: &model.Sender{
			Type: "user",
			User: &model.UserRef{
				UserID: userID,
				Name:   userName,
			},
		},
		Topic:   "survey",
		Subject: "Illinois",
		Body:    fmt.Sprintf("Survey '%s' has been created", survey.Title),
		Data: map[string]string{
			"group_id":    groupID,
			"type":        "survey",
			"operation":   "survey_created",
			"entity_type": "survey",
			"entity_id":   survey.ID,
			"entity_name": survey.Title,
		},
	})
}

// newAppShared creates new appShared
func newAppShared(app *Application) appShared {
	return appShared{app: app}
}
