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
	"encoding/json"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"time"
	//"github.com/ODIM-Project/ODIM/lib-utilities/common"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpresponse"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
)

// GetPluginStatus defines the GetPluginStatus iris handler.
// and returns status
func GetPluginStatus(ctx iris.Context) {
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
	var messageQueueInfo []dpresponse.EmbQueue
	var resp = dpresponse.PluginStatusResponse{
		Comment: "Plugin Status Response",
		Name:    "Common Redfish Plugin Status",
		Version: pluginConfig.Data.FirmwareVersion,
	}
	resp.Status = dputilities.Status
	resp.Status.TimeStamp = time.Now().Format(time.RFC3339)
	resp.EventMessageBus = dpresponse.EventMessageBus{
		EmbType: pluginConfig.Data.MessageBusConf.EmbType,
	}
	//messageQueueInfo := make([]dpresponse.EmbQueue, 0)
	for i := 0; i < len(pluginConfig.Data.MessageBusConf.EmbQueue); i++ {
		messageQueueInfo = append(messageQueueInfo, dpresponse.EmbQueue{
			QueueName: pluginConfig.Data.MessageBusConf.EmbQueue[i],
			QueueDesc: "Queue for redfish events",
		})
	}
	resp.EventMessageBus.EmbQueue = messageQueueInfo

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)

}

// GetPluginStartup ...
func GetPluginStartup(ctx iris.Context) {
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

	var startup dpmodel.StartUpData
	err := ctx.ReadJSON(&startup)
	if err != nil {
		log.Error("While trying to collect data from request, got: " + err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Error: bad request.")
		return
	}

	if len(startup.Devices) <= 0 {
		log.Info("startup devices list is empty")
		ctx.StatusCode(http.StatusOK)
		return
	}

	errorCh := make(chan error)
	startUpResponse := make(chan map[string]string)
	respBody := make(map[string]string)
	quit := make(chan bool)
	var writeWG sync.WaitGroup
	go func() {
		for {
			select {
			case err = <-errorCh:
				ctx.StatusCode(http.StatusInternalServerError)
				ctx.WriteString(err.Error())
				//close(startUpResponse)
				//close(respHeader)
				//close(errorCh)
				writeWG.Done()
				return
			case startResp := <-startUpResponse:
				for k, v := range startResp {
					respBody[k] = v
				}
				writeWG.Done()
			case <-quit:
				//close(startUpResponse)
				//close(respHeader)
				//close(errorCh)
				//close(quit)
				break
			}
		}
	}()

	for uuid, device := range startup.Devices {
		if device.Operation == "add" {
			dpmodel.AddDeviceToInventory(uuid, device)
		}
		if device.Operation == "del" {
			dpmodel.DeleteDeviceInInventory(uuid)
		}
		if startup.ResyncEvtSubscription && startup.RequestType == "full" {
			writeWG.Add(1)
			go checkCreateSub(device, startUpResponse, errorCh, &writeWG)
		}
	}

	writeWG.Wait()
	quit <- true
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(respBody)
	return
}

func checkCreateSub(server dpmodel.DeviceData, startUpResponse chan map[string]string, errorCh chan error, writeWG *sync.WaitGroup) {
	var respBody = make(map[string]string)

	device := &dputilities.RedfishDevice{
		Host:     server.Address,
		Username: server.UserName,
		Password: string(server.Password),
		Location: server.EventSubscriptionInfo.Location,
	}
	redfishClient, err := dputilities.GetRedfishClient()
	if err != nil {
		errorCh <- err
		return
	}

	//Get Subscription details
	resp, err := redfishClient.GetSubscriptionDetail(device)
	if err != nil {
		errorCh <- err
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorCh <- err
			return
		}
		var obj dpmodel.EvtSubPost
		if err = json.Unmarshal([]byte(body), &obj); err != nil {
			errorCh <- err
			return
		}

		res := reflect.DeepEqual(obj.EventTypes, server.EventSubscriptionInfo.EventTypes)
		if !res {
			//Delete Subscription details
			resp, err := redfishClient.DeleteSubscriptionDetail(device)
			if err != nil {
				errorCh <- err
				return
			}
			defer resp.Body.Close()

			//Create new Subscription with details in odimra
			req := dpmodel.EvtSubPost{
				Destination: "https://" + pluginConfig.Data.LoadBalancerConf.Host + ":" + pluginConfig.Data.LoadBalancerConf.Port + pluginConfig.Data.EventConf.DestURI,
				EventTypes:  server.EventSubscriptionInfo.EventTypes,
				Context:     "Event Subscription",
				//      HTTPHeaders: reqPostBody.HTTPHeaders,
				Protocol: "Redfish",
			}
			device.PostBody, err = json.Marshal(req)
			if err != nil {
				errorCh <- err
				return
			}

			//Subscribe to Events
			resp, err = redfishClient.SubscribeForEvents(device)
			if err != nil {
				errorCh <- err
				return
			}
			defer resp.Body.Close()

		}

	} else if resp.StatusCode == http.StatusNotFound {
		req := dpmodel.EvtSubPost{
			Destination: "https://" + pluginConfig.Data.LoadBalancerConf.Host + ":" + pluginConfig.Data.LoadBalancerConf.Port + pluginConfig.Data.EventConf.DestURI,
			EventTypes:  []string{"Alert"},
			Context:     "Event Subscription",
			//	HTTPHeaders: reqPostBody.HTTPHeaders,
			Protocol: "Redfish",
		}
		device.PostBody, err = json.Marshal(req)
		if err != nil {
			errorCh <- err
			return
		}

		//Subscribe to Events
		resp, err = redfishClient.SubscribeForEvents(device)
		if err != nil {
			errorCh <- err
			return
		}
		defer resp.Body.Close()
	}

	respBody[device.Host] = resp.Header.Get("location")
	startUpResponse <- respBody
	return
}
