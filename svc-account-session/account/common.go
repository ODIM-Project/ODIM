package account

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
)

// ExternalInterface holds all the external connections account package functions uses
type ExternalInterface struct {
	CreateUser         func(asmodel.User) *errors.Error
	GetUserDetails     func(string) (asmodel.User, *errors.Error)
	GetRoleDetailsByID func(string) (asmodel.Role, *errors.Error)
	UpdateUserDetails  func(asmodel.User, asmodel.User) *errors.Error
}

// GetExternalInterface retrieves all the external connections account package functions uses
func GetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		CreateUser:         asmodel.CreateUser,
		GetUserDetails:     asmodel.GetUserDetails,
		GetRoleDetailsByID: asmodel.GetRoleDetailsByID,
		UpdateUserDetails:  asmodel.UpdateUserDetails,
	}
}
