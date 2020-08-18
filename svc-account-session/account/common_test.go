package account

import (
	"testing"
	"encoding/base64"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"golang.org/x/crypto/sha3"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
)

func TestGetExternalInterface(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExternalInterface(); got == nil {
				t.Errorf("Result of GetExternalInterface() should not be equal to nil")
			}
		})
	}
}

func getMockExternalInterface() *ExternalInterface {
	return &ExternalInterface {
		GetUserDetails: mockGetUserDetails,
		GetRoleDetailsByID: asmodel.GetRoleDetailsByID,
	}
}

func mockGetUserDetails(userName string) (asmodel.User, *errors.Error) {
	hash := sha3.New512()
	hash.Write([]byte("P@$$w0rd"))
	hashSum := hash.Sum(nil)
	hashedPassword := base64.URLEncoding.EncodeToString(hashSum)

	user := asmodel.User{
		UserName: userName,
		Password: hashedPassword,
	}

	if userName == "testUser1" || userName == "testUser2" {
		user.RoleID = common.RoleAdmin
	} else if userName == "testUser3" {
		user.RoleID = "PrivilegeLogin"
	} else {
		return user, errors.PackError(errors.DBKeyNotFound, "error while trying to get user: ", "fdfdfdfdfdfdfdf")
	}
	return user, nil
}

//func mockGetRoleDetailsByID(roleID string) (Role, *errors.Error) {

//}
