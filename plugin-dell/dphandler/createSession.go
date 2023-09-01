//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

// Package dphandler ...
package dphandler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"

	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	iris "github.com/kataras/iris/v12"
	"golang.org/x/crypto/sha3"
)

// TokenMap is used to define the plugin generated tokens
type TokenMap struct {
	Token    string
	LastUsed time.Time
}

var tokenDetails []TokenMap
var tokenSpec TokenMap

// CreateSession is used to create session for odimra to interact with plugin
func CreateSession(ctx iris.Context) {
	ctxt := ctx.Request().Context()
	var userCreds dpmodel.Users
	rawBodyAsBytes, err := IoUtilReadAll(ctx.Request().Body)
	if err != nil {
		errorMessage := "While trying to validate the credentials, got: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString(errorMessage)
	}
	err = json.Unmarshal(rawBodyAsBytes, &userCreds)
	if err != nil {
		errorMessage := "While trying to unmarshal user details, got: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errorMessage)
	}
	//Validate the credentials
	userName := userCreds.Username
	password := userCreds.Password
	validateResponse := validate(ctxt, userName, password)
	if !validateResponse {
		errorMessage := "Invalid credentials for session creation"
		l.LogWithFields(ctxt).Error(errorMessage)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.WriteString(errorMessage)
		return
	}
	//Create token
	token := createToken()
	currentTime := time.Now()
	tokenSpec.Token = token
	tokenSpec.LastUsed = currentTime
	tokenDetails = append(tokenDetails, tokenSpec)
	ctx.StatusCode(http.StatusCreated)
	ctx.Header("X-Auth-Token", token)
}

func validate(ctxt context.Context, userName, password string) bool {
	//var err error
	username := pluginConfig.Data.PluginConf.UserName
	passwd := pluginConfig.Data.PluginConf.Password
	if username != userName {
		return false
	}
	hash := sha3.New512()
	hash.Write([]byte(password))
	hashSum := hash.Sum(nil)
	hashedPassword := base64.URLEncoding.EncodeToString(hashSum)
	if passwd != hashedPassword {
		l.LogWithFields(ctxt).Error("Username/Password does not match")
		return false
	}

	return true
}

func createToken() string {
	token := uuid.NewV4().String()
	return token
}
