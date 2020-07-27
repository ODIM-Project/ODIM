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

// Package asmodel ...
package asmodel

import (
	"encoding/json"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

var sessionStore = common.InMemory

// Session will hold the data assosiated with the session
type Session struct {
	ID           string
	Token        string
	UserName     string
	RoleID       string
	Privileges   map[string]bool
	Origin       string
	CreatedTime  time.Time
	LastUsedTime time.Time
}

//CreateSession will hold input request for creating a session
type CreateSession struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}

// Persist will create a session in the DB
func (s *Session) Persist() *errors.Error {
	connPool, err := common.GetDBConnection(sessionStore)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connecting to DB: ", err.Error())
	}
	if err = connPool.Create("session", s.Token, s); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to create new session: ", err.Error())
	}
	return nil
}

// Update will update a session in the DB
func (s *Session) Update() *errors.Error {
	connPool, err := common.GetDBConnection(sessionStore)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connecting to DB: ", err.Error())
	}
	if _, err = connPool.Update("session", s.Token, s); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to update session: ", err.Error())
	}
	return nil
}

// GetSession will get the session details from db if available
func GetSession(token string) (Session, *errors.Error) {
	var session Session
	connPool, err := common.GetDBConnection(sessionStore)
	if err != nil {
		return session, errors.PackError(err.ErrNo(), "error while trying to connecting to DB: ", err.Error())
	}
	sessionData, err := connPool.Read("session", token)
	if err != nil {
		return session, errors.PackError(err.ErrNo(), "error while trying to get the session from DB: ", err.Error())
	}
	if jerr := json.Unmarshal([]byte(sessionData), &session); jerr != nil {
		return session, errors.PackError(errors.UndefinedErrorType, "error while trying to unmarshal session data: ", jerr)
	}
	return session, nil
}

// Delete will delete a session from the DB
func (s *Session) Delete() *errors.Error {
	connPool, err := common.GetDBConnection(sessionStore)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connecting to DB: ", err.Error())
	}
	if err = connPool.Delete("session", s.Token); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to delete session: ", err.Error())
	}
	return nil
}

// GetAllSessionKeys will collect all session keys available in the DB
func GetAllSessionKeys() ([]string, *errors.Error) {
	connPool, err := common.GetDBConnection(sessionStore)
	if err != nil {
		return nil, errors.PackError(err.ErrNo(), "error while trying to connecting to DB: ", err.Error())
	}
	sessionIDs, err := connPool.GetAllDetails("session")
	if err != nil {
		return nil, err
	}
	return sessionIDs, nil
}
