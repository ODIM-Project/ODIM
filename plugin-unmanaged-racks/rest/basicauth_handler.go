/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rest

import (
	"encoding/base64"
	"net/http"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"golang.org/x/crypto/sha3"
)

func newBasicAuthHandler(username, hashedPass string) context.Handler {
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
