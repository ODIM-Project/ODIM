// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/logs"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-api/apicommon"
	"github.com/ODIM-Project/ODIM/svc-api/router"
	"github.com/ODIM-Project/ODIM/svc-api/rpc"
	iris "github.com/kataras/iris/v12"
)

func main() {
	// setting up the logging framework
	hostName := os.Getenv("HOST_NAME")
	podName := os.Getenv("POD_NAME")
	pid := os.Getpid()
	logs.Adorn(logrus.Fields{
		"host":   hostName,
		"procid": podName + fmt.Sprintf("_%d", pid),
	})

	// log should be initialized after Adorn is invoked
	// as Adorn will assign new pointer to Log variable in logs package.
	log := logs.Log
	configWarnings, err := config.SetConfiguration()
	if err != nil {
		log.Logger.SetFormatter(&logs.SysLogFormatter{})
		log.Fatal("Error while trying set up configuration: " + err.Error())
	}
	logs.SetFormatter(config.Data.LogFormat)
	log.Logger.SetOutput(os.Stdout)
	log.Logger.SetLevel(config.Data.LogLevel)

	// verifying the uid of the user
	if uid := os.Geteuid(); uid == 0 {
		log.Fatal("Api Service should not be run as the root user")
	}
	router := router.Router()

	//WrapRouter method removes the trailing slash from the URL if present in the request and convert the URL to lower case.
	router.WrapRouter(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		path := r.URL.Path
		path = strings.Replace(strings.Replace(path, "\n", "", -1), "\r", "", -1)
		if len(path) > 1 && path[len(path)-1] == '/' && path[len(path)-2] != '/' {
			path = path[:len(path)-1]
			r.RequestURI = path
			r.URL.Path = path
		}
		// generating transaction ID
		transactionID := uuid.New()
		ctx := createContext(r, transactionID, podName)
		r = r.WithContext(ctx)
		basicAuth := r.Header.Get("Authorization")
		var basicAuthToken string

		if basicAuth != "" {
			var urlNoBasicAuth = []string{"/redfish/v1", "/redfish/v1/SessionService"}
			var authRequired bool
			authRequired = true
			for _, item := range urlNoBasicAuth {
				if item == path {
					authRequired = false
					logs.LogWithFields(ctx).Warn("Basic auth is provided but not used as URL is: " + path)
					break
				}
			}
			if authRequired {
				var username, password string
				yes := strings.Contains(basicAuth, "Basic")
				if yes {
					spl := strings.Split(basicAuth, " ")
					if len(spl) != 2 {
						errorMessage := "Invalid basic auth provided"
						logs.LogWithFields(ctx).Error(errorMessage)
						invalidAuthResp(errorMessage, w)
						return
					}
					data, err := base64.StdEncoding.DecodeString(spl[1])
					if err != nil {
						errorMessage := "Decoding the authorization failed: " + err.Error()
						logs.LogWithFields(ctx).Error(err.Error())
						invalidAuthResp(errorMessage, w)
						return
					}
					userCred := strings.SplitN(string(data), ":", 2)
					if len(userCred) < 2 {
						errorMessage := "Invalid basic auth provided"
						logs.LogWithFields(ctx).Error(errorMessage)
						invalidAuthResp(errorMessage, w)
						return
					}
					username = userCred[0]
					password = userCred[1]
				} else {
					errorMessage := "Invalid basic auth provided"
					logs.LogWithFields(ctx).Error(errorMessage)
					invalidAuthResp(errorMessage, w)
					return
				}

				//Converting the request into a map
				sessionReq := map[string]interface{}{
					"UserName": username,
					"Password": password,
				}
				//Marshalling input to get bytes since session create request accepts bytes
				sessionReqData, err := json.Marshal(sessionReq)

				var req sessionproto.SessionCreateRequest
				req.RequestBody = sessionReqData
				resp, err := rpc.DoSessionCreationRequest(ctx, req)
				if err != nil && resp == nil {
					errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
					logs.LogWithFields(ctx).Error(errorMessage)
					common.SetCommonHeaders(w)
					w.WriteHeader(http.StatusInternalServerError)
					body, _ := json.Marshal(common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil).Body)
					w.Write([]byte(body))
					return
				}
				if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
					common.SetCommonHeaders(w)
					w.WriteHeader(int(resp.StatusCode))
					if resp.StatusCode == http.StatusServiceUnavailable {
						logs.LogWithFields(ctx).Error("error: unable to establish connection with db")
						w.Write(resp.Body)
						return
					}
					errorMessage := "error: failed to create a sesssion"
					logs.LogWithFields(ctx).Info(errorMessage)
					body, _ := json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body)
					w.Write([]byte(body))
					return
				}
				var sessionHeader map[string]string
				var sessionID string
				sessionHeader = resp.Header
				basicAuthToken = sessionHeader["X-Auth-Token"]
				sessionLocation := sessionHeader["Link"]
				sessionLocationSlice := strings.Split(sessionLocation, "/")
				if len(sessionLocationSlice) > 1 {
					sessionID = sessionLocationSlice[len(sessionLocationSlice)-2]
				}
				r.Header.Set("X-Auth-Token", basicAuthToken)
				r.Header.Set("Session-ID", sessionID)
			}
		}
		// r.URL.Path = strings.ToLower(path)
		next(w, r)
	})

	// TODO: uncomment the following line after the migration
	config.CollectCLArgs(&configWarnings)
	for _, warning := range configWarnings {
		log.Warn(warning)
	}

	err = services.InitializeClient(services.APIClient)
	if err != nil {
		logs.Log.Fatal("service initialization failed: " + err.Error())
	}

	conf := &config.HTTPConfig{
		Certificate:   &config.Data.APIGatewayConf.Certificate,
		PrivateKey:    &config.Data.APIGatewayConf.PrivateKey,
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
		ServerAddress: config.Data.APIGatewayConf.Host,
		ServerPort:    config.Data.APIGatewayConf.Port,
	}
	apiServer, err := conf.GetHTTPServerObj()
	if err != nil {
		logs.Log.Fatal("service initialization failed: " + err.Error())
	}

	apicommon.ConfigFilePath = os.Getenv("CONFIG_FILE_PATH")
	if apicommon.ConfigFilePath == "" {
		logs.Log.Fatal("error: no value get the environment variable CONFIG_FILE_PATH")
	}

	errChan := make(chan error)
	// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
	go apicommon.TrackConfigFileChanges(errChan)

	router.Run(iris.Server(apiServer))
}

