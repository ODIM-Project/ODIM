package telemetry

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tmodel"
)

func MockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) (response.RPC, error) {
	if sessionToken == "InvalidToken" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error while trying to authenticate session", nil, nil), nil
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil), nil
}

func MockContactClient(ctx context.Context, url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
	if url == "https://localhost:9091/ODIM/v1/Sessions" {
		body := `{"Token": "12345"}`
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header: http.Header{
				"X-Auth-Token": []string{"12345"},
			},
		}, nil
	} else if url == "https://localhost:9092/ODIM/v1/Sessions" {
		body := `{"Token": ""}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1" && token == "12345" {
		body := `{"@odata.id":"/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1","@odata.type":"#MetricReport.v1_0_0.MetricReport","Id":"CPUUtilCustom1","Name":"Metric report of CPU Utilization for 10 minutes with sensing interval of 20 seconds.","MetricReportDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom1"},"MetricValues":[{"MetricDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"},"MetricId":"CPUUtil","MetricValue":"0","Timestamp":"2021-06-16T07:59:43Z"},{"MetricDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"},"MetricId":"CPUUtil","MetricValue":"0","Timestamp":"2021-06-16T08:00:04Z"}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func MockGetResource(table, key string, dbType common.DbType) (string, *errors.Error) {
	if key == "error" {
		return "", &errors.Error{}
	}
	return "body", nil
}

func MockGetAllKeysFromTable(table string, dbType common.DbType) ([]string, error) {
	if table == "Plugin" {
		return []string{"ILO", "GRF"}, nil
	}
	return []string{"/redfish/v1/TelemetryService/Triggers/uuid.1"}, nil
}

func GetEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func MockGetPluginData(pluginID string) (tmodel.Plugin, *errors.Error) {
	var t *testing.T
	password := GetEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := tmodel.Plugin{
		IP:                "localhost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		ID:                pluginID,
		PreferredAuthType: "XAuthToken",
		PluginType:        "Compute",
	}
	return plugin, nil
}

func MockGetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		External: External{
			Auth:          MockIsAuthorized,
			ContactClient: MockContactClient,
			GetPluginData: MockGetPluginData,
		},
		DB: DB{
			GetAllKeysFromTable: MockGetAllKeysFromTable,
			GetResource:         MockGetResource,
		},
	}
}
