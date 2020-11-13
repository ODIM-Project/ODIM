package rest

import (
	"log"
	"net/http"
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

func (s *Status) Init() {
	s.Available = "Yes"
	s.TimeStamp = time.Now().Format(time.RFC3339)
	s.Uptime = s.TimeStamp
}

func (s *Status) Refresh() {
	s.TimeStamp = time.Now().Format(time.RFC3339)
}

//TokenMap is used to define the plugin generated tokens
type TokenMap struct {
	Token    string
	LastUsed time.Time
}

var tokenDetails []TokenMap

//TokenValidation validates sent token with the list of plugin generated tokens
func TokenValidation(token string, sessionTimeoutInMinutes float64) bool {
	var flag bool
	flag = false
	for _, v := range tokenDetails {
		if token == v.Token {
			flag = true
			log.Println(time.Since(v.LastUsed).Minutes())
			if time.Since(v.LastUsed).Minutes() > sessionTimeoutInMinutes {
				return flag
			}
		}
	}
	return flag
}

type pluginStatusController struct {
	status       *Status
	pluginConfig *config.PluginConfig
}

func NewPluginStatusController(pc *config.PluginConfig) context.Handler {
	s := &Status{}
	s.Init()

	return pluginStatusController{
		status:       s,
		pluginConfig: pc,
	}.getPluginStatus
}

func (p pluginStatusController) getPluginStatus(ctx iris.Context) {
	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	//Validating the token
	if token != "" {
		flag := TokenValidation(token, p.pluginConfig.SessionTimeoutInMinutes)
		if !flag {
			log.Println("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}
	var resp = PluginStatusResponse{
		Comment: "Plugin Status Response",
		Name:    "Common Redfish Plugin Status",
		Version: p.pluginConfig.FirmwareVersion,
	}
	resp.Status = *p.status
	resp.Status.TimeStamp = time.Now().Format(time.RFC3339)
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(resp)
}
