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
	"application/utils"

	"github.com/rokwire/core-auth-library-go/v3/authservice"
	"github.com/rokwire/core-auth-library-go/v3/authutils"
	"github.com/rokwire/core-auth-library-go/v3/coreservice"
	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logs"
	"github.com/rokwire/logging-library-go/v2/logutils"
)

type storageListener struct {
	app *Application
	model.DefaultStorageListener
}

// OnExampleUpdated notifies that the example collection has changed
func (s *storageListener) OnExampleUpdated() {
	s.app.logger.Infof("OnExampleUpdated")

	// TODO: Implement listener
}

// Application represents the core application code based on hexagonal architecture
type Application struct {
	version   string
	build     string
	serviceID string

	Default   interfaces.Default   // expose to the drivers adapters
	Client    interfaces.Client    // expose to the drivers adapters
	Admin     interfaces.Admin     // expose to the drivers adapters
	Analytics interfaces.Analytics // expose to the drivers adapters
	BBs       interfaces.BBs       // expose to the drivers adapters
	TPS       interfaces.TPS       // expose to the drivers adapters
	System    interfaces.System    // expose to the drivers adapters
	shared    Shared

	logger *logs.Logger

	storage       interfaces.Storage
	notifications interfaces.Notifications
	calendar      interfaces.Calendar

	coreService *coreservice.CoreService
}

// Start starts the core part of the application
func (a *Application) Start() {
	//set storage listener
	storageListener := storageListener{app: a}
	a.storage.RegisterStorageListener(&storageListener)

	if a.coreService != nil {
		a.coreService.StartDeletedMembershipsTimer()
	}
}

// GetEnvConfigs retrieves the cached database env configs
func (a *Application) GetEnvConfigs() (*model.EnvConfigData, error) {
	// Load env configs from database
	config, err := a.storage.FindConfig(model.ConfigTypeEnv, authutils.AllApps, authutils.AllOrgs)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionGet, model.TypeConfig, nil, err)
	}
	if config == nil {
		return nil, errors.ErrorData(logutils.StatusMissing, model.TypeConfig, &logutils.FieldArgs{"type": model.ConfigTypeEnv, "app_id": authutils.AllApps, "org_id": authutils.AllOrgs})
	}
	return model.GetConfigData[model.EnvConfigData](*config)
}

func (a *Application) handleDeletedMemberships(deleted []coreservice.DeletedOrgAppMemberships) error {
	for _, orgAppMemberships := range deleted {
		noContextAccountIDs := make([]string, 0)
		for _, context := range orgAppMemberships.Memberships {
			if context.Context != nil {
				retainData, err := a.getRetainSurveyResponsesData(context.Context)
				if err != nil {
					err = errors.WrapErrorAction(logutils.ActionGet, "retain survey responses data", &logutils.FieldArgs{"org_id": orgAppMemberships.OrgID, "app_id": orgAppMemberships.AppID, "account_id": context.AccountID}, err)
					a.logger.Error(err.Error())
					continue
				}
				if retainData == nil {
					continue
				}
				err = a.storage.DeleteSurveyResponsesExcept(orgAppMemberships.OrgID, orgAppMemberships.AppID, context.AccountID, retainData.SurveyIDs, retainData.SurveyTypes, nil, nil, true)
				if err != nil && a.logger != nil {
					err = errors.WrapErrorAction(logutils.ActionDelete, model.TypeSurveyResponse, &logutils.FieldArgs{"org_id": orgAppMemberships.OrgID, "app_id": orgAppMemberships.AppID}, err)
					a.logger.Error(err.Error())
				}
			} else {
				noContextAccountIDs = append(noContextAccountIDs, context.AccountID)
			}
		}

		if len(noContextAccountIDs) > 0 {
			err := a.storage.DeleteSurveyResponses(orgAppMemberships.OrgID, orgAppMemberships.AppID, noContextAccountIDs, nil, nil, nil, nil, true)
			if err != nil && a.logger != nil {
				err = errors.WrapErrorAction(logutils.ActionDelete, model.TypeSurveyResponse, &logutils.FieldArgs{"org_id": orgAppMemberships.OrgID, "app_id": orgAppMemberships.AppID}, err)
				a.logger.Error(err.Error())
			}
		}
	}

	return nil
}

func (a *Application) getRetainSurveyResponsesData(context map[string]interface{}) (*retainSurveyResponsesData, error) {
	surveyResponseContext, err := utils.JSONConvert[deleteSurveyResponsesContext](context)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionParse, "delete survey responses context", nil, err)
	}
	if surveyResponseContext == nil {
		return nil, nil
	}

	return &surveyResponseContext.RetainResponses, nil
}

type deleteSurveyResponsesContext struct {
	RetainResponses retainSurveyResponsesData `json:"retain_responses"`
}

type retainSurveyResponsesData struct {
	SurveyIDs   []string `json:"survey_ids"`
	SurveyTypes []string `json:"survey_types"`
}

// NewApplication creates new Application
func NewApplication(version string, build string, serviceID string, storage interfaces.Storage, notifications interfaces.Notifications, calendar interfaces.Calendar,
	serviceAccountManager *authservice.ServiceAccountManager, logger *logs.Logger) *Application {

	application := Application{version: version, build: build, serviceID: serviceID, storage: storage, notifications: notifications, calendar: calendar, logger: logger}

	var err error
	if serviceAccountManager != nil {
		deletedAccountsConfig := coreservice.DeletedMembershipsConfig{
			Callback: application.handleDeletedMemberships,
		}
		application.coreService, err = coreservice.NewCoreService(serviceAccountManager, &deletedAccountsConfig, logger)
		if err != nil && logger != nil {
			logger.Errorf("error creating core service: %v", err)
		}
	}

	//add the drivers ports/interfaces
	application.Default = newAppDefault(&application)
	application.Client = newAppClient(&application)
	application.Admin = newAppAdmin(&application)
	application.Analytics = newAppAnalytics(&application)
	application.BBs = newAppBBs(&application)
	application.TPS = newAppTPS(&application)
	application.System = newAppSystem(&application)
	application.shared = newAppShared(&application)

	return &application
}
