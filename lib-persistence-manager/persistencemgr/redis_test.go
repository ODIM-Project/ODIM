//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.
package persistencemgr

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	redisSentinel "github.com/go-redis/redis"
	"github.com/gomodule/redigo/redis"
)

type sample struct {
	Data1 string
	Data2 string
	Data3 string
}

func TestConnection(t *testing.T) {
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

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

	if cerr := c.Create("table", "key", "sample"); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

}

func TestCreate_invalidData(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}

	defer func() {
		c.Delete("table", "key")
	}()

	if cerr := c.Create("table", "key", math.Inf(1)); cerr != nil {
		if !(strings.Contains(cerr.Error(), "unsupported")) {
			t.Errorf("Error while making data entry: %v\n", cerr.Error())
		}
	}

}

func TestCreate_existingData(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB coonection:", err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key", data); cerr != nil {
		if errors.DBKeyAlreadyExist != cerr.ErrNo() {
			t.Errorf("Data already exists")
		}
	}

	data = sample{Data1: "Value4", Data2: "Value5", Data3: "Value6"}
	if cerr := c.Create("table", "key", data); cerr != nil {
		if errors.DBKeyAlreadyExist != cerr.ErrNo() {
			t.Errorf("Data already exists")
		}
	}
	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

}

func TestRead(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key", data); cerr != nil {
		t.Errorf("Error: %v\n", cerr.Error())
	}
	got, rerr := c.Read("table", "key")
	if rerr != nil {
		t.Errorf("Error while reading data: %v\n", rerr.Error())
	}
	var res sample
	if jerr := json.Unmarshal([]byte(string(got)), &res); jerr != nil {
		t.Errorf("Error while unmarshaling data : %v\n", jerr)
	}

	if res != data {
		t.Errorf("Mismatch in fetched data")
	}
	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

}

func TestRead_nonExistingData(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	_, rerr := c.Read("othertable", "key")

	if rerr != nil {
		if !(strings.Contains(rerr.Error(), "no data with")) {
			t.Errorf("Error while making data entry: %v\n", rerr.Error())
		}
	}

}

func TestUpdate(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key", data); cerr != nil {
		t.Errorf("Error: %v\n", cerr.Error())
	}
	data = sample{Data1: "Value1", Data2: "Value2", Data3: "Value4"}
	uid, uerr := c.Update("table", "key", data)
	if uerr != nil {
		t.Errorf("Error while updating data: %v\n", uerr.Error())
	}
	updateResponse := strings.Split(uid, ":")
	got, rerr := c.Read(updateResponse[0], updateResponse[1])
	if rerr != nil {
		t.Errorf("Error while read data: %v\n", rerr.Error())
	}
	var res sample
	if jerr := json.Unmarshal([]byte(string(got)), &res); jerr != nil {
		t.Errorf("Error while unmarshaling data : %v\n", jerr)
	}

	if res != data {
		t.Errorf("Mismatch in fetched data")
	}

	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

}

func TestUpdate_invalidData(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key", data); cerr != nil {
		t.Errorf("Error: %v\n", cerr.Error())
	}

	if _, uerr := c.Update("table", "key", make(chan int, 1)); uerr != nil {
		if !(strings.Contains(uerr.Error(), "unsupported")) {
			t.Errorf("Error while making data entry: %v\n", uerr.Error())
		}
	}
	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()
}

func TestUpdate_nonExistingData(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}

	data := sample{Data1: "Value5", Data2: "Value6", Data3: "Value4"}

	_, uerr := c.Update("table", "nonExistingKey", data)

	if uerr.ErrNo() != errors.DBKeyNotFound {
		t.Errorf("Error while updating data: %v\n", uerr.Error())
	}
}

