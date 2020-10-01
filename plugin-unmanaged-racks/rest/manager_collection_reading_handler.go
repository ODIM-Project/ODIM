package rest

import (
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"
	"github.com/kataras/iris/v12/context"
	"net/http"
)

type getManagersHandler struct {
	managersCollection redfish.Collection
	pc                 *config.PluginConfig
}

func (m *getManagersHandler) handle(ctx context.Context) {
	ctx.JSON(m.managersCollection)
	ctx.StatusCode(http.StatusOK)
}

func NewGetManagersHandler(pc *config.PluginConfig) context.Handler {
	collection := redfish.NewCollection("/ODIM/v1/Managers", "#ManagerCollection.ManagerCollection", redfish.Link{Oid: "/ODIM/v1/Managers/" + pc.RootServiceUUID})
	return (&getManagersHandler{collection, pc}).handle
}
