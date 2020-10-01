package rest

import (
	"encoding/base64"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"golang.org/x/crypto/sha3"
	"log"
	"net/http"
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
		errorMsg := "error: not a valid basic auth"
		log.Println("error:", errorMsg)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.WriteString(errorMsg)
		return
	}

	if username != b.username {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.WriteString("Invalid Username/Password")
		return
	}

	hash := sha3.New512()
	hash.Write([]byte(password))
	hashSum := hash.Sum(nil)
	hashedPassword := base64.URLEncoding.EncodeToString(hashSum)
	if b.hashedPass != hashedPassword {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.WriteString("Invalid Username/Password")
		return
	}

	ctx.Next()
}
