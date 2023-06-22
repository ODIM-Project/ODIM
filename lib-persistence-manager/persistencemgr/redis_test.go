// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package persistencemgr

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"

	redis "github.com/go-redis/redis"
)

const (
	deleteDataErrMsg    string = "Error while deleting Data: %v\n"
	createDataErrMsg    string = "Error while creating data: %v\n"
	fetchDataErrMsg     string = "Error while fetching data: %v\n"
	readDataErrMsg      string = "Error while reading data: %v\n"
	dataExistErrMsg     string = "Data already exists"
	getDataErrMsg       string = "Error while getting data: %v\n"
	genericErrorMsg     string = "Error: %v\n"
	unmarshalingErrMsg  string = "Error while unmarshaling data : %v\n"
	dataMismatchErrMsg  string = "Mismatch in fetched data"
	dataIncrementErrMsg string = "Error while incrementing data: %v\n"
	dataEntryFailed     string = "Error while making data entry: %v\n"
	dataCleanUpfailed   string = "Error while cleaning up data in DB: %v\n"
	addPluginFailed     string = "Error while adding plugin task to set"
	dataNotFound        string = "data not found"
	mockDBConnection    string = "Error while making mock DB connection:"
	pluginTask          string = "PluginTask:task1"
	pluginTaskIndex     string = "PluginTaskIndex"
	hostIPAddress       string = "10.24.0.1"
	locationURL         string = "https://10.24.1.23/redfish/v1/EventService/Subscriptions/123"
	originResourcesURI  string = "/redfish/v1/Systems/uuid.1"
)

type MockConn struct {
	MockClose   func() error
	MockErr     func() error
	MockDo      func(string, ...interface{}) (interface{}, error)
	MockSend    func(string, ...interface{}) error
	MockFlush   func() error
	MockReceive func() (interface{}, error)
}

func (mc MockConn) Close() error {
	return mc.MockClose()
}

func (mc MockConn) Err() error {
	return mc.MockErr()
}

func (mc MockConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	return mc.MockDo(commandName, args...)
}

func (mc MockConn) Send(commandName string, args ...interface{}) error {
	return mc.MockSend(commandName, args...)
}

func (mc MockConn) Flush() error {
	return mc.MockFlush()
}

func (mc MockConn) Receive() (interface{}, error) {
	return mc.MockReceive()
}

type sample struct {
	Data1 string
	Data2 string
	Data3 string
}

func TestConnection(t *testing.T) {
	config.SetUpMockConfig(t)
	persistConfig, err := GetMockDBConfig()
	if err != nil {
		t.Fatal("Error while initializing config:", err)
	}
	_, errs := persistConfig.Connection()
	if errs != nil {
		t.Fatal("Error while making DB connection:", errs)
	}

}
func TestCreate(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

	if cerr := c.Create("table", "key", "sample"); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

}

func TestCreateInvalidData(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}

	defer func() {
		c.Delete("table", "key")
	}()

	if cerr := c.Create("table", "key", math.Inf(1)); cerr != nil {
		if !(strings.Contains(cerr.Error(), "unsupported")) {
			t.Errorf(dataEntryFailed, cerr.Error())
		}
	}

}

func TestCreateExistingData(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key", data); cerr != nil {
		if errors.DBKeyAlreadyExist != cerr.ErrNo() {
			t.Errorf(dataExistErrMsg)
		}
	}

	data = sample{Data1: "Value4", Data2: "Value5", Data3: "Value6"}
	if cerr := c.Create("table", "key", data); cerr != nil {
		if errors.DBKeyAlreadyExist != cerr.ErrNo() {
			t.Errorf(dataExistErrMsg)
		}
	}
	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

}

func TestRead(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key", data); cerr != nil {
		t.Errorf(genericErrorMsg, cerr.Error())
	}
	got, rerr := c.Read("table", "key")
	if rerr != nil {
		t.Errorf(readDataErrMsg, rerr.Error())
	}
	var res sample
	if jerr := json.Unmarshal([]byte(string(got)), &res); jerr != nil {
		t.Errorf(unmarshalingErrMsg, jerr)
	}

	if res != data {
		t.Errorf(dataMismatchErrMsg)
	}
	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

}

func TestReadNonExistingData(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	_, rerr := c.Read("othertable", "key")

	if rerr != nil {
		if !(strings.Contains(rerr.Error(), "no data with")) {
			t.Errorf(dataEntryFailed, rerr.Error())
		}
	}

}