// invalidAuthResp function is used to generate an invalid credentials response
func invalidAuthResp(errMsg string, w http.ResponseWriter) {
	common.SetCommonHeaders(w)
	w.WriteHeader(http.StatusUnauthorized)
	body, _ := json.Marshal(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil).Body)
	w.Write([]byte(body))
}

// getContext is used to create new context with fields which are required for logging (transcationID, actionID, actionName,
// ThreadID, ThreadName and ProcessName)
func createContext(r *http.Request, transactionID uuid.UUID, podName string) context.Context {
	ctx := context.Background()
	var serviceName string
	val := strings.Split(r.URL.Path, "/")
	if len(val) >= 4 && val[2] != "" {
		serviceName = val[3]
	} else {
		serviceName = ""
	}
	// Add Action ID and Action Name in logs
	if action, ok := common.Actions[common.ActionKey{Service: serviceName, Uri: val[len(val)-1], Method: r.Method}]; ok {
		ctx = context.WithValue(ctx, common.ActionName, action.ActionName)
		ctx = context.WithValue(ctx, common.ActionID, action.ActionID)
	} else if action, ok := common.Actions[common.ActionKey{Service: serviceName, Uri: val[len(val)-2] + "/{id}", Method: r.Method}]; ok {
		ctx = context.WithValue(ctx, common.ActionName, action.ActionName)
		ctx = context.WithValue(ctx, common.ActionID, action.ActionID)
	} else {
		ctx = context.WithValue(ctx, common.ActionName, common.InvalidActionName)
		ctx = context.WithValue(ctx, common.ActionID, common.InvalidActionID)
	}
	// Add values in context (TransactionID, ThreadName, ThreadID)
	ctx = context.WithValue(ctx, common.TransactionID, transactionID.String())
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	ctx = context.WithValue(ctx, common.ThreadName, common.ApiService)
	ctx = context.WithValue(ctx, common.ThreadID, common.DefaultThreadID)

	return ctx
}
