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

// Package auth ...
package auth

import (
	"context"
	"encoding/base64"
	"strconv"
	"sync"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"golang.org/x/crypto/sha3"
)

var lastExpiredSessionCleanUpTime time.Time

// Lock defines mutex lock to avoid race conditions
var Lock sync.Mutex

// CheckSessionCreationCredentials defines the auth at the time of session creation
func CheckSessionCreationCredentials(ctx context.Context, userName, password string) (*asmodel.User, *errors.Error) {
	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.CheckSessionCreation)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go expiredSessionCleanUp(ctxt)
	threadID++
	if userName == "" || password == "" {
		return nil, errors.PackError(errors.UndefinedErrorType, "error while checking session credentials: username or password is empty")
	}
	user, err := asmodel.GetUserDetails(userName)
	if err != nil {
		return nil, errors.PackError(err.ErrNo(), "error: Invalid username or password :", err.Error())
	}
	hash := sha3.New512()
	hash.Write([]byte(password))
	hashSum := hash.Sum(nil)
	hashedPassword := base64.URLEncoding.EncodeToString(hashSum)
	if user.Password != hashedPassword {
		return nil, errors.PackError(errors.UndefinedErrorType, "error while checking session credentials: input password is not matching user password")
	}
	return &user, nil
}

// CheckSessionTimeOut defines the session validity check
func CheckSessionTimeOut(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.CheckSessionTimeout)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go expiredSessionCleanUp(ctxt)
	threadID++
	if sessionToken == "" {
		return nil, errors.PackError(errors.InvalidAuthToken, "error: no session token found in header")
	}
	session, err := asmodel.GetSession(sessionToken)
	if err != nil {
		return nil, errors.PackError(err.ErrNo(), "error while trying to get session details with the token ", sessionToken, ": ", err.Error())
	}
	if time.Since(session.LastUsedTime).Minutes() > config.Data.AuthConf.SessionTimeOutInMins {
		return nil, errors.PackError(errors.InvalidAuthToken, "error: session is timed out", sessionToken)
	}

	return &session, nil
}

// expiredSessionCleanUp is for deleting timed out sessions from the db
func expiredSessionCleanUp(ctx context.Context) {
	Lock.Lock()
	defer Lock.Unlock()
	// checking whether the db is cleaned up recently
	if time.Since(lastExpiredSessionCleanUpTime).Minutes() > config.Data.AuthConf.ExpiredSessionCleanUpTimeInMins {
		sessionTokens, err := asmodel.GetAllSessionKeys()
		if err != nil {
			l.LogWithFields(ctx).Error("Unable to get all session tokens from DB: %v" + err.Error())
			return
		}

		for _, token := range sessionTokens {
			session, err := asmodel.GetSession(token)
			if err != nil {
				l.LogWithFields(ctx).Error("Unable to get session details with the token " + token + ": " + err.Error())
				continue
			}
			// checking for the timed out sessions
			if time.Since(session.LastUsedTime).Minutes() > config.Data.AuthConf.SessionTimeOutInMins {
				err = session.Delete()
				if err != nil {
					l.LogWithFields(ctx).Printf("Unable to delete expired session with token " + token + ": " + err.Error())
					continue
				}
			}
		}
		lastExpiredSessionCleanUpTime = time.Now()
		sessionTokens = nil
	}
}
