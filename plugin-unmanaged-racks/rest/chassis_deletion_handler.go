package rest

import (
	stdCtx "context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12/context"
)

func NewChassisDeletionHandler(cm *db.ConnectionManager) context.Handler {
	return (&chassisDeletionHandler{cm}).handle
}

type chassisDeletionHandler struct {
	connectionManager *db.ConnectionManager
}

func (c *chassisDeletionHandler) handle(ctx context.Context) {
	requestedChassis := ctx.Request().RequestURI
	requestedChassisKey := db.CreateKey("Chassis", requestedChassis)

	bytes, err := c.connectionManager.DAO().Get(stdCtx.TODO(), requestedChassisKey.String()).Bytes()
	if err != nil && err == redis.Nil {
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(redfish.NewResourceNotFoundMsg("Chassis", requestedChassis, ""))
		return
	}
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	chassisToBeDeleted := new(redfish.Chassis)
	err = json.Unmarshal(bytes, chassisToBeDeleted)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	validator := c.createValidator(chassisToBeDeleted)
	r := validator.Validate()
	if r.HasErrors() {
		ctx.StatusCode(http.StatusConflict)
		ctx.JSON(r.Error())
		return
	}

	switch chassisToBeDeleted.ChassisType {
	case "RackGroup":
		ctx := stdCtx.TODO()
		err = c.connectionManager.DAO().Watch(ctx, func(tx *redis.Tx) error {
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				_, err = pipe.Del(ctx, requestedChassisKey.String()).Result()
				return err
			})
			return err
		}, requestedChassis)

	case "Rack":
		ctx := stdCtx.TODO()

		transactional := func(tx *redis.Tx) error {
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				if _, err = pipe.Del(ctx, requestedChassisKey.String()).Result(); err != nil {
					return fmt.Errorf("del: %s error: %w", requestedChassisKey, err)
				}

				if _, err = pipe.Del(ctx, db.CreateContainedInKey("Chassis", requestedChassis).String()).Result(); err != nil {
					return fmt.Errorf("del: %s error: %w", db.CreateContainedInKey("Chassis", requestedChassis), err)
				}

				parentContainsId := db.CreateContainsKey("Chassis", chassisToBeDeleted.Links.ContainedBy[0].Oid).String()
				_, err = pipe.SRem(ctx, parentContainsId, requestedChassis).Result()
				return err
			})
			return err
		}

		err = c.connectionManager.DAO().
			Watch(
				ctx,
				transactional,
				requestedChassis, db.CreateContainedInKey(requestedChassisKey.String()).String(), db.CreateContainsKey(chassisToBeDeleted.Links.ContainedBy[0].Oid).String(),
			)

		mem, e := c.connectionManager.DAO().SMembers(ctx, db.CreateContainsKey("Chassis", "CONTAINS:Chassis:/ODIM/v1/Chassis/1f5780bc-1c86-52cb-b2ed-ba67cd2345f7").String()).Result()
		fmt.Println(e)
		fmt.Println(mem)
	}

	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	ctx.StatusCode(http.StatusOK)
}

func (c *chassisDeletionHandler) createValidator(chassis *redfish.Chassis) *redfish.CompositeValidator {
	return &redfish.CompositeValidator{
		redfish.Validator{
			ValidationRule: func() bool {
				return !strings.Contains(strings.Join([]string{"", "RackGroup", "Rack"}, "#"), chassis.ChassisType)
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewPropertyValueNotInListMsg(chassis.ChassisType, "ChassisType", "supported ChassisTypes are: RackGroup|Rack")
			},
		},
		redfish.Validator{
			ValidationRule: func() bool {
				hasChildren, err := c.connectionManager.DAO().Exists(stdCtx.TODO(), db.CreateContainsKey("Chassis", chassis.Oid).String()).Result()
				return err != nil || hasChildren == 1
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewResourceInUseMsg("there are existing elements(Links.Contains) under requested chassis")
			},
		},
	}
}
