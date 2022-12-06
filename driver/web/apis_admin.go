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

// AdminAPIsHandler handles the rest Admin APIs implementation
type AdminAPIsHandler struct {
	app *core.Application
}

func (h AdminAPIsHandler) getExample(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	params := mux.Vars(r)
	id := params["id"]
	if len(id) <= 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypePathParam, logutils.StringArgs("id"), nil, http.StatusBadRequest, false)
	}

	example, err := h.app.Admin.GetExample(claims.OrgID, claims.AppID, id)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionGet, model.TypeExample, nil, err, http.StatusInternalServerError, true)
	}

	response, err := json.Marshal(example)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionMarshal, model.TypeExample, nil, err, http.StatusInternalServerError, false)
	}
	return l.HTTPResponseSuccessJSON(response)
}

func (h AdminAPIsHandler) createExample(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	var requestData model.Example
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionUnmarshal, model.TypeExample, nil, err, http.StatusBadRequest, true)
	}

	requestData.OrgID = claims.OrgID
	requestData.AppID = claims.AppID
	example, err := h.app.Admin.CreateExample(requestData)
	if err != nil || example == nil {
		return l.HTTPResponseErrorAction(logutils.ActionCreate, model.TypeExample, nil, err, http.StatusInternalServerError, true)
	}

	response, err := json.Marshal(example)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionMarshal, model.TypeExample, nil, err, http.StatusInternalServerError, false)
	}
	return l.HTTPResponseSuccessJSON(response)
}

func (h AdminAPIsHandler) updateExample(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	params := mux.Vars(r)
	id := params["id"]
	if len(id) <= 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypePathParam, logutils.StringArgs("id"), nil, http.StatusBadRequest, false)
	}

	var requestData model.Example
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionUnmarshal, model.TypeExample, nil, err, http.StatusBadRequest, true)
	}

	requestData.ID = id
	requestData.OrgID = claims.OrgID
	requestData.AppID = claims.AppID
	err = h.app.Admin.UpdateExample(requestData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionUpdate, model.TypeExample, nil, err, http.StatusInternalServerError, true)
	}

	return l.HTTPResponseSuccess()
}

func (h AdminAPIsHandler) deleteExample(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	params := mux.Vars(r)
	id := params["id"]
	if len(id) <= 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypePathParam, logutils.StringArgs("id"), nil, http.StatusBadRequest, false)
	}

	err := h.app.Admin.DeleteExample(claims.OrgID, claims.AppID, id)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionDelete, model.TypeExample, nil, err, http.StatusInternalServerError, true)
	}

	return l.HTTPResponseSuccess()
}

// NewAdminAPIsHandler creates new rest Handler instance
func NewAdminAPIsHandler(app *core.Application) AdminAPIsHandler {
	return AdminAPIsHandler{app: app}
}
