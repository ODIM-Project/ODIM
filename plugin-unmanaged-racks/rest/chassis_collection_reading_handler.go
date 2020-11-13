package rest

import (
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12/context"
)

func newGetChassisCollectionHandler(cm *db.ConnectionManager) context.Handler {
	return (&getChassisCollectionHandler{cm}).handle
}

type getChassisCollectionHandler struct {
	cm *db.ConnectionManager
}

func (c *getChassisCollectionHandler) handle(ctx context.Context) {
	conn := c.cm.GetConnection()
	defer db.NewConnectionCloser(&conn)

	searchKey := db.CreateKey("Chassis")
	keys, err := redis.Strings(conn.Do("KEYS", searchKey.WithWildcard()))
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
	}

	collection := createChassisCollection()
	for _, k := range keys {
		collection.Members = append(
			collection.Members,
			redfish.Link{
				Oid: strings.TrimPrefix(k, searchKey.Prefix()),
			},
		)
		collection.MembersCount++
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(&collection)
}

func createChassisCollection() redfish.Collection {
	return redfish.NewCollection("/ODIM/v1/Chassis", "#ManagerCollection.ManagerCollection")
}
