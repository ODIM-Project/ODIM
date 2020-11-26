/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rest

import (
	"time"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type PluginStatusResponse struct {
	Comment         string          `json:"_comment"`
	Name            string          `json:"Name"`
	Version         string          `json:"Version"`
	Status          Status          `json:"Status"`
	EventMessageBus EventMessageBus `json:"EventMessageBus"`
}

type Status struct {
	Available string `json:"Available"`
	Uptime    string `json:"Uptime"`
	TimeStamp string `json:"TimeStamp"`
}

type EventMessageBus struct {
	EmbType  string     `json:"EmbType"`
	EmbQueue []EmbQueue `json:"EmbQueue"`
}

type EmbQueue struct {
	QueueName string `json:"EmbQueueName"`
	QueueDesc string `json:"EmbQueueDesc"`
}

func (s *Status) Init() *Status {
	s.Refresh()
	s.Available = "Yes"
	s.Uptime = s.TimeStamp
	return s
}

func (s *Status) Refresh() {
	s.TimeStamp = time.Now().Format(time.RFC3339)
}

type pluginStatusController struct {
	status       *Status
	pluginConfig *config.PluginConfig
}

func newPluginStatusController(pc *config.PluginConfig) context.Handler {
	return pluginStatusController{status: new(Status).Init(), pluginConfig: pc}.getPluginStatus
}

func (p pluginStatusController) getPluginStatus(ctx iris.Context) {
	p.status.Refresh()
	var resp = PluginStatusResponse{
		Comment: "Unmanaged Racks Plugin",
		Name:    _PLUGIN_NAME,
		Version: p.pluginConfig.FirmwareVersion,
		Status:  *p.status,
	}
	ctx.JSON(resp)
}
