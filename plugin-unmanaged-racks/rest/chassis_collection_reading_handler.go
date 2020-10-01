package rest

import (
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"
	"github.com/kataras/iris/v12/context"
	"net/http"
)

func NewGetChassisCollectionHandler(cm *db.ConnectionManager) context.Handler {
	return (&getChassisCollectionHandler{cm}).handle
}

type getChassisCollectionHandler struct {
	cm *db.ConnectionManager
}

func (c *getChassisCollectionHandler) handle(ctx context.Context) {
	keys, err := c.cm.GetAllKeys("Chassis")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		return
	}

	collection := createChassisCollection()
	for _, k := range keys {
		collection.Members = append(collection.Members, redfish.Link{Oid: k})
		collection.MembersCount++
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(&collection)
}

func createChassisCollection() redfish.Collection {
	return redfish.NewCollection("/ODIM/v1/Chassis", "#ManagerCollection.ManagerCollection")
}
