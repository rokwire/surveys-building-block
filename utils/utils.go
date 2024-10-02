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

package utils

import (
	"crypto/sha256"
	"encoding/json"
	"reflect"
	"time"

	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logutils"
)

// GetInt gives the value which this pointer points. Gives 0 if the pointer is nil
func GetInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

// GetString gives the value which this pointer points. Gives empty string if the pointer is nil
func GetString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// GetTime gives the value which this pointer points. Gives empty string if the pointer is nil
func GetTime(time *time.Time) string {
	if time == nil {
		return ""
	}
	return time.String()
}

// SHA256Hash computes the SHA256 hash of a byte slice
func SHA256Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// JSONConvert json marshals and unmarshals data into result (result should be passed as a pointer)
func JSONConvert[T any, F any](val F) (*T, error) {
	if IsNil(val) {
		return nil, nil
	}

	bytes, err := json.Marshal(val)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionMarshal, "value", nil, err)
	}

	var out T
	err = json.Unmarshal(bytes, &out)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionUnmarshal, "value", nil, err)
	}

	return &out, nil
}

// IsNil determines whether the given interface has a nil value
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
