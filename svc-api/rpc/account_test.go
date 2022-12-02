package rpc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	"google.golang.org/grpc"
)

// type fakeStruct struct{}

func TestDoAccountCreationRequest(t *testing.T) {
	type args struct {
		req accountproto.CreateAccountRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewAccountClientFunc func(cc *grpc.ClientConn) accountproto.AccountClient
		want                 *accountproto.AccountResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{accountproto.CreateAccountRequest{}},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAccountClientFunc: func(cc *grpc.ClientConn) accountproto.AccountClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "Account Create error",
			args:                 args{accountproto.CreateAccountRequest{}},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAccountClientFunc: func(cc *grpc.ClientConn) accountproto.AccountClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAccountClientFunc = tt.NewAccountClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoAccountCreationRequest(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoAccountCreationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoAccountCreationRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoAccountDeleteRequest(t *testing.T) {
	type args struct {
		req accountproto.DeleteAccountRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewAccountClientFunc func(cc *grpc.ClientConn) accountproto.AccountClient
		want                 *accountproto.AccountResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{accountproto.DeleteAccountRequest{}},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAccountClientFunc: func(cc *grpc.ClientConn) accountproto.AccountClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "Account Delete error",
			args:                 args{accountproto.DeleteAccountRequest{}},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAccountClientFunc: func(cc *grpc.ClientConn) accountproto.AccountClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAccountClientFunc = tt.NewAccountClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoAccountDeleteRequest(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoAccountDeleteRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoAccountDeleteRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetAccountRequest(t *testing.T) {
	type args struct {
		req accountproto.GetAccountRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewAccountClientFunc func(cc *grpc.ClientConn) accountproto.AccountClient
		want                 *accountproto.AccountResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{accountproto.GetAccountRequest{}},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAccountClientFunc: func(cc *grpc.ClientConn) accountproto.AccountClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetAccount error",
			args:                 args{accountproto.GetAccountRequest{}},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAccountClientFunc: func(cc *grpc.ClientConn) accountproto.AccountClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAccountClientFunc = tt.NewAccountClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetAccountRequest(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetAccountRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetAccountRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetAccountServiceRequest(t *testing.T) {
	type args struct {
		req accountproto.AccountRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewAccountClientFunc func(cc *grpc.ClientConn) accountproto.AccountClient
		want                 *accountproto.AccountResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{accountproto.AccountRequest{}},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAccountClientFunc: func(cc *grpc.ClientConn) accountproto.AccountClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetAccountServices error",
			args:                 args{accountproto.AccountRequest{}},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAccountClientFunc: func(cc *grpc.ClientConn) accountproto.AccountClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAccountClientFunc = tt.NewAccountClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetAccountServiceRequest(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetAccountServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetAccountServiceRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetAllAccountRequest(t *testing.T) {
	type args struct {
		req accountproto.AccountRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewAccountClientFunc func(cc *grpc.ClientConn) accountproto.AccountClient
		want                 *accountproto.AccountResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{accountproto.AccountRequest{}},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAccountClientFunc: func(cc *grpc.ClientConn) accountproto.AccountClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetAllAccounts error",
			args:                 args{accountproto.AccountRequest{}},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAccountClientFunc: func(cc *grpc.ClientConn) accountproto.AccountClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAccountClientFunc = tt.NewAccountClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetAllAccountRequest(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetAllAccountRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetAllAccountRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoUpdateAccountRequest(t *testing.T) {
	type args struct {
		req accountproto.UpdateAccountRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewAccountClientFunc func(cc *grpc.ClientConn) accountproto.AccountClient
		want                 *accountproto.AccountResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{accountproto.UpdateAccountRequest{}},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAccountClientFunc: func(cc *grpc.ClientConn) accountproto.AccountClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "Account Update error",
			args:                 args{accountproto.UpdateAccountRequest{}},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAccountClientFunc: func(cc *grpc.ClientConn) accountproto.AccountClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAccountClientFunc = tt.NewAccountClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoUpdateAccountRequest(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoUpdateAccountRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoUpdateAccountRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
