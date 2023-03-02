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

package web

import (
	"application/core"
	"application/core/model"
	"application/utils"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rokwire/core-auth-library-go/v2/tokenauth"
	"github.com/rokwire/logging-library-go/v2/logs"
	"github.com/rokwire/logging-library-go/v2/logutils"
)

// AnalyticsAPIsHandler handles the analytics rest APIs implementation
type AnalyticsAPIsHandler struct {
	app *core.Application
}

func (h AnalyticsAPIsHandler) getAnonymousSurveyResponses(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	// validate static token by comparing it against env config
	token, _, err := tokenauth.GetRequestTokens(r)
	if err != nil {
		return l.HTTPResponseErrorData(logutils.StatusInvalid, logutils.TypeToken, nil, nil, http.StatusUnauthorized, false)
	}
	envConfig, err := h.app.GetEnvConfigs()
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionGet, model.TypeConfig, logutils.StringArgs(model.ConfigTypeEnv), nil, http.StatusInternalServerError, false)
	}
	hashedToken := utils.SHA256Hash([]byte(token))
	if subtle.ConstantTimeCompare([]byte(base64.StdEncoding.EncodeToString(hashedToken)), []byte(envConfig.SplunkToken)) == 0 {
		return l.HTTPResponseErrorAction(logutils.ActionValidate, logutils.TypeToken, nil, nil, http.StatusUnauthorized, false)
	}

	surveyTypesRaw := r.URL.Query().Get("survey_types")
	var surveyTypes []string
	if len(surveyTypesRaw) > 0 {
		surveyTypes = strings.Split(surveyTypesRaw, ",")
	}

	timeOffsetRaw := r.URL.Query().Get("time_offset")
	var timeOffset int // hours
	if len(timeOffsetRaw) > 0 {
		intParsed, err := strconv.Atoi(timeOffsetRaw)
		if err != nil {
			return l.HTTPResponseErrorData(logutils.StatusInvalid, logutils.TypeQueryParam, logutils.StringArgs("time_offset"), nil, http.StatusBadRequest, false)
		}
		timeOffset = intParsed
	}

	startDateRaw := r.URL.Query().Get("start_date")
	var startDate *time.Time
	if len(startDateRaw) > 0 {
		dateParsed, err := time.Parse(time.RFC3339, startDateRaw)
		if err != nil {
			return l.HTTPResponseErrorData(logutils.StatusInvalid, logutils.TypeQueryParam, logutils.StringArgs("start_date"), nil, http.StatusBadRequest, false)
		}
		startDate = &dateParsed
	} else if timeOffset == 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypeQueryParam, &logutils.ListArgs{"start_date", "time_offset"}, nil, http.StatusBadRequest, false)
	}

	endDateRaw := r.URL.Query().Get("end_date")
	var endDate *time.Time
	if len(endDateRaw) > 0 {
		dateParsed, err := time.Parse(time.RFC3339, endDateRaw)
		if err != nil {
			return l.HTTPResponseErrorData(logutils.StatusInvalid, logutils.TypeQueryParam, logutils.StringArgs("end_date"), nil, http.StatusBadRequest, false)
		}
		endDate = &dateParsed
	} else if timeOffset == 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypeQueryParam, &logutils.ListArgs{"end_date", "time_offset"}, nil, http.StatusBadRequest, false)
	}

	if startDate == nil && endDate == nil {
		now := time.Now()
		offsetHours := now.Add(time.Duration(-timeOffset) * time.Hour)
		startDate = &offsetHours
		endDate = &now
	}

	resData, err := h.app.Analytics.GetSurveyResponses(surveyTypes, startDate, endDate)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionGet, model.TypeSurveyResponse, nil, err, http.StatusInternalServerError, true)
	}

	anonResData := make([]model.SurveyResponseAnonymous, len(resData))
	for i, surveyRes := range resData {
		anonResData[i] = model.SurveyResponseAnonymous{ID: surveyRes.Survey.ID, CreatorID: surveyRes.Survey.CreatorID, AppID: surveyRes.Survey.AppID,
			OrgID: surveyRes.Survey.OrgID, Title: surveyRes.Survey.Title, Type: surveyRes.Survey.Type, SurveyStats: surveyRes.Survey.SurveyStats,
			DateCreated: surveyRes.Survey.DateCreated, DateUpdated: surveyRes.Survey.DateUpdated}
	}

	data, err := json.Marshal(anonResData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionMarshal, logutils.TypeResponseBody, nil, err, http.StatusInternalServerError, false)
	}

	return l.HTTPResponseSuccessJSON(data)
}

// NewAnalyticsAPIsHandler creates new analytics API handler instance
func NewAnalyticsAPIsHandler(app *core.Application) AnalyticsAPIsHandler {
	return AnalyticsAPIsHandler{app: app}
}
