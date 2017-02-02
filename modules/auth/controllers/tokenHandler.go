/*
Copyright Pascal Limeux. 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"encoding/json"
	"errors"
	"github.com/pascallimeux/ocms2/modules/common"
	"github.com/pascallimeux/ocms2/modules/log"
	"net/http"
	"strings"
)

// create a token for a user
//HTTP Post - /o/auth
func (a *AppContext) getToken(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "getToken() : calling method -")

	type Authent struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var authent Authent
	err := json.NewDecoder(r.Body).Decode(&authent)
	if err != nil {
		common.SendError(log.Here(), w, err)
		return
	}

	if authent.Username == "" {
		common.SendError(log.Here(), w, errors.New("no username given"))
		return
	}
	if authent.Password == "" {
		common.SendError(log.Here(), w, errors.New("no password given"))
		return
	}

	user, err2 := a.SqlContext.GetUserByCredentials(authent.Username, authent.Password)
	if err2 != nil {
		common.SendError(log.Here(), w, err2)
		return
	}

	if !user.Activated {
		common.SendError(log.Here(), w, errors.New("user not activated"))
		return
	}

	userRole, _ := a.SqlContext.GetRole(user.Role_id)
	authorized := a.SqlContext.IsPermitted4User(userRole, user.Id, user.Username, "getToken", "")
	if authorized {
		log.Trace(log.Here(), "The user: ", user.Username, " is granted to get a token")
	} else {
		log.Trace(log.Here(), "The user: ", user.Username, " is not authorized get a token")
		common.SendError(log.Here(), w, errors.New("User not authorized for this resource!"))
		return
	}
	expire_in := a.Settings.ExpireInToken
	token, err3 := a.SqlContext.CreateToken(user, expire_in)
	if err3 != nil {
		common.SendError(log.Here(), w, err3)
		return
	}
	tokenString, _ := json.Marshal(token)
	log.Trace(log.Here(), "create token:", string(tokenString))
	common.BuildHttp201Response(w, token)
}

func (a *AppContext) CheckPermissionFromToken(w http.ResponseWriter, r *http.Request, resourceName, resourceId string) error {
	log.Trace(log.Here(), "checkPermissionFromToken() : calling method -")
	tokenValue, err1 := extractTokenFromHeader(r)
	if err1 != nil {
		log.Trace(log.Here(), err1.Error())
		return err1
	}
	err := a.SqlContext.IsAuthorized4Token(tokenValue, resourceName, resourceId)
	if err != nil {
		common.BuildHttp401Response(w)
	}
	return err
}

func extractTokenFromHeader(r *http.Request) (string, error) {
	log.Trace(log.Here(), "extractTokenFromHeader() : calling method -")
	tokenValue := r.Header.Get("authorization")
	if tokenValue == "" {
		return "", errors.New("no token in header")
	}
	space := strings.LastIndex(tokenValue, " ")
	return tokenValue[space+1:], nil
}