func TestGetall(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	data1 := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if errs := c.Create("table", "key1", data1); errs != nil {
		t.Errorf("Error while creating data: %v\n", errs.Error())
	}
	data2 := sample{Data1: "Value4", Data2: "Value5", Data3: "Value6"}
	if errs := c.Create("table", "key2", data2); errs != nil {
		t.Errorf("Error while creating data: %v\n", errs.Error())
	}
	keys, errs := c.GetAllDetails("table")
	if errs != nil {
		t.Errorf("Error while fetching data: %v\n", errs.Error())
	}

	if len(keys) < 2 {
		t.Errorf("Error in fetching all the keys")
	}
	for _, key := range keys {
		got, rerr := c.Read("table", key)
		if rerr != nil {
			t.Errorf("Error while read data: %v\n", rerr.Error())
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
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
		if derr := c.Delete("table", "key2"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()
}

func TestGetall_nonExistingtable(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	keys, errs := c.GetAllDetails("nonExistingtable")
	if errs != nil {
		t.Errorf("Error while fetching data: %v\n", errs.Error())
	}
	if len(keys) != 0 {
		t.Errorf("Error, fetching data even if table does not exist")
	}
}
func TestDelete(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	data1 := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key1", data1); cerr != nil {
		t.Errorf("Error while creating data: %v\n", cerr.Error())
	}
	if derr := c.Delete("table", "key1"); derr != nil {
		t.Errorf("Error while deleting Data: %v\n", derr.Error())
	}
	if _, rerr := c.Read("table", "key1"); rerr == nil {
		t.Errorf("Error, data still exists post delete operation")
	}

}

func TestDelete_nonExistingKey(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}

	derr := c.Delete("table", "key")
	if derr.ErrNo() != errors.DBKeyNotFound {
		t.Errorf("table should not exist: %v\n", derr.Error())
	}
}

func TestCleanUpDB(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	data1 := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if errs := c.Create("table", "key1", data1); errs != nil {
		t.Errorf("Error while creating data: %v\n", errs.Error())
	}
	if errs := c.CleanUpDB(); errs != nil {
		t.Errorf("error while trying to flush db: %v\n", errs.Error())
	}
	keys, errs := c.GetAllDetails("*")
	if errs != nil {
		t.Errorf("Error while fetching data: %v\n", errs.Error())
	}
	if len(keys) != 0 {
		t.Error("database was not fully cleaned")
	}

}

/*
func TestFilterSearch(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if errs := c.Create("table", "key", data); errs != nil {
		t.Errorf("Error: %v\n", errs.Error())
	}
	got, errs := c.FilterSearch("table", "key", ".Data3")
	// t.Errorf("HERE: %v, %v",len(string(got.([]uint8))),len("Value3"))
	if errs != nil {
		t.Errorf("error while looking up data: %v", errs.Error())
	}
	if string(got.([]uint8)) != `"Value3"` {
		t.Errorf("Mismatch in fetched data")
	}
	defer func() {
		if errs := c.Delete("table", "key"); errs != nil {
			t.Errorf("Error while deleting Data: %v\n", errs.Error())
		}
	}()

}
*/
func TestGetAllMatchingDetails(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	data1 := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key1", data1); cerr != nil {
		t.Errorf("Error while creating data: %v\n", cerr.Error())
	}
	data2 := sample{Data1: "Value4", Data2: "Value5", Data3: "Value6"}
	if cerr := c.Create("table", "key2", data2); cerr != nil {
		t.Errorf("Error while creating data: %v\n", cerr.Error())
	}
	keys, err := c.GetAllMatchingDetails("table", "key")

	if len(keys) < 2 {
		t.Errorf("Error in fetching all the keys")
	}
	for _, key := range keys {
		got, rerr := c.Read("table", key)
		if rerr != nil {
			t.Errorf("Error while read data: %v\n", rerr.Error())
		}
		var res sample
		if jerr := json.Unmarshal([]byte(string(got)), &res); jerr != nil {
			t.Errorf("Error while unmarshaling data : %v\n", jerr)
		}

		if res != data1 {
			if res != data2 {
				t.Errorf("Mismatch in data saved")
			}
		}
	}
	defer func() {
		if derr := c.Delete("table", "key1"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
		if derr := c.Delete("table", "key2"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()
}

func TestGetAllMatchingDetails_nonExistingtable(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	keys, err := c.GetAllMatchingDetails("nonExistingtable", "key")
	if len(keys) != 0 {
		t.Errorf("Error, fetching data even if table does not exist")
	}
}

func TestTransaction(t *testing.T) {
	const threadCount = 10
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	state := [5]string{"Running", "Completed", "Cancelled", "cancelling", "Pending"}
	data := sample{Data1: "Running", Data2: "Value2", Data3: "Value3"}
	if errs := c.Create("table", "key", data); errs != nil {
		t.Fatal("Error while creating data:", errs)
	}

	updateTransaction := func(key string, state string) error {
		testCallBack := func(key string) error {
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
			err := c.Transaction("key", testCallBack)
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
				t.Fatal("error: update transaction failed")

			}
		}(status, t)
	}
	wg.Wait()
	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()
}

func TestGetResourceDetails(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	data := sample{Data1: "Value1", Data2: "Value2", Data3: "Value3"}
	if cerr := c.Create("table", "key", data); cerr != nil {
		t.Errorf("Error: %v\n", cerr.Error())
	}
	got, rerr := c.GetResourceDetails("key")
	if rerr != nil {
		t.Errorf("Error while reading data: %v\n", rerr.Error())
	}
	var res sample
	if jerr := json.Unmarshal([]byte(string(got)), &res); jerr != nil {
		t.Errorf("Error while unmarshaling data : %v\n", jerr)
	}

	if res != data {
		t.Errorf("Mismatch in fetched data")
	}
	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

}

func TestGetResourceDetails_nonExistingData(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	_, rerr := c.GetResourceDetails("key")

	if rerr != nil {
		if !(strings.Contains(rerr.Error(), "no data with")) {
			t.Errorf("Error while making data entry: %v\n", rerr.Error())
		}
	}

}

func TestAddResourceData(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if derr := c.Delete("table", "key"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

	if cerr := c.Create("table", "key", "sample"); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

}

func TestAddResourceData_invalidData(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}

	defer func() {
		c.Delete("table", "key")
	}()

	if cerr := c.Create("table", "key", math.Inf(1)); cerr != nil {
		if !(strings.Contains(cerr.Error(), "unsupported")) {
			t.Errorf("Error while making data entry: %v\n", cerr.Error())
		}
	}

}

func TestPing(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Ping(); err != nil {
		t.Errorf("Error while pinging DB: %v\n", err)
	}
}

func TestIndexCreate(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if derr := c.Del("table", "123::sample"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()
	form := map[string]interface{}{"table": 123}
	if cerr := c.CreateIndex(form, "sample"); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}
	form = map[string]interface{}{"table": 2.179699264}
	if cerr := c.CreateIndex(form, "sample"); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}
	form = map[string]interface{}{"table": []float64{2.179699264, 4.543211223}}
	if cerr := c.CreateIndex(form, "sample"); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}
	form = map[string]interface{}{"table": []string{"sample", "sample2"}}
	if cerr := c.CreateIndex(form, "sample"); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

}
func TestCreateIndex_invalidData(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}

	defer func() {
		c.Del("table", "key::sample")
	}()

	if cerr := c.CreateIndex(map[string]interface{}{"table": []string{"key"}}, "sample"); cerr != nil {
		if !(strings.Contains(cerr.Error(), "unsupported")) {
			t.Errorf("Error while making data entry: %v\n", cerr.Error())
		}
	}

}

