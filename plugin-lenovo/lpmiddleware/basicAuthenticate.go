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

//Package lpmiddleware ...
package lpmiddleware

import (
	"encoding/base64"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/config"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"

	"net/http"
	"strings"
)

//BasicAuth is used to validate REST API calls with plugin with basic autherization
func BasicAuth(ctx iris.Context) {
	basicAuth := ctx.GetHeader("Authorization")
	if basicAuth != "" {
		var username, password string
		yes := strings.Contains(basicAuth, "Basic")
		if yes {
			spl := strings.Split(basicAuth, " ")
			data, err := base64.StdEncoding.DecodeString(spl[1])
			if err != nil {
				log.Error(err.Error())
				ctx.StatusCode(http.StatusInternalServerError)
				ctx.WriteString(err.Error())
				return
			}
			userCred := strings.SplitN(string(data), ":", 2)
			if len(userCred) < 2 {
				log.Error("Not a valid basic auth")
				ctx.StatusCode(http.StatusUnauthorized)
				ctx.WriteString("Not a valid basic auth")
				return
			}
			username = userCred[0]
			password = userCred[1]
		} else {
			log.Error("Not a valid basic auth")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Not a valid basic auth")
			return
		}
		userName := config.Data.PluginConf.UserName
		passwd := config.Data.PluginConf.Password
		if username != userName {
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid Username/Password")
			return
		}
		hash := sha3.New512()
		hash.Write([]byte(password))
		hashSum := hash.Sum(nil)
		hashedPassword := base64.URLEncoding.EncodeToString(hashSum)
		if passwd != hashedPassword {
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid Username/Password")
			return
		}
	}
	ctx.Next()

}
