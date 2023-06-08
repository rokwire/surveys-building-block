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

func (a appShared) getAllSurveyResponses(id string, orgID string, appID string, userToken string, userID string, groupID string, startDate *time.Time, endDate *time.Time, limit *int, offset *int) ([]model.SurveyResponse, error) {
	group, err := a.app.groups.GetGroupDetails(userToken, groupID)
	if err != nil {
		return nil, nil
	}

	if group.IsCurrentUserAdmin(userID) || group.CreatorID == userID {
		return a.app.storage.GetAllSurveyResponses(&orgID, &appID, &id, startDate, endDate, limit, offset)
	}

	return nil, nil
}

func (a appShared) createSurvey(survey model.Survey, user model.User) (*model.Survey, error) {
	survey.ID = uuid.NewString()
	survey.DateCreated = time.Now().UTC()
	survey.DateUpdated = nil

	if len(survey.UserIDs) > 0 {
		sendNotificationsToUserList(a, survey, user)
	} else if len(survey.GroupIDs) > 0 {
		for _, groupID := range survey.GroupIDs {
			sendNotificationsToGroup(a, survey, user, groupID)
		}
	}

	return a.app.storage.CreateSurvey(survey)
}

func (a appShared) updateSurvey(survey model.Survey, admin bool) error {
	return a.app.storage.UpdateSurvey(survey, admin)
}

func (a appShared) deleteSurvey(id string, orgID string, appID string, creatorID *string) error {
	return a.app.storage.DeleteSurvey(id, orgID, appID, creatorID)
}

func sendNotificationsToUserList(a appShared, survey model.Survey, user model.User) {
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
				UserID: user.Claims.Subject,
				Name:   user.Claims.Name,
			},
		},
		Topic:   "survey",
		Subject: "Illinois",
		Body:    fmt.Sprintf("Survey '%s' has been created", survey.Title),
		Data: map[string]string{
			"type":        survey.Type,
			"operation":   "survey_created",
			"entity_type": "survey",
			"entity_id":   survey.ID.Hex(),
			"entity_name": survey.Title,
		},
	})
}

func sendNotificationsToGroup(a appShared, survey model.Survey, user model.User, groupID string) {
	members, err := a.app.groups.GetGroupMembers(user.Token, groupID)
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
				UserID: user.Claims.Subject,
				Name:   user.Claims.Name,
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
			"entity_id":   survey.ID.Hex(),
			"entity_name": survey.Title,
		},
	})
}

// newAppShared creates new appShared
func newAppShared(app *Application) appShared {
	return appShared{app: app}
}
