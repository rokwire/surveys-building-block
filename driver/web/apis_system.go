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
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rokwire/core-auth-library-go/v2/tokenauth"
	"github.com/rokwire/logging-library-go/v2/logs"
	"github.com/rokwire/logging-library-go/v2/logutils"
)

// SystemAPIsHandler handles the rest system admin APIs implementation
type SystemAPIsHandler struct {
	app *core.Application
}

func (h SystemAPIsHandler) getConfig(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	params := mux.Vars(r)
	id := params["id"]
	if len(id) <= 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypePathParam, logutils.StringArgs("id"), nil, http.StatusBadRequest, false)
	}

	configs, err := h.app.System.GetConfig(id)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionGet, model.TypeConfig, nil, err, http.StatusInternalServerError, true)
	}

	response, err := json.Marshal(configs)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionMarshal, model.TypeConfig, nil, err, http.StatusInternalServerError, false)
	}
	return l.HTTPResponseSuccessJSON(response)
}

func (h SystemAPIsHandler) saveConfig(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	params := mux.Vars(r)
	id := params["id"]
	if len(id) <= 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypePathParam, logutils.StringArgs("id"), nil, http.StatusBadRequest, false)
	}

	var requestData model.Config
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionUnmarshal, model.TypeConfig, nil, err, http.StatusBadRequest, true)
	}

	requestData.ID = id
	err = h.app.System.SaveConfig(requestData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionSave, model.TypeConfig, nil, err, http.StatusInternalServerError, true)
	}

	return l.HTTPResponseSuccess()
}

func (h SystemAPIsHandler) deleteConfig(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	params := mux.Vars(r)
	id := params["id"]
	if len(id) <= 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypePathParam, logutils.StringArgs("id"), nil, http.StatusBadRequest, false)
	}

	err := h.app.System.DeleteConfig(id)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionDelete, model.TypeConfig, nil, err, http.StatusInternalServerError, true)
	}
	return l.HTTPResponseSuccess()
}

// NewSystemAPIsHandler creates new system admin API handler instance
func NewSystemAPIsHandler(app *core.Application) SystemAPIsHandler {
	return SystemAPIsHandler{app: app}
}
