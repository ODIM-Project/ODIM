package rpc

import (
	"context"
	"reflect"
	"testing"

	authproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/auth"
)

func TestAuth_IsAuthorized(t *testing.T) {
	type args struct {
		ctx context.Context
		req *authproto.AuthRequest
	}
	tests := []struct {
		name     string
		args     args
		AuthFunc func(ctx context.Context, req *authproto.AuthRequest) (int32, string)
		want     *authproto.AuthResponse
		wantErr  bool
	}{
		{
			name:     "No error",
			args:     args{context.Background(), &authproto.AuthRequest{}},
			AuthFunc: func(ctx context.Context, req *authproto.AuthRequest) (int32, string) { return 200, "123Success" },
			want:     &authproto.AuthResponse{StatusCode: 200, StatusMessage: "123Success"},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		AuthFunc = tt.AuthFunc
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{}
			got, err := a.IsAuthorized(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsAuthorized() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsAuthorized() got = %v, want %v", got, tt.want)
			}
		})
	}
}
