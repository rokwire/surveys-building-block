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

// appAnalytics contains analytics implementations
type appAnalytics struct {
	app *Application
}

// Survey Response
// GetSurveyResponses returns the survey responses matching the provided filters
func (a appAnalytics) GetSurveyResponses(surveyType string, startDate *time.Time, endDate *time.Time) ([]model.SurveyResponse, error) {
	return a.app.storage.GetSurveyResponses(nil, nil, nil, nil, []string{surveyType}, startDate, endDate, nil, nil)
}

// newAppAnalytics creates new appAnalytics
func newAppAnalytics(app *Application) appAnalytics {
	return appAnalytics{app: app}
}