func TestGet(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	if cerr := c.CreateIndex(map[string]interface{}{"model": "intel"}, "abc-123"); cerr != nil {
		t.Errorf("Error: %v\n", cerr.Error())
	}
	for i := 0; i < 10000; {
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
		t.Errorf("Error while reading data: %v\n", rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf("data not found")
	}
	got, rerr = c.GetString("ProcessorSummary/Model", 0, "value*", false)
	if rerr != nil {
		t.Errorf("Error while reading data: %v\n", rerr.Error())
	}
	if len(got) != 10000 {
		t.Errorf("data not found")
	}
	defer func() {
		if derr := c.Del("model", "intel::abc-123"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

}

func TestGetTaskList(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	if cerr := c.CreateIndex(map[string]interface{}{"count": 4}, "abc-123"); cerr != nil {
		t.Errorf("Error: %v\n", cerr.Error())
	}
	got, rerr := c.GetTaskList("count", 0, -1)
	if rerr != nil {
		t.Errorf("Error while reading data: %v\n", rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf("data not found")
	}
	defer func() {
		if derr := c.Del("count", "4::abc-123"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

}
func TestGetRange(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	if cerr := c.CreateIndex(map[string]interface{}{"count": 4}, "abc-123"); cerr != nil {
		t.Errorf("Error: %v\n", cerr.Error())
	}
	got, rerr := c.GetRange("count", 0, 5, false)
	if rerr != nil {
		t.Errorf("Error while reading data: %v\n", rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf("data not found")
	}
	defer func() {
		if derr := c.Del("count", "4::abc-123"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

}
func TestGetStorageList(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB connection:", err)
	}
	if cerr := c.CreateIndex(map[string]interface{}{"storage": "[2.179699264 4.543211223]"}, "abc-123"); cerr != nil {
		t.Errorf("Error: %v\n", cerr.Error())
	}
	got, rerr := c.GetStorageList("storage", 0, 2.179699264, "eq", false)
	if rerr != nil {
		t.Errorf("Error while reading data: %v\n", rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf("data not found")
	}
	got, rerr = c.GetStorageList("storage", 0, 2.179699264, "ge", false)
	if rerr != nil {
		t.Errorf("Error while reading data: %v\n", rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf("data not found")
	}
	got, rerr = c.GetStorageList("storage", 0, 0, "gt", false)
	if rerr != nil {
		t.Errorf("Error while reading data: %v\n", rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf("data not found")
	}
	got, rerr = c.GetStorageList("storage", 0, 4, "lt", false)
	if rerr != nil {
		t.Errorf("Error while reading data: %v\n", rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf("data not found")
	}
	got, rerr = c.GetStorageList("storage", 0, 4, "le", false)
	if rerr != nil {
		t.Errorf("Error while reading data: %v\n", rerr.Error())
	}
	if len(got) == 0 {
		t.Errorf("data not found")
	}

	defer func() {
		if derr := c.Del("storage", "*::abc-123"); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

}

func TestCreateEvtSubscriptions(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}
	js := `{"Name":"Subscriptions", "hostip":"10.10.10.10"}`
	defer func() {
		if derr := c.DeleteEvtSubscriptions("subscriptions", js); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

	if cerr := c.CreateEvtSubscriptionIndex("subscriptions", js); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}
}

func TestCreateEvtSubscriptions_existingData(t *testing.T) {

	c, err := MockDBConnection()
	if err != nil {
		t.Fatal("Error while making mock DB coonection:", err)
	}
	js := `{"Name":"Subscriptions", "hosts":["10.10.10.10"]}`
	defer func() {
		if derr := c.DeleteEvtSubscriptions("subscriptions", js); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

	if cerr := c.CreateEvtSubscriptionIndex("subscriptions", js); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

	if cerr := c.CreateEvtSubscriptionIndex("subscriptions", js); cerr == nil {
		t.Errorf("Error while making data entry")
	}

}

func TestGetEvtSubscriptions(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}
	js := `{"Name":"Subscriptions", "hostip":"10.10.10.10"}`
	defer func() {
		if derr := c.DeleteEvtSubscriptions("subscriptions", js); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()

	if cerr := c.CreateEvtSubscriptionIndex("subscriptions", js); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

	subscriptions, gerr := c.GetEvtSubscriptions("subscriptions", "*10.10.10.10*")
	if gerr != nil {
		t.Errorf("Error while getting data: %v\n", gerr.Error())
	}
	if len(subscriptions) < 1 {
		t.Errorf("No data found for the host ip")
	}

	if subscriptions[0] != js {
		t.Errorf("Error while trying to get data")
	}

	subscriptions, gerr = c.GetEvtSubscriptions("subscriptions", "*10.10.10.100*")
	if err != nil {
		t.Errorf("Error while getting data: %v\n", err.Error())
	}
	if len(subscriptions) > 0 {
		t.Errorf("data shouldn't be there")
	}

}

func TestDeleteEvtSubscriptions(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}
	js := `{"Name":"Subscriptions", "hostip":"10.10.10.10"}`

	if cerr := c.CreateEvtSubscriptionIndex("subscriptions", js); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

	if derr := c.DeleteEvtSubscriptions("subscriptions", js); derr != nil {
		t.Errorf("Error while deleting Data: %v\n", derr.Error())
	}

}

func TestDeleteEvtSubscriptions_nonexisting_data(t *testing.T) {
	c, err := MockDBConnection()
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
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}

	js := `{"Name":"Subscriptions", "SubscriptionID":"12345","hostip":["10.10.10.10"]}`

	if cerr := c.CreateEvtSubscriptionIndex("subscriptions", js); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

	js = `{"Name":"Subscriptions", "SubscriptionID":"12345","hostip":["10.10.10.19"]}`

	if cerr := c.UpdateEvtSubscriptions("subscriptions", "*12345*", js); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

}

func TestCreateDeviceSubscription(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}
	location := "https://10.24.1.23/redfish/v1/EventService/Subscriptions/123"
	originResources := []string{"/redfish/v1/Systems/uuid:1"}
	hostIP := "10.24.0.1"
	defer func() {
		if derr := c.DeleteDeviceSubscription("DeviceSubscription", hostIP); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()
	if cerr := c.CreateDeviceSubscriptionIndex("DeviceSubscription", hostIP, location, originResources); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

}

func TestCreateDeviceSubscription_existingData(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}
	location := "https://10.24.1.23/redfish/v1/EventService/Subscriptions/123"
	originResources := []string{"/redfish/v1/Systems/uuid:1"}
	hostIP := "10.24.0.1"
	defer func() {
		if derr := c.DeleteDeviceSubscription("DeviceSubscription", hostIP); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()
	if cerr := c.CreateDeviceSubscriptionIndex("DeviceSubscription", hostIP, location, originResources); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

	if cerr := c.CreateDeviceSubscriptionIndex("DeviceSubscription", hostIP, location, originResources); cerr == nil {
		t.Errorf("Data Already Exist: %v\n", cerr.Error())
	}
}

func TestGetDeviceSubscription(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}
	location := "https://10.24.1.23/redfish/v1/EventService/Subscriptions/123"
	originResources := []string{"/redfish/v1/Systems/uuid:1"}
	hostIP := "10.24.0.1"
	defer func() {
		if derr := c.DeleteDeviceSubscription("DeviceSubscription", hostIP); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()
	if cerr := c.CreateDeviceSubscriptionIndex("DeviceSubscription", hostIP, location, originResources); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

	subscriptions, gerr := c.GetDeviceSubscription("DeviceSubscription", "10.24.0.1*")
	if gerr != nil {
		t.Errorf("Error while getting data: %v\n", gerr.Error())
	}
	if len(subscriptions) < 1 {
		t.Errorf("No data found for the host ip")
	}
	devSub := strings.Split(subscriptions[0], "::")

	if devSub[0] != hostIP && devSub[1] != location {
		t.Errorf("HostIP/Location didn't matched")
	}

	subscriptions, gerr = c.GetDeviceSubscription("DeviceSubscription", "10.10.10.100*")
	if err != nil {
		t.Errorf("Error while getting data: %v\n", err.Error())
	}
	if len(subscriptions) > 0 {
		t.Errorf("data shouldn't be there")
	}

}

func TestDeleteDeviceSubscription(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}
	location := "https://10.24.1.23/redfish/v1/EventService/Subscriptions/123"
	originResources := []string{"/redfish/v1/Systems/uuid:1"}
	hostIP := "10.24.0.1"
	if cerr := c.CreateDeviceSubscriptionIndex("DeviceSubscription", hostIP, location, originResources); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}
	if derr := c.DeleteDeviceSubscription("DeviceSubscription", hostIP); derr != nil {
		t.Errorf("Error while deleting Data: %v\n", derr.Error())
	}

}

func TestDeleteDeviceSubscriptions_nonexisting_data(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}
	if derr := c.DeleteDeviceSubscription("DeviceSubscription", "10.24.0.10"); derr == nil {
		t.Errorf("No Data found: %v\n", derr.Error())
	}

}

func TestUpdateDeviceSubscriptions(t *testing.T) {
	c, err := MockDBConnection()
	if err != nil {
		t.Fatal(err)
	}

	location := "https://10.24.1.23/redfish/v1/EventService/Subscriptions/123"
	originResources := []string{"/redfish/v1/Systems/uuid:1"}
	hostIP := "10.24.0.1"
	defer func() {
		if derr := c.DeleteDeviceSubscription("DeviceSubscription", hostIP); derr != nil {
			t.Errorf("Error while deleting Data: %v\n", derr.Error())
		}
	}()
	if cerr := c.CreateDeviceSubscriptionIndex("DeviceSubscription", hostIP, location, originResources); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

	location = "https://10.24.1.23/redfish/v1/EventService/Subscriptions/12345"
	if cerr := c.UpdateDeviceSubscription("DeviceSubscription", hostIP, location, originResources); cerr != nil {
		t.Errorf("Error while making data entry: %v\n", cerr.Error())
	}

}

type redisExtCallsImpMock struct{}

func (r redisExtCallsImpMock) newSentinelClient(opt *redisSentinel.Options) *redisSentinel.SentinelClient {
	return newSentinelClientMock(opt)
}
func newSentinelClientMock(opt *redisSentinel.Options) *redisSentinel.SentinelClient {
	strSlice := strings.Split(opt.Addr, ":")
	sentinelHost := strSlice[0]
	sentinelPort := strSlice[1]
	if sentinelHost == "ValidHost" && sentinelPort == "ValidSentinelPort" {
		return &redisSentinel.SentinelClient{}
	}
	return nil
}
func (r redisExtCallsImpMock) getMasterAddrByName(masterSet string, snlClient *redisSentinel.SentinelClient) []string {
	return getMasterAddbyNameMock(masterSet, snlClient)
}

func getMasterAddbyNameMock(masterSet string, snlClient *redisSentinel.SentinelClient) []string {
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
			got, got1 := GetCurrentMasterHostPort(tt.args.dbConfig)
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
func TestGetDBConnection_HAEnabled(t *testing.T) {
	GetMockDBConfig()
	// Enableing HA
	config.Data.DBConf.RedisHAEnabled = true

	inMemDBConnPool = &ConnPool{
		ReadPool:        &redis.Pool{},
		WritePool:       nil,
		MasterIP:        "NotValid",
		PoolUpdatedTime: time.Now(),
	}
	onDiskDBConnPool = &ConnPool{
		ReadPool:        &redis.Pool{},
		WritePool:       nil,
		MasterIP:        "NotValid",
		PoolUpdatedTime: time.Now(),
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
