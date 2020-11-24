package rest

import (
	"net/http"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/kataras/iris/v12/context"
)

type getPluginManagerHandler struct {
	pluginManager redfish.Manager
}

func (m *getPluginManagerHandler) handle(ctx context.Context) {
	requestedManager := ctx.Request().RequestURI
	if requestedManager == m.pluginManager.OdataID {
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(m.pluginManager)
		return
	}
	ctx.StatusCode(http.StatusNotFound)
}

func createPluginManager(pc *config.PluginConfig) redfish.Manager {
	return redfish.Manager{
		OdataContext:    "/ODIM/v1/$metadata#Manager.Manager",
		OdataID:         "/ODIM/v1/Managers/" + pc.RootServiceUUID,
		OdataType:       "#Manager.v1_3_3.Manager",
		Name:            _PLUGIN_NAME,
		ManagerType:     "Service",
		ID:              pc.RootServiceUUID,
		UUID:            pc.RootServiceUUID,
		FirmwareVersion: pc.FirmwareVersion,
		Status: &redfish.ManagerStatus{
			State: "Enabled",
		},
	}
}

func NewGetPluginManagerHandler(pc *config.PluginConfig) context.Handler {
	return (&getPluginManagerHandler{
		pluginManager: createPluginManager(pc),
	}).handle
}
