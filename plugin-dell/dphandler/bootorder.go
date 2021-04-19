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

//Package dphandler ...
package dphandler

import (
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

// SetDefaultBootOrder : sets the defult boot order
func SetDefaultBootOrder(ctx iris.Context) {

	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	//Validating the token
	if token != "" {
		flag := TokenValidation(token)
		if !flag {
			log.Error("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}

	var deviceDetails dpmodel.Device
	uri := ctx.Request().RequestURI
	//replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	}
	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		log.Error("While trying to collect data from request: " + err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Error: bad request.")
		return
	}

	device := &dputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
	}
	/*
		priv := []byte(dpmodel.PluginPrivateKey)
		block, _ := pem.Decode(priv)
		enc := x509.IsEncryptedPEMBlock(block)
		b := block.Bytes
		if enc {
			log.Println("is encrypted pem block")
			b, err = x509.DecryptPEMBlock(block, nil)
			if err != nil {
				log.Println("Error: " + err.Error())
			}
		}
		key, err := x509.ParsePKCS1PrivateKey(b)
		if err != nil {
			log.Println("Error: " + err.Error())
		}

		hash := sha512.New()

		plainText, err := rsa.DecryptOAEP(
			hash,
			rand.Reader,
			key,
			device.Password,
			nil,
		)
		if err != nil {
			log.Println("Error while trying decrypt data: ", err)
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString("Error while trying to decypt data")
			return
		}

		device.Password = plainText
	*/
	redfishClient, err := dputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to create the redfish client, got:" + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}

	//Subscribe to Events
	resp, err := redfishClient.SetDefaultBootOrder(device, uri)
	if err != nil {
		errorMessage := "while trying to set default boot order, got: " + err.Error()
		log.Error(errorMessage)
		if resp == nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(errorMessage)
			return
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		body = []byte("while trying to set default boot order, got: " + err.Error())
		log.Error(string(body))
	}

	ctx.StatusCode(resp.StatusCode)
	ctx.Write(body)
}
