package chassis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/stretchr/testify/assert"
)

func mockAddChassis(body []byte, table, key string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err := connPool.Create(table, key, string(body)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", table, err.Error())
	}
	return nil
}

func TestHandle(t *testing.T) {
	create := Create{}
	req := chassisproto.CreateChassisRequest{}
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	response := create.Handle(&req)
	assert.Equal(t, http.StatusInternalServerError, int(response.StatusCode), "Request with empty data , Status code should be StatusInternalServerError")
	req = chassisproto.CreateChassisRequest{
		RequestBody: []byte(`{
			"ChassisType": "RackGroup",
			"Description": "My RackGroup",
			"Links": {
			  "ManagedBy": [
				{
				  "@odata.id": "/redfish/v1/Managers/1"
				}
			  ]
			},
			"Name": "RG5"
		  }`),
	}
	GetDbConnectFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	response = create.Handle(&req)
	assert.Equal(t, http.StatusInternalServerError, int(response.StatusCode), "Can not acquire database connection")

	GetDbConnectFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	response = create.Handle(&req)
	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Managers/1"})

	err := mockAddChassis(reqData, "Managers", "/redfish/v1/Managers/1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	assert.Equal(t, http.StatusBadRequest, int(response.StatusCode), "error occured during database access")

}
