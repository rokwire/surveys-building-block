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
	"github.com/rokwire/core-auth-library-go/v3/authutils"
	"github.com/rokwire/core-auth-library-go/v3/tokenauth"
	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logutils"
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

// GetSurvey returns surveys matching the provided query
func (a appAdmin) GetSurveys(orgID string, appID string, surveyIDs []string, surveyTypes []string, limit *int, offset *int) ([]model.Survey, error) {
	return a.app.shared.getSurveys(orgID, appID, surveyIDs, surveyTypes, limit, offset)
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
	alertContact.ID = uuid.NewString()
	alertContact.DateCreated = time.Now().UTC()
	alertContact.DateUpdated = nil
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

func (a appAdmin) GetConfig(id string, claims *tokenauth.Claims) (*model.Config, error) {
	config, err := a.app.storage.FindConfigByID(id)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, model.TypeConfig, nil, err)
	}
	if config == nil {
		return nil, errors.ErrorData(logutils.StatusMissing, model.TypeConfig, &logutils.FieldArgs{"id": id})
	}

	err = claims.CanAccess(config.AppID, config.OrgID, config.System)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionValidate, "config access", nil, err)
	}

	return config, nil
}

func (a appAdmin) GetConfigs(configType *string, claims *tokenauth.Claims) ([]model.Config, error) {
	configs, err := a.app.storage.FindConfigs(configType)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, model.TypeConfig, nil, err)
	}

	allowedConfigs := make([]model.Config, 0)
	for _, config := range configs {
		if err := claims.CanAccess(config.AppID, config.OrgID, config.System); err == nil {
			allowedConfigs = append(allowedConfigs, config)
		}
	}
	return allowedConfigs, nil
}

func (a appAdmin) CreateConfig(config model.Config, claims *tokenauth.Claims) (*model.Config, error) {
	// must be a system config if applying to all orgs
	if config.OrgID == authutils.AllOrgs && !config.System {
		return nil, errors.ErrorData(logutils.StatusInvalid, "config system status", &logutils.FieldArgs{"config.org_id": authutils.AllOrgs})
	}

	err := claims.CanAccess(config.AppID, config.OrgID, config.System)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionValidate, "config access", nil, err)
	}

	config.ID = uuid.NewString()
	config.DateCreated = time.Now().UTC()
	err = a.app.storage.InsertConfig(config)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionInsert, model.TypeConfig, nil, err)
	}
	return &config, nil
}

func (a appAdmin) UpdateConfig(config model.Config, claims *tokenauth.Claims) error {
	// must be a system config if applying to all orgs
	if config.OrgID == authutils.AllOrgs && !config.System {
		return errors.ErrorData(logutils.StatusInvalid, "config system status", &logutils.FieldArgs{"config.org_id": authutils.AllOrgs})
	}

	oldConfig, err := a.app.storage.FindConfig(config.Type, config.AppID, config.OrgID)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionFind, model.TypeConfig, nil, err)
	}
	if oldConfig == nil {
		return errors.ErrorData(logutils.StatusMissing, model.TypeConfig, &logutils.FieldArgs{"type": config.Type, "app_id": config.AppID, "org_id": config.OrgID})
	}

	// cannot update a system config if not a system admin
	if !claims.System && oldConfig.System {
		return errors.ErrorData(logutils.StatusInvalid, "system claim", nil)
	}
	err = claims.CanAccess(config.AppID, config.OrgID, config.System)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionValidate, "config access", nil, err)
	}

	now := time.Now().UTC()
	config.ID = oldConfig.ID
	config.DateUpdated = &now

	err = a.app.storage.UpdateConfig(config)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionUpdate, model.TypeConfig, nil, err)
	}
	return nil
}

func (a appAdmin) DeleteConfig(id string, claims *tokenauth.Claims) error {
	config, err := a.app.storage.FindConfigByID(id)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionFind, model.TypeConfig, nil, err)
	}
	if config == nil {
		return errors.ErrorData(logutils.StatusMissing, model.TypeConfig, &logutils.FieldArgs{"id": id})
	}

	err = claims.CanAccess(config.AppID, config.OrgID, config.System)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionValidate, "config access", nil, err)
	}

	err = a.app.storage.DeleteConfig(id)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionDelete, model.TypeConfig, nil, err)
	}
	return nil
}

// newAppAdmin creates new appAdmin
func newAppAdmin(app *Application) appAdmin {
	return appAdmin{app: app}
}
