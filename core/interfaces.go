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
	"application/driven/storage"
	"time"
)

// Default exposes client APIs for the driver adapters
type Default interface {
	GetVersion() string
}

// Client exposes client APIs for the driver adapters
type Client interface {
	// Surveys
	GetSurvey(id string, orgID string, appID string) (*model.Survey, error)
	GetSurveys(orgID string, appID string, surveyIDs []string, surveyTypes []string, limit *int, offset *int) ([]model.Survey, error)
	CreateSurvey(survey model.Survey) (*model.Survey, error)
	UpdateSurvey(survey model.Survey) error
	DeleteSurvey(id string, orgID string, appID string, userID string) error

	// Survey Response
	GetSurveyResponse(id string, orgID string, appID string, userID string) (*model.SurveyResponse, error)
	GetSurveyResponses(orgID string, appID string, userID string, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time, limit *int, offset *int) ([]model.SurveyResponse, error)
	CreateSurveyResponse(surveyResponse model.SurveyResponse) (*model.SurveyResponse, error)
	UpdateSurveyResponse(surveyResponse model.SurveyResponse) error
	DeleteSurveyResponse(id string, orgID string, appID string, userID string) error
	DeleteSurveyResponses(orgID string, appID string, userID string, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time) error

	// Survey Alerts
	CreateSurveyAlert(surveyAlert model.SurveyAlert) error
}

// Admin exposes administrative APIs for the driver adapters
type Admin interface {
	// Surveys
	GetSurvey(id string, orgID string, appID string) (*model.Survey, error)
	GetSurveys(orgID string, appID string, surveyIDs []string, surveyTypes []string, limit *int, offset *int) ([]model.Survey, error)
	CreateSurvey(survey model.Survey) (*model.Survey, error)
	UpdateSurvey(survey model.Survey) error
	DeleteSurvey(id string, orgID string, appID string) error

	// Alert Contacts
	GetAlertContacts(orgID string, appID string) ([]model.AlertContact, error)
	GetAlertContact(id string, orgID string, appID string) (*model.AlertContact, error)
	CreateAlertContact(alertContact model.AlertContact) (*model.AlertContact, error)
	UpdateAlertContact(alertContact model.AlertContact) error
	DeleteAlertContact(id string, orgID string, appID string) error
}

// BBs exposes Building Block APIs for the driver adapters
type BBs interface {
}

// TPS exposes third-party service APIs for the driver adapters
type TPS interface {
}

// System exposes system administrative APIs for the driver adapters
type System interface {
	GetConfig(id string) (*model.Config, error)
	SaveConfig(configs model.Config) error
	DeleteConfig(id string) error
}

// Shared exposes shared APIs for other interface implementations
type Shared interface {
	// Surveys
	getSurvey(id string, orgID string, appID string) (*model.Survey, error)
	getSurveys(orgID string, appID string, surveyIDs []string, surveyTypes []string, limit *int, offset *int) ([]model.Survey, error)
	createSurvey(survey model.Survey) (*model.Survey, error)
	updateSurvey(survey model.Survey, admin bool) error
	deleteSurvey(id string, orgID string, appID string, creatorID *string) error
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	RegisterStorageListener(storageListener storage.Listener)
	PerformTransaction(func(adapter storage.Adapter) error) error

	GetConfig(id string) (*model.Config, error)
	SaveConfig(configs model.Config) error
	DeleteConfig(id string) error

	GetSurvey(id string, orgID string, appID string) (*model.Survey, error)
	GetSurveys(orgID string, appID string, surveyIDs []string, surveyTypes []string, limit *int, offset *int) ([]model.Survey, error)
	CreateSurvey(survey model.Survey) (*model.Survey, error)
	UpdateSurvey(survey model.Survey, admin bool) error
	DeleteSurvey(id string, orgID string, appID string, creatorID *string) error

	GetSurveyResponse(id string, orgID string, appID string, userID string) (*model.SurveyResponse, error)
	GetSurveyResponses(orgID string, appID string, userID string, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time, limit *int, offset *int) ([]model.SurveyResponse, error)
	CreateSurveyResponse(surveyResponse model.SurveyResponse) (*model.SurveyResponse, error)
	UpdateSurveyResponse(surveyResponse model.SurveyResponse) error
	DeleteSurveyResponse(id string, orgID string, appID string, userID string) error
	DeleteSurveyResponses(orgID string, appID string, userID string, surveyIDs []string, surveyTypes []string, startDate *time.Time, endDate *time.Time) error

	GetAlertContacts(orgID string, appID string) ([]model.AlertContact, error)
	GetAlertContact(id string, orgID string, appID string) (*model.AlertContact, error)
	GetAlertContactsByKey(key string, orgID string, appID string) ([]model.AlertContact, error)
	CreateAlertContact(alertContact model.AlertContact) (*model.AlertContact, error)
	UpdateAlertContact(alertContact model.AlertContact) error
	DeleteAlertContact(id string, orgID string, appID string) error
}

// Notifications is the interface for accessing the Notifications BB
type Notifications interface {
	SendNotification(notification model.NotificationMessage)
	SendMail(toEmail string, subject string, body string)
}