func TestUpdate(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key", data); cerr != nil {
		t.Errorf(genericErrorMsg, cerr.Error())
	}
	data = sample{Data1: "Value1", Data2: "Value2", Data3: "Value4"}
	uid, uerr := c.Update("table", "key", data)
	if uerr != nil {
		t.Errorf("Error while updating data: %v\n", uerr.Error())
	}
	updateResponse := strings.Split(uid, ":")
	got, rerr := c.Read(updateResponse[0], updateResponse[1])
	if rerr != nil {
		t.Errorf(readDataErrMsg, rerr.Error())
	}
	var res sample
	if jerr := json.Unmarshal([]byte(string(got)), &res); jerr != nil {
		t.Errorf(unmarshalingErrMsg, jerr)
	}

	if res != data {
		t.Errorf(dataMismatchErrMsg)
	}

	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

}

func TestUpdateInvalidData(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key", data); cerr != nil {
		t.Errorf(genericErrorMsg, cerr.Error())
	}

	if _, uerr := c.Update("table", "key", make(chan int, 1)); uerr != nil {
		if !(strings.Contains(uerr.Error(), "unsupported")) {
			t.Errorf(dataEntryFailed, uerr.Error())
		}
	}
	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()
}

func TestGetall(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	data1 := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if errs := c.Create("table", "key1", data1); errs != nil {
		t.Errorf(createDataErrMsg, errs.Error())
	}
	data2 := sample{Data1: "Value4", Data2: "Value5", Data3: "Value6"}
	if errs := c.Create("table", "key2", data2); errs != nil {
		t.Errorf(createDataErrMsg, errs.Error())
	}
	keys, errs := c.GetAllDetails("table")
	if errs != nil {
		t.Errorf(fetchDataErrMsg, errs.Error())
	}

	if len(keys) < 2 {
		t.Errorf("Error in fetching all the keys")
	}
	for _, key := range keys {
		got, rerr := c.Read("table", key)
		if rerr != nil {
			t.Errorf(readDataErrMsg, rerr.Error())
		}
		var res sample
		if err := json.Unmarshal([]byte(string(got)), &res); err != nil {
			t.Errorf("Error during json unmarshal: %v\n", err)
		}

		if res != data1 {
			if res != data2 {
				t.Errorf("Mismatch in data saved")
			}
		}
	}
	defer func() {
		if derr := c.Delete("table", "key1"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
		if derr := c.Delete("table", "key2"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()
}

func TestGetallNonExistingtable(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	keys, errs := c.GetAllDetails("nonExistingtable")
	if errs != nil {
		t.Errorf(fetchDataErrMsg, errs.Error())
	}
	if len(keys) != 0 {
		t.Errorf("Error, fetching data even if table does not exist")
	}
}
func TestDelete(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	data1 := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key1", data1); cerr != nil {
		t.Errorf(createDataErrMsg, cerr.Error())
	}
	if derr := c.Delete("table", "key1"); derr != nil {
		t.Errorf(deleteDataErrMsg, derr.Error())
	}
	if _, rerr := c.Read("table", "key1"); rerr == nil {
		t.Errorf("Error, data still exists post delete operation")
	}

}

func TestDeleteNonExistingKey(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	want := errors.PackError(errors.DBKeyNotFound, "no data with the with key key found")
	derr := c.Delete("table", "key")
	if !reflect.DeepEqual(derr, want) {
		t.Errorf("table should not exist: %v, want: %v\n", derr.Error(), want)
	}
}

func TestCleanUpDB(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	data1 := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if errs := c.Create("table", "key1", data1); errs != nil {
		t.Errorf(createDataErrMsg, errs.Error())
	}
	if errs := c.CleanUpDB(); errs != nil {
		t.Errorf("error while trying to flush db: %v\n", errs.Error())
	}
	keys, errs := c.GetAllDetails("*")
	if errs != nil {
		t.Errorf(fetchDataErrMsg, errs.Error())
	}
	if len(keys) != 0 {
		t.Error("database was not fully cleaned")
	}

}

/*
func TestFilterSearch(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if errs := c.Create("table", "key", data); errs != nil {
		t.Errorf(genericErrorMsg, errs.Error())
	}
	got, errs := c.FilterSearch("table", "key", ".Data3")
	// t.Errorf("HERE: %v, %v",len(string(got.([]uint8))),len("Value3"))
	if errs != nil {
		t.Errorf("error while looking up data: %v", errs.Error())
	}
	if string(got.([]uint8)) != `"Value3"` {
		t.Errorf(dataMismatchErrMsg)
	}
	defer func() {
		if errs := c.Delete("table", "key"); errs != nil {
			t.Errorf(deleteDataErrMsg, errs.Error())
		}
	}()

}
*/
func TestGetAllMatchingDetails(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	data1 := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key1", data1); cerr != nil {
		t.Errorf(createDataErrMsg, cerr.Error())
	}
	data2 := sample{Data1: "Value4", Data2: "Value5", Data3: "Value6"}
	if cerr := c.Create("table", "key2", data2); cerr != nil {
		t.Errorf(createDataErrMsg, cerr.Error())
	}
	keys, _ := c.GetAllMatchingDetails("table", "key")

	if len(keys) < 2 {
		t.Errorf("Error in fetching all the keys")
	}
	for _, key := range keys {
		got, rerr := c.Read("table", key)
		if rerr != nil {
			t.Errorf(readDataErrMsg, rerr.Error())
		}
		var res sample
		if jerr := json.Unmarshal([]byte(string(got)), &res); jerr != nil {
			t.Errorf(unmarshalingErrMsg, jerr)
		}

		if res != data1 {
			if res != data2 {
				t.Errorf("Mismatch in data saved")
			}
		}
	}
	defer func() {
		if derr := c.Delete("table", "key1"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
		if derr := c.Delete("table", "key2"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()
}

func TestGetAllMatchingDetailsNonExistingtable(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	keys, _ := c.GetAllMatchingDetails("nonExistingtable", "key")
	if len(keys) != 0 {
		t.Errorf("Error, fetching data even if table does not exist")
	}
}

func TestTransaction(t *testing.T) {
	const threadCount = 10
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	state := [5]string{"Running", "Completed", "Cancelled", "cancelling", "Pending"}
	data := sample{Data1: "Running", Data2: "Value2", Data3: "Value3"}
	if errs := c.Create("table", "key", data); errs != nil {
		t.Fatal("Error while creating data:", errs)
	}

	updateTransaction := func(key string, state string) error {
		testCallBack := func(ctx context.Context, key string) error {
			got, err := c.Read("table", "key")
			if err != nil {
				t.Fatal("Error while reading data:", err)
				return err
			}
			var res sample
			if emsg := json.Unmarshal([]byte(string(got)), &res); emsg != nil {
				t.Fatal("Error while unmarshaling data:", emsg)
				return emsg
			}
			res.Data1 = state
			if _, err = c.Update("table", "key", res); err != nil {
				t.Fatal("Error while updating data:", err)
				return err
			}
			return nil
		}
		for retries := threadCount / 10; retries > 0; retries-- {
			err := c.Transaction(context.TODO(), "key", testCallBack)
			if err != nil {
				t.Fatal("Error while making a transaction:", err)
			} else {
				return nil
			}
		}
		err := fmt.Errorf("error: updateTransaction reached max retries")
		t.Fatal(err)
		return err
	}

	var wg sync.WaitGroup
	wg.Add(threadCount)
	for i := 0; i < threadCount; i++ {
		status := state[4%(i+1)]
		go func(status string, t *testing.T) {
			defer wg.Done()
			if err := updateTransaction("key", status); err != nil {
				t.Error("error: update transaction failed")

			}
		}(status, t)
	}
	wg.Wait()
	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()
}

func TestGetResourceDetails(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key", data); cerr != nil {
		t.Errorf(genericErrorMsg, cerr.Error())
	}
	got, rerr := c.GetResourceDetails("key")
	if rerr != nil {
		t.Errorf(readDataErrMsg, rerr.Error())
	}
	var res sample
	if jerr := json.Unmarshal([]byte(string(got)), &res); jerr != nil {
		t.Errorf(unmarshalingErrMsg, jerr)
	}

	if res != data {
		t.Errorf(dataMismatchErrMsg)
	}
	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

}

func TestGetResourceDetailsNonExistingData(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	_, rerr := c.GetResourceDetails("keys")
	if rerr != nil {
		if !(strings.Contains(rerr.Error(), "no data with")) {
			t.Errorf(dataEntryFailed, rerr.Error())
		}
	}

}

func TestAddResourceData(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

	if cerr := c.Create("table", "key", "sample"); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

}

func TestAddResourceDataInvalidData(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}

	defer func() {
		c.Delete("table", "key")
	}()

	if cerr := c.Create("table", "key", math.Inf(1)); cerr != nil {
		if !(strings.Contains(cerr.Error(), "unsupported")) {
			t.Errorf(dataEntryFailed, cerr.Error())
		}
	}

}

func TestPing(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Ping(); err != nil {
		t.Errorf("Error while pinging DB: %v\n", err)
	}
}

func TestIndexCreate(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if derr := c.Del("table", "123::sample"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()
	form := map[string]interface{}{"table": 123}
	if cerr := c.CreateIndex(form, "sample"); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}
	form = map[string]interface{}{"table": 2.179699264}
	if cerr := c.CreateIndex(form, "sample"); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}
	form = map[string]interface{}{"table": []float64{2.179699264, 4.543211223}}
	if cerr := c.CreateIndex(form, "sample"); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}
	form = map[string]interface{}{"table": []string{"sample", "sample2"}}
	if cerr := c.CreateIndex(form, "sample"); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

}
func TestCreateIndexInvalidData(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}

	defer func() {
		c.Del("table", "key::sample")
	}()

	if cerr := c.CreateIndex(map[string]interface{}{"table": []string{"key"}}, "sample"); cerr != nil {
		if !(strings.Contains(cerr.Error(), "unsupported")) {
			t.Errorf(dataEntryFailed, cerr.Error())
		}
	}

}

