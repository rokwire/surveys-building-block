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
)

// appAdmin contains admin implementations
type appAdmin struct {
	app *Application
}

// Surveys
// GetSurvey returns the survey with the provided ID
func (a appAdmin) GetSurvey(id string, orgID string, appID string) (*model.Survey, error) {
	return a.app.shared.getSurvey(id, orgID, appID)
}

// CreateSurvey creates a new survey
func (a appAdmin) CreateSurvey(survey model.Survey) (*model.Survey, error) {
	return a.app.shared.createSurvey(survey)
}

// UpdateSurvey updates the provided survey
func (a appAdmin) UpdateSurvey(survey model.Survey) error {
	return a.app.shared.updateSurvey(survey, true)
}

// DeleteSurvey deletes the survey with the specified ID
func (a appAdmin) DeleteSurvey(id string, orgID string, appID string) error {
	return a.app.shared.deleteSurvey(id, orgID, appID, nil)
}

// GetAlertContacts returns all alert contacts for the provided app/org
func (a appAdmin) GetAlertContacts(orgID string, appID string) ([]model.AlertContact, error) {
	return a.app.storage.GetAlertContacts(orgID, appID)
}

// GetAlertContacts returns the alert contacts for the provided id
func (a appAdmin) GetAlertContact(id string, orgID string, appID string) (*model.AlertContact, error) {
	return a.app.storage.GetAlertContact(id, orgID, appID)
}

// CreateAlertContact creates a new alert contact
func (a appAdmin) CreateAlertContact(alertContact model.AlertContact) (*model.AlertContact, error) {
	return a.app.storage.CreateAlertContact(alertContact)
}

// UpdateAlertContact updates an existing alert contact
func (a appAdmin) UpdateAlertContact(alertContact model.AlertContact) error {
	return a.app.storage.UpdateAlertContact(alertContact)
}

// DeleteAlertContact deletes an existing alert contact with the provided id
func (a appAdmin) DeleteAlertContact(id string, orgID string, appID string) error {
	return a.app.storage.DeleteAlertContact(id, orgID, appID)
}

// newAppAdmin creates new appAdmin
func newAppAdmin(app *Application) appAdmin {
	return appAdmin{app: app}
}
