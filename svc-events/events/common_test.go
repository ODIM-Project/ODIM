// Package events have the functionality of
// - Create Event Subscription
// - Delete Event Subscription
// - Get Event Subscription
// - Post Event Subscription to destination
// - Post TestEvent (SubmitTestEvent)
// and corresponding unit test cases
package events

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
)

func TestUpdateTaskData(t *testing.T) {

	type args struct {
		taskData common.TaskData
	}
	tests := []struct {
		name              string
		args              args
		UpdateTaskService func(ctx context.Context, taskID, taskState, taskStatus string, percentComplete int32, payLoad *task.Payload, endTime time.Time) error
		wantErr           error
	}{
		{
			name: "Update Task Service",
			args: args{
				taskData: common.TaskData{},
			},
			UpdateTaskService: func(ctx context.Context, taskID, taskState, taskStatus string, percentComplete int32, payLoad *task.Payload, endTime time.Time) error {
				return errors.New(common.Cancelling)
			},
			wantErr: errors.New(common.Cancelling),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateTaskService = tt.UpdateTaskService
			if err := UpdateTaskData(mockContext(), tt.args.taskData); errors.Is(err, tt.wantErr) {
				t.Errorf("UpdateTaskData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_removeElement(t *testing.T) {
	type args struct {
		slice   []string
		element string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "Positive test case",
			args: args{
				slice:   []string{"data1", "data2"},
				element: "data1",
			},
			want: []string{"data2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeElement(tt.args.slice, tt.args.element); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalInterfaces_PluginCall(t *testing.T) {
	config.SetUpMockConfig(t)
	e := getMockMethods()
	e.ContactClient = func(ctx context.Context, s1, s2, s3, s4 string, i interface{}, m map[string]string) (*http.Response, error) {
		return &http.Response{Status: "Save", StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, nil
	}
	IOUtilReadAllFunc = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("")
	}
	e.PluginCall(evcommon.PluginContactRequest{Plugin: &evmodel.Plugin{IP: "10.10.10"}, HTTPMethodType: http.MethodPost})
	IOUtilReadAllFunc = func(r io.Reader) ([]byte, error) {
		return ioutil.ReadAll(r)
	}
	isHostPresentInEventForward([]string{}, "test")
	isHostPresent([]string{}, "test")
}

func mockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, "xyz")
	ctx = context.WithValue(ctx, common.ActionID, "001")
	ctx = context.WithValue(ctx, common.ActionName, "xyz")
	ctx = context.WithValue(ctx, common.ThreadID, "0")
	ctx = context.WithValue(ctx, common.ThreadName, "xyz")
	ctx = context.WithValue(ctx, common.ProcessName, "xyz")
	return ctx
}
