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
package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bharath-b-hpe/odimra/lib-utilities/common"
	"github.com/bharath-b-hpe/odimra/lib-utilities/config"
	sessionproto "github.com/bharath-b-hpe/odimra/lib-utilities/proto/session"
	"github.com/bharath-b-hpe/odimra/lib-utilities/response"
	"github.com/bharath-b-hpe/odimra/lib-utilities/services"
	"github.com/bharath-b-hpe/odimra/svc-api/router"
	"github.com/bharath-b-hpe/odimra/svc-api/rpc"
	iris "github.com/kataras/iris/v12"
)

func main() {
	// verifying the uid of the user
	if uid := os.Geteuid(); uid == 0 {
		log.Fatalln("Api Service should not be run as the root user")
	}
	router := router.Router()

	//WrapRouter method removes the trailing slash from the URL if present in the request and convert the URL to lower case.
	router.WrapRouter(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		path := r.URL.Path
		if len(path) > 1 && path[len(path)-1] == '/' && path[len(path)-2] != '/' {
			path = path[:len(path)-1]
			r.RequestURI = path
			r.URL.Path = path
		}
		basicAuth := r.Header.Get("Authorization")
		var basicAuthToken string
		if basicAuth != "" {
			var username, password string
			yes := strings.Contains(basicAuth, "Basic")
			if yes {
				spl := strings.Split(basicAuth, " ")
				data, err := base64.StdEncoding.DecodeString(spl[1])
				if err != nil {
					log.Println("error:", err)
					return
				}
				userCred := strings.SplitN(string(data), ":", 2)
				if len(userCred) < 2 {
					log.Println("error: not a valid basic auth")
					return
				}
				username = userCred[0]
				password = userCred[1]
			} else {
				log.Println("error: not a valid basic auth")
				return
			}

			var req sessionproto.SessionCreateRequest
			req.UserName = username
			req.Password = password
			resp, err := rpc.DoSessionCreationRequest(req)
			if err != nil && resp == nil {
				errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
				log.Println(errorMessage)
				w.Header().Set("Content-type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				body, _ := json.Marshal(common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil).Body)
				w.Write([]byte(body))
				return
			}
			if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
				errorMessage := "error: failed to create a sesssion"
				log.Println(errorMessage)
				w.Header().Set("Content-type", "application/json; charset=utf-8")
				w.WriteHeader(int(resp.StatusCode))
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
		// r.URL.Path = strings.ToLower(path)
		next(w, r)
	})

	err := config.SetConfiguration()
	if err != nil {
		log.Fatalf("fatal: %v", err)
	}

	err = services.InitializeService(services.APIClient)
	if err != nil {
		log.Fatalf("fatal: error while trying to initialize service: %v", err)
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
		log.Fatalf("fatal: error while initializing server: %v", err)
	}
	router.Run(iris.Server(apiServer))
}