func TestGet(t *testing.T) {
	c, err := MockDBConnection(t)
	defer func() {
		for i := 0; i < 10; {
			f := strconv.Itoa(i)
			key := "value::" + f
			err := c.Del("ProcessorSummary/Model", key)
			if err != nil {
				break
			}

		}
	}()
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	if cerr := c.CreateIndex(map[string]interface{}{"model": "intel"}, "abc-123"); cerr != nil {
		t.Errorf(genericErrorMsg, cerr.Error())
	}
	for i := 0; i < 10; {
		f := strconv.Itoa(i)
		key := "value::" + f
		createErr := c.CreateTaskIndex("ProcessorSummary/Model", 0, key)
		if createErr != nil {
			t.Errorf("Error while trying to create index: %v\n", createErr.Error())
		}
		i++
	}
	got, rerr := c.GetString("model", 0, "*intel*", false)
	if rerr != nil {
		t.Errorf(readDataErrMsg, rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf(dataNotFound)
	}
	got, rerr = c.GetString("ProcessorSummary/Model", 0, "value*", false)
	if rerr != nil {
		t.Errorf(readDataErrMsg, rerr.Error())
	}
	if len(got) != 10000 {
		t.Errorf(dataNotFound)
	}
	defer func() {
		if derr := c.Del("model", "intel::abc-123"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

}

func TestGetTaskList(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	if cerr := c.CreateIndex(map[string]interface{}{"count": 4}, "abc-123"); cerr != nil {
		t.Errorf(genericErrorMsg, cerr.Error())
	}
	got, rerr := c.GetTaskList("count", 0, -1)
	if rerr != nil {
		t.Errorf(readDataErrMsg, rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf(dataNotFound)
	}
	defer func() {
		if derr := c.Del("count", "4::abc-123"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

}
func TestGetRange(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	if cerr := c.CreateIndex(map[string]interface{}{"count": 4}, "abc-123"); cerr != nil {
		t.Errorf(genericErrorMsg, cerr.Error())
	}
	got, rerr := c.GetRange("count", 0, 5, false)
	if rerr != nil {
		t.Errorf(readDataErrMsg, rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf(dataNotFound)
	}
	defer func() {
		if derr := c.Del("count", "4::abc-123"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

}
func TestGetStorageList(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	if cerr := c.CreateIndex(map[string]interface{}{"storage": "[2.179699264 4.543211223]"}, "abc-123"); cerr != nil {
		t.Errorf(genericErrorMsg, cerr.Error())
	}
	got, rerr := c.GetStorageList("storage", 0, 2.179699264, "eq", false)
	if rerr != nil {
		t.Errorf(readDataErrMsg, rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf(dataNotFound)
	}
	got, rerr = c.GetStorageList("storage", 0, 2.179699264, "ge", false)
	if rerr != nil {
		t.Errorf(readDataErrMsg, rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf(dataNotFound)
	}
	got, rerr = c.GetStorageList("storage", 0, 0, "gt", false)
	if rerr != nil {
		t.Errorf(readDataErrMsg, rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf(dataNotFound)
	}
	got, rerr = c.GetStorageList("storage", 0, 4, "lt", false)
	if rerr != nil {
		t.Errorf(readDataErrMsg, rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf(dataNotFound)
	}
	got, rerr = c.GetStorageList("storage", 0, 4, "le", false)
	if rerr != nil {
		t.Errorf(readDataErrMsg, rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf(dataNotFound)
	}

	defer func() {
		if derr := c.Del("storage", "*::abc-123"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

}

func TestCreateEvtSubscriptions(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}
	js := `{"Name":"Subscriptions", "hostip":"10.10.10.10"}`
	defer func() {
		if derr := c.DeleteEvtSubscriptions("subscriptions", js); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

	if cerr := c.CreateEvtSubscriptionIndex("subscriptions", js); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}
}

func TestCreateEvtSubscriptionsExistingData(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	js := `{"Name":"Subscriptions", "hosts":["10.10.10.10"]}`
	defer func() {
		if derr := c.DeleteEvtSubscriptions("subscriptions", js); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

	if cerr := c.CreateEvtSubscriptionIndex("subscriptions", js); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

	if cerr := c.CreateEvtSubscriptionIndex("subscriptions", js); cerr == nil {
		t.Errorf("Error while making data entry")
	}

}

func TestGetEvtSubscriptions(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}
	js := `{"Name":"Subscriptions", "hostip":"10.10.10.10"}`
	defer func() {
		if derr := c.DeleteEvtSubscriptions("subscriptions", js); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

	if cerr := c.CreateEvtSubscriptionIndex("subscriptions", js); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

	subscriptions, gerr := c.GetEvtSubscriptions("subscriptions", "*10.10.10.10*")
	if gerr != nil {
		t.Errorf(getDataErrMsg, gerr.Error())
	}
	if len(subscriptions) < 1 {
		t.Errorf("No data found for the host ip")
	}

	if subscriptions[0] != js {
		t.Errorf("Error while trying to get data")
	}

	subscriptions, _ = c.GetEvtSubscriptions("subscriptions", "*10.10.10.100*")
	if err != nil {
		t.Errorf(getDataErrMsg, err.Error())
	}
	if len(subscriptions) > 0 {
		t.Errorf("data shouldn't be there")
	}

}

func TestDeleteEvtSubscriptions(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}
	js := `{"Name":"Subscriptions", "hostip":"10.10.10.10"}`

	if cerr := c.CreateEvtSubscriptionIndex("subscriptions", js); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

	if derr := c.DeleteEvtSubscriptions("subscriptions", js); derr != nil {
		t.Errorf(deleteDataErrMsg, derr.Error())
	}

}

func TestDeleteEvtSubscriptionsNonexistingData(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}
	js := `{"Name":"Subscriptions", "hostip":"10.10.10.10"}`

	derr := c.DeleteEvtSubscriptions("subscriptions", js)
	if derr == nil {
		t.Errorf("No data found: %v\n", derr.Error())
	}

}

func TestUpdateEvtSubscriptions(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}

	js := `{"Name":"Subscriptions", "SubscriptionID":"12345","hostip":["10.10.10.10"]}`

	if cerr := c.CreateEvtSubscriptionIndex("subscriptions", js); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

	js = `{"Name":"Subscriptions", "SubscriptionID":"12345","hostip":["10.10.10.19"]}`

	if cerr := c.UpdateEvtSubscriptions("subscriptions", "*12345*", js); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

}

func TestCreateDeviceSubscription(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}
	location := locationURL
	originResources := []string{originResourcesURI}
	hostIP := hostIPAddress
	defer func() {
		if derr := c.DeleteDeviceSubscription("DeviceSubscription", hostIP); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()
	if cerr := c.CreateDeviceSubscriptionIndex("DeviceSubscription", hostIP, location, originResources); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

}

func TestCreateDeviceSubscriptionExistingData(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}
	location := locationURL
	originResources := []string{originResourcesURI}
	hostIP := hostIPAddress
	defer func() {
		if derr := c.DeleteDeviceSubscription("DeviceSubscription", hostIP); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()
	if cerr := c.CreateDeviceSubscriptionIndex("DeviceSubscription", hostIP, location, originResources); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

	if cerr := c.CreateDeviceSubscriptionIndex("DeviceSubscription", hostIP, location, originResources); cerr == nil {
		t.Errorf("Data Already Exist: %v\n", cerr.Error())
	}
}

func TestGetDeviceSubscription(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}
	location := locationURL
	originResources := []string{originResourcesURI}
	hostIP := hostIPAddress
	defer func() {
		if derr := c.DeleteDeviceSubscription("DeviceSubscription", hostIP); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()
	if cerr := c.CreateDeviceSubscriptionIndex("DeviceSubscription", hostIP, location, originResources); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

	subscriptions, gerr := c.GetDeviceSubscription("DeviceSubscription", "10.24.0.1*")
	if gerr != nil {
		t.Errorf(getDataErrMsg, gerr.Error())
	}
	if len(subscriptions) < 1 {
		t.Errorf("No data found for the host ip")
	}
	devSub := strings.Split(subscriptions[0], "||")

	if devSub[0] != hostIP && devSub[1] != location {
		t.Errorf("HostIP/Location didn't matched")
	}

	subscriptions, _ = c.GetDeviceSubscription("DeviceSubscription", "10.10.10.100*")
	if err != nil {
		t.Errorf(getDataErrMsg, err.Error())
	}
	if len(subscriptions) > 0 {
		t.Errorf("data shouldn't be there")
	}

}

func TestDeleteDeviceSubscription(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}
	location := locationURL
	originResources := []string{originResourcesURI}
	hostIP := hostIPAddress
	if cerr := c.CreateDeviceSubscriptionIndex("DeviceSubscription", hostIP, location, originResources); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}
	if derr := c.DeleteDeviceSubscription("DeviceSubscription", hostIP); derr != nil {
		t.Errorf(deleteDataErrMsg, derr.Error())
	}

}

func TestDeleteDeviceSubscriptionsNonexistingData(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}
	if derr := c.DeleteDeviceSubscription("DeviceSubscription", "10.24.0.10"); derr == nil {
		t.Errorf("No Data found: %v\n", derr.Error())
	}

}

func TestUpdateDeviceSubscriptions(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}

	location := locationURL
	originResources := []string{originResourcesURI}
	hostIP := hostIPAddress
	defer func() {
		if derr := c.DeleteDeviceSubscription("DeviceSubscription", hostIP); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()
	if cerr := c.CreateDeviceSubscriptionIndex("DeviceSubscription", hostIP, location, originResources); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

	location = "https://10.24.1.23/redfish/v1/EventService/Subscriptions/12345"
	if cerr := c.UpdateDeviceSubscription("DeviceSubscription", hostIP, location, originResources); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

}

type redisExtCallsImpMock struct{}

func (r redisExtCallsImpMock) newSentinelClient(opt *redis.Options) *redis.SentinelClient {
	return newSentinelClientMock(opt)
}
func newSentinelClientMock(opt *redis.Options) *redis.SentinelClient {
	strSlice := strings.Split(opt.Addr, ":")
	sentinelHost := strSlice[0]
	sentinelPort := strSlice[1]
	if sentinelHost == "ValidHost" && sentinelPort == "ValidSentinelPort" {
		return &redis.SentinelClient{}
	}
	return nil
}
func (r redisExtCallsImpMock) getMasterAddrByName(masterSet string, snlClient *redis.SentinelClient) []string {
	return getMasterAddbyNameMock(masterSet, snlClient)
}

func getMasterAddbyNameMock(masterSet string, snlClient *redis.SentinelClient) []string {
	if masterSet == "ValidMasterSet" && snlClient != nil {
		return []string{"ValidMasterIP", "ValidMasterPort"}
	}
	return []string{"", ""}
}

func TestGetCurrentMasterHostPort(t *testing.T) {
	redisExtCalls = redisExtCallsImpMock{}
	type args struct {
		dbConfig *Config
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "Positive Case: All is well, valid sentinel Host",
			args: args{
				dbConfig: &Config{
					Host:         "ValidHost",
					SentinelPort: "ValidSentinelPort",
					MasterSet:    "ValidMasterSet",
				},
			},
			want:  "ValidMasterIP",
			want1: "ValidMasterPort",
		},
		{
			name: "Negative Case: Invalid sentinel Host",
			args: args{
				dbConfig: &Config{
					Host:         "InvalidHost",
					SentinelPort: "ValidSentinelPort",
					MasterSet:    "ValidMasterSet",
				},
			},
			want:  "",
			want1: "",
		},
		{
			name: "Negative Case: Invalid sentinel Port",
			args: args{
				dbConfig: &Config{
					Host:         "ValidHost",
					SentinelPort: "InvalidSentinelPort",
					MasterSet:    "ValidMasterSet",
				},
			},
			want:  "",
			want1: "",
		},
		{
			name: "Negative Case: Invalid MasterSet",
			args: args{
				dbConfig: &Config{
					Host:         "ValidHost",
					SentinelPort: "ValidSentinelPort",
					MasterSet:    "InvalidMasterSet",
				},
			},
			want:  "",
			want1: "",
		},
		{
			name: "Negative Case: empty sentinel Host",
			args: args{
				dbConfig: &Config{
					Host:         "",
					SentinelPort: "ValidSentinelPort",
					MasterSet:    "ValidMasterSet",
				},
			},
			want:  "",
			want1: "",
		},
		{
			name: "Negative Case: empty MasterSet",
			args: args{
				dbConfig: &Config{
					Host:         "ValidHost",
					SentinelPort: "ValidSentinelPort",
					MasterSet:    "",
				},
			},
			want:  "",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, _ := GetCurrentMasterHostPort(tt.args.dbConfig)
			if got != tt.want {
				t.Errorf("GetCurrentMasterHostPort() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetCurrentMasterHostPort() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
func TestGetDBConnection(t *testing.T) {
	GetMockDBConfig()
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	redisExtCalls = redisExtCallsImpMock{}
	type args struct {
		dbFlag DbType
	}
	tests := []struct {
		name  string
		args  args
		want  *ConnPool
		want1 *errors.Error
	}{
		{
			name: "Positive case: All is well, inmemory db type",
			args: args{
				dbFlag: InMemory,
			},
			want:  &ConnPool{ReadPool: c.ReadPool},
			want1: nil,
		},
		{
			name: "Positive case: All is well, OnDisk db type",
			args: args{
				dbFlag: OnDisk,
			},
			want:  &ConnPool{ReadPool: c.ReadPool},
			want1: nil,
		},
		{
			name: "Negative case: invalid db type",
			args: args{
				dbFlag: 3,
			},
			want:  nil,
			want1: errors.PackError(errors.UndefinedErrorType, "error invalid db type selection"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetDBConnection(tt.args.dbFlag)
			if (got != nil) != (tt.want != nil) {
				t.Errorf("GetDBConnection() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetDBConnection() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetDBConnectionHAEnabled(t *testing.T) {
	GetMockDBConfig()
	// Enableing HA
	config.Data.DBConf.RedisHAEnabled = true

	inMemDBConnPool = &ConnPool{
		ReadPool: &redis.Client{},
		MasterIP: "NotValid",
	}
	onDiskDBConnPool = &ConnPool{
		ReadPool: &redis.Client{},
		MasterIP: "NotValid",
	}
	redisExtCalls = redisExtCallsImpMock{}
	type args struct {
		dbFlag DbType
	}
	tests := []struct {
		name  string
		args  args
		want  *ConnPool
		want1 *errors.Error
	}{
		{
			name: "Positive case: All is well, inmemory db type",
			args: args{
				dbFlag: InMemory,
			},
			want:  &ConnPool{},
			want1: nil,
		},
		{
			name: "Positive case: All is well, OnDisk db type",
			args: args{
				dbFlag: OnDisk,
			},
			want:  &ConnPool{},
			want1: nil,
		},
		{
			name: "Negative case: invalid db type",
			args: args{
				dbFlag: 3,
			},
			want:  nil,
			want1: errors.PackError(errors.UndefinedErrorType, "error invalid db type selection"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetDBConnection(tt.args.dbFlag)
			if (got != nil) != (tt.want != nil) {
				t.Errorf("GetDBConnection() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetDBConnection() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestIncr(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}

	got, rerr := c.Incr("table", "key")
	if rerr != nil {
		t.Errorf(dataIncrementErrMsg, rerr.Error())
	}
	if got != 1 {
		t.Errorf(dataMismatchErrMsg)
	}

	got, rerr = c.Incr("table", "key")
	if rerr != nil {
		t.Errorf(dataIncrementErrMsg, rerr.Error())
	}
	if got != 2 {
		t.Errorf(dataMismatchErrMsg)
	}

	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

}

func TestDecr(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	_, rerr := c.Incr("table", "key")
	if rerr != nil {
		t.Errorf(dataIncrementErrMsg, rerr.Error())
	}
	_, rerr = c.Incr("table", "key")
	if rerr != nil {
		t.Errorf(dataIncrementErrMsg, rerr.Error())
	}

	got, rerr := c.Decr("table", "key")
	if rerr != nil {
		t.Errorf(dataIncrementErrMsg, rerr.Error())
	}
	if got != 1 {
		t.Errorf(dataMismatchErrMsg)
	}

	got, rerr = c.Decr("table", "key")
	if rerr != nil {
		t.Errorf(dataIncrementErrMsg, rerr.Error())
	}
	if got != 0 {
		t.Errorf(dataMismatchErrMsg)
	}

	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

}

func TestSetExpire(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

	if cerr := c.SetExpire("table", "key", "sample", 1); cerr != nil {
		t.Errorf(dataEntryFailed, cerr.Error())
	}

}

func TestSetExpireInvalidData(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}

	defer func() {
		c.Delete("table", "key")
	}()

	if cerr := c.SetExpire("table", "key", math.Inf(1), 1); cerr != nil {
		if !(strings.Contains(cerr.Error(), "unsupported")) {
			t.Errorf(dataEntryFailed, cerr.Error())
		}
	}

}

func TestSetExpireExistingData(t *testing.T) {

	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.SetExpire("table", "key", data, 1); cerr != nil {
		if errors.DBKeyAlreadyExist != cerr.ErrNo() {
			t.Errorf(dataExistErrMsg)
		}
	}

	data = sample{Data1: "Value4", Data2: "Value5", Data3: "Value6"}
	if cerr := c.SetExpire("table", "key", data, 1); cerr != nil {
		if errors.DBKeyAlreadyExist != cerr.ErrNo() {
			t.Errorf(dataExistErrMsg)
		}
	}
	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

}

func TestTTL(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}

	rerr := c.SetExpire("table", "key", "", 2)
	if rerr != nil {
		t.Errorf("Error while setting data: %v\n", rerr.Error())
	}
	time, rerr := c.TTL("table", "key")
	if rerr != nil {
		t.Errorf("Error while TTL jey: %v\n", rerr.Error())
	}
	if time < 0 {
		t.Errorf("Time should not be elapsed")
	}
	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf(deleteDataErrMsg, derr.Error())
		}
	}()

}

func TestConnPoolGetWriteConnection(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	tests := []struct {
		name    string
		p       *ConnPool
		wantErr bool
	}{
		{
			name:    "get write connection to db",
			p:       c,
			wantErr: false,
		},
		{
			name: "fail while getting write connection to DB if write pool is nil",
			p: &ConnPool{
				ReadPool: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.p.GetWriteConnection()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWriteConnection() : error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnUpdateTransaction(t *testing.T) {
	c, err := MockDBWriteConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}
	type args struct {
		data map[string]interface{}
	}
	tests := []struct {
		name    string
		c       *Conn
		args    args
		wantErr bool
	}{
		{
			name: "db update operation using pipelined transaction",
			c:    c,
			args: args{
				data: map[string]interface{}{
					"TASK:1": "progress",
					"TASK:2": "completed",
				},
			},
			wantErr: false,
		},
		{
			name: "failure while db update operation",
			c: &Conn{
				WriteConn: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UpdateTransaction(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTransaction() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetSortedMapKeys(t *testing.T) {
	type args struct {
		m interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "success case 1",
			args: args{
				m: map[string]interface{}{
					"TASK:1": "task1",
					"TASK:2": "task1",
					"TASK:3": "task1",
				},
			},
			want: []string{"TASK:1", "TASK:2", "TASK:3"},
		},
		{
			name: "success case 2",
			args: args{
				m: map[string]int64{
					"TASK:1": 1,
					"TASK:2": 2,
					"TASK:3": 3,
				},
			},
			want: []string{"TASK:1", "TASK:2", "TASK:3"},
		},
		{
			name: "invalid case",
			args: args{
				m: map[string]string{
					"TASK:1": "task1",
					"TASK:2": "task1",
					"TASK:3": "task1",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSortedMapKeys(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSortedMapKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRetriable(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success case 1",
			args: args{
				err: fmt.Errorf("EOF"),
			},
			want: true,
		},
		{
			name: "success case 2",
			args: args{
				err: fmt.Errorf("error LOADING redis"),
			},
			want: true,
		},
		{
			name: "success case 3",
			args: args{
				err: fmt.Errorf("ERR max number of clients reached"),
			},
			want: true,
		},
		{
			name: "failure case",
			args: args{
				err: fmt.Errorf("unexpected error"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetriable(tt.args.err); got != tt.want {
				t.Errorf("IsRetriable() = %v, want %v", got, tt.want)
			}
		})
	}
}

type timeOutError struct {
	error
}

func (e timeOutError) Timeout() bool {
	return true
}

func (e timeOutError) Temporary() bool {
	return true
}

func (e timeOutError) Error() string {
	return ""
}

func TestIsTimeOutError(t *testing.T) {

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success case",
			args: args{
				err: &timeOutError{
					error: fmt.Errorf("timeout error"),
				},
			},
			want: true,
		},
		{
			name: "failure case",
			args: args{
				err: fmt.Errorf("db error"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTimeOutError(tt.args.err); got != tt.want {
				t.Errorf("isTimeOutError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnPoolAddMemberToSet(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}

	defer func() {
		if derr := c.CleanUpDB(); derr != nil {
			t.Errorf(dataCleanUpfailed, derr.Error())
		}
	}()

	type args struct {
		key    string
		member string
	}
	tests := []struct {
		name    string
		p       *ConnPool
		args    args
		wantErr bool
	}{
		{
			name: "add members to a set",
			p:    c,
			args: args{
				key:    pluginTaskIndex,
				member: pluginTask,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.AddMemberToSet(tt.args.key, tt.args.member); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTransaction() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnPoolGetAllMembersInSet(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}

	err = c.AddMemberToSet(pluginTaskIndex, pluginTask)
	if err != nil {
		t.Fatal(addPluginFailed, err)
	}

	defer func() {
		if derr := c.CleanUpDB(); derr != nil {
			t.Errorf(dataCleanUpfailed, derr.Error())
		}
	}()

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		p       *ConnPool
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "get all members from plugin task",
			p:    c,
			args: args{
				key: pluginTaskIndex,
			},
			want:    []string{pluginTask},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetAllMembersInSet(tt.args.key)
			fmt.Println("got", got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConnPool.GetAllMembersInSet() got = %v, want %v", got, tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnPool.GetAllMembersInSet() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnPoolRemoveMemberFromSet(t *testing.T) {
	c, err := MockDBConnection(t)
	if err != nil {
		t.Fatal(mockDBConnection, err)
	}

	err = c.AddMemberToSet(pluginTaskIndex, pluginTask)
	if err != nil {
		t.Fatal(addPluginFailed, err)
	}

	err = c.AddMemberToSet(pluginTaskIndex, "PluginTask:task2")
	if err != nil {
		t.Fatal(addPluginFailed, err)
	}

	defer func() {
		if derr := c.CleanUpDB(); derr != nil {
			t.Errorf(dataCleanUpfailed, derr.Error())
		}
	}()

	type args struct {
		key    string
		member string
	}
	tests := []struct {
		name    string
		p       *ConnPool
		args    args
		wantErr bool
	}{
		{
			name: "remove member from set",
			p:    c,
			args: args{
				key:    pluginTaskIndex,
				member: pluginTask,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.RemoveMemberFromSet(tt.args.key, tt.args.member); (err != nil) != tt.wantErr {
				t.Errorf("ConnPool.RemoveMemberFromSet() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
