package rpc

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"testing"

	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-telemetry/telemetry"
	tm "github.com/ODIM-Project/ODIM/svc-telemetry/telemetry"
)

func Test_GetTele(t *testing.T) {
	testCases := []struct {
		desc string
		want *Telemetry
	}{
		{
			desc: "Get telemetry instance",
			want: &Telemetry{
				connector: telemetry.GetExternalInterface(),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := GetTele()
			if reflect.TypeOf(got) != reflect.TypeOf(tC.want) {
				t.Error(fmt.Errorf("GetTele()- want: %+v, Got: %+v", tC.want, got))
			}
		})
	}
}

func Test_generateResponse(t *testing.T) {
	testCases := []struct {
		desc  string
		input interface{}
		want  string
	}{
		{
			desc:  "success case",
			input: "input",
			want: func() string {
				w, _ := json.Marshal("input")
				return string(w)
			}(),
		},
		{
			desc:  "error case",
			input: math.Inf(1),
			want:  "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := generateResponse(tC.input)
			if string(got) != tC.want {
				t.Error(fmt.Errorf("generateResponse()- want: %+v, Got: %+v", tC.want, string(got)))
			}
		})
	}
}

func Test_fillProtoResponse(t *testing.T) {
	type args struct {
		resp *teleproto.TelemetryResponse
		data response.RPC
	}
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	testCases := []struct {
		desc string
		args *args
	}{
		{
			desc: "fill proto response",
			args: &args{
				resp: &teleproto.TelemetryResponse{},
				data: telemetry.connector.GetTelemetryService(),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			fillProtoResponse(tC.args.resp, tC.args.data)
		})
	}
}

func Test_generateRPCResponse(t *testing.T) {
	type args struct {
		rpcResp  response.RPC
		teleResp *teleproto.TelemetryResponse
	}
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	testCases := []struct {
		desc string
		args *args
	}{
		{
			desc: "generate rpc response",
			args: &args{
				rpcResp:  telemetry.connector.GetTelemetryService(),
				teleResp: &teleproto.TelemetryResponse{},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			generateRPCResponse(tC.args.rpcResp, tC.args.teleResp)
		})
	}
}

func Test_generateTaskRespone(t *testing.T) {
	type args struct {
		rpcResp response.RPC
		taskID  string
		taskURI string
	}
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	testCases := []struct {
		desc string
		args *args
	}{
		{
			desc: "generate rpc response",
			args: &args{
				rpcResp: telemetry.connector.GetTelemetryService(),
				taskID:  "taskID",
				taskURI: "/task/uri",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			generateTaskRespone(tC.args.taskID, tC.args.taskURI, &tC.args.rpcResp)
		})
	}
}
