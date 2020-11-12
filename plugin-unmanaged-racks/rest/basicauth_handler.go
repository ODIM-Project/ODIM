package rest

import (
	"encoding/base64"
	"net/http"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"golang.org/x/crypto/sha3"
)

func NewBasicAuthHandler(username, hashedPass string) context.Handler {
	bah := basicAuthHandler{
		username:   username,
		hashedPass: hashedPass,
	}
	return bah.handle
}

type basicAuthHandler struct {
	username, hashedPass string
}

func (b basicAuthHandler) handle(ctx iris.Context) {
	username, password, ok := ctx.Request().BasicAuth()
	if !ok {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(redfish.NewError().AddExtendedInfo(redfish.NewResourceAtURIUnauthorizedMsg(ctx.Request().RequestURI, "Cannot decode Authorization header")))
		return
	}

	if username != b.username {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(redfish.NewError().AddExtendedInfo(redfish.NewResourceAtURIUnauthorizedMsg(ctx.Request().RequestURI, "Invalid user or password")))
		return
	}

	hash := sha3.New512()
	hash.Write([]byte(password))
	hashSum := hash.Sum(nil)
	hashedPassword := base64.URLEncoding.EncodeToString(hashSum)
	if b.hashedPass != hashedPassword {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(redfish.NewError().AddExtendedInfo(redfish.NewResourceAtURIUnauthorizedMsg(ctx.Request().RequestURI, "Invalid user or password")))
		return
	}

	ctx.Next()
}
