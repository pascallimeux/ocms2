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
	"github.com/gorilla/mux"
	"github.com/pascallimeux/ocms2/modules/common"
	"github.com/pascallimeux/ocms2/modules/log"
	"net/http"
)

//HTTP Post - /o/register
func (a *AppContext) registerUser(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "registerUser() : calling method -")

	err1 := a.CheckPermissionFromToken(w, r, "registerUser", "")
	if err1 != nil {
		return
	}

	type User struct {
		Username  string `json:"username"`
		Lastname  string `json:"lastname"`
		Firstname string `json:"firstname"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Role_id   int    `json:"role_id"`
	}
	var httpUser User
	err := json.NewDecoder(r.Body).Decode(&httpUser)
	if err != nil {
		common.SendError(log.Here(), w, err)
		return
	}
	if httpUser.Username == "" {
		common.SendError(log.Here(), w, errors.New("no username given"))
		return
	}
	if httpUser.Email == "" {
		common.SendError(log.Here(), w, errors.New("no email given"))
		return
	}
	if httpUser.Password == "" {
		common.SendError(log.Here(), w, errors.New("no password given"))
		return
	}
	user, err := a.SqlContext.CreateUser(httpUser.Username, httpUser.Lastname, httpUser.Firstname, httpUser.Email, httpUser.Password, httpUser.Role_id)
	if err != nil {
		common.SendError(log.Here(), w, err)
		return
	}
	user.Password = nil
	userString, _ := json.Marshal(user)
	log.Trace(log.Here(), "register user:", string(userString))
	common.BuildHttp201Response(w, user)
}

//HTTP Get - /o/user/{id}
func (a *AppContext) getUser(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "getUser() : calling method -")

	vars := mux.Vars(r)
	userid := vars["id"]
	err1 := a.CheckPermissionFromToken(w, r, "getUser", userid)
	if err1 != nil {
		return
	}
	user, err := a.SqlContext.GetUser(userid)
	if err != nil {
		common.SendError(log.Here(), w, err)
		return
	}
	common.BuildHttp200Response(w, user)
}

//HTTP Get - /o/user
func (a *AppContext) getUsers(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "getUsers() : calling method -")

	err1 := a.CheckPermissionFromToken(w, r, "getUsers", "")
	if err1 != nil {
		return
	}

	users, err := a.SqlContext.GetUsers()
	if err != nil {
		common.SendError(log.Here(), w, err)
		return
	}
	common.BuildHttp200Response(w, users)
}

//HTTP Put - /o/user
func (a *AppContext) updateUser(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "updateUser() : calling method -")

	type User struct {
		Id        string `json:"id"`
		Lastname  string `json:"lastname"`
		Firstname string `json:"firstname"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	var httpUser User
	err := json.NewDecoder(r.Body).Decode(&httpUser)
	if err != nil {
		common.SendError(log.Here(), w, err)
		return
	}

	err0 := a.CheckPermissionFromToken(w, r, "updateUser", httpUser.Id)
	if err0 != nil {
		return
	}

	user, err1 := a.SqlContext.PutUser(httpUser.Id, httpUser.Lastname, httpUser.Firstname, httpUser.Email, httpUser.Password)
	if err1 != nil {
		common.SendError(log.Here(), w, err1)
		return
	}
	user.Password = nil
	userString, _ := json.Marshal(user)
	log.Trace(log.Here(), "updated user:", string(userString))
	common.BuildHttp200Response(w, user)
}

//HTTP Delete - /o/user/{id}
func (a *AppContext) deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "deleteUser() : calling method -")

	err0 := a.CheckPermissionFromToken(w, r, "deleteUser", "")
	if err0 != nil {
		return
	}

	vars := mux.Vars(r)
	userid := vars["id"]
	err := a.SqlContext.UnactivateUser(userid)
	if err != nil {
		common.SendError(log.Here(), w, err)
		return
	}
	common.BuildHttp204Response(w)
}
