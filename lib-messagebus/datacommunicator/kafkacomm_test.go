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
package datacommunicator

import (
	"crypto/tls"
	"reflect"
	"testing"

	"github.com/segmentio/kafka-go"
)

func TestTLS(t *testing.T) {
	type args struct {
		cCert  string
		cKey   string
		caCert string
	}
	tests := []struct {
		name    string
		args    args
		want    *tls.Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TLS(tt.args.cCert, tt.args.cKey, tt.args.caCert)
			if (err != nil) != tt.wantErr {
				t.Errorf("TLS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TLS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKafkaConnect(t *testing.T) {
	type args struct {
		kp *KafkaPacket
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := kafkaConnect(tt.args.kp); (err != nil) != tt.wantErr {
				t.Errorf("kafkaConnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKafkaPacket_Distribute(t *testing.T) {
	type fields struct {
		Packet     Packet
		DialerConn *kafka.Dialer
		Server     []string
	}
	type args struct {
		d interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kp := &KafkaPacket{
				Packet:      tt.fields.Packet,
				DialerConn:  tt.fields.DialerConn,
				ServersInfo: tt.fields.Server,
			}
			kp.Distribute(tt.args.d)
		})
	}
}

func TestKafkaPacket_Accept(t *testing.T) {
	type fields struct {
		Packet     Packet
		DialerConn *kafka.Dialer
		Server     []string
	}
	type args struct {
		fn MsgProcess
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kp := &KafkaPacket{
				Packet:      tt.fields.Packet,
				DialerConn:  tt.fields.DialerConn,
				ServersInfo: tt.fields.Server,
			}
			if err := kp.Accept(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("KafkaPacket.Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKafkaPacket_Read(t *testing.T) {
	type fields struct {
		Packet     Packet
		DialerConn *kafka.Dialer
		Server     []string
	}
	type args struct {
		fn MsgProcess
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kp := &KafkaPacket{
				Packet:      tt.fields.Packet,
				DialerConn:  tt.fields.DialerConn,
				ServersInfo: tt.fields.Server,
			}
			if err := kp.Read(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("KafkaPacket.Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKafkaPacket_Get(t *testing.T) {
	type fields struct {
		Packet     Packet
		DialerConn *kafka.Dialer
		Server     []string
	}
	type args struct {
		pipe string
		d    interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kp := &KafkaPacket{
				Packet:      tt.fields.Packet,
				DialerConn:  tt.fields.DialerConn,
				ServersInfo: tt.fields.Server,
			}
			if got := kp.Get(tt.args.pipe, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KafkaPacket.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKafkaPacket_Close(t *testing.T) {
	type fields struct {
		Packet     Packet
		DialerConn *kafka.Dialer
		Server     []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kp := &KafkaPacket{
				Packet:      tt.fields.Packet,
				DialerConn:  tt.fields.DialerConn,
				ServersInfo: tt.fields.Server,
			}
			kp.Close()
		})
	}
}
