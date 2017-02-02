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
	"github.com/gorilla/mux"
	"github.com/pascallimeux/ocms2/modules/auth/model"
	"github.com/pascallimeux/ocms2/modules/common"
	"github.com/pascallimeux/ocms2/modules/log"
	"net/http"
	"strconv"
)

//HTTP Post - /o/role
func (a *AppContext) postRole(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "postRole() : calling method -")

	err1 := a.CheckPermissionFromToken(w, r, "postRole", "")
	if err1 != nil {
		return
	}

	var role model.Role
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		common.SendError(log.Here(), w, err)
		return
	}
	role, err = a.SqlContext.CreateRole(role)
	if err != nil {
		common.SendError(log.Here(), w, err)
		return
	}
	roleString, _ := json.Marshal(role)
	log.Trace(log.Here(), "create role:", string(roleString))
	common.BuildHttp201Response(w, role)
}

//HTTP Get - /o/role/{id}
func (a *AppContext) getRole(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "getRole() : calling method -")

	err1 := a.CheckPermissionFromToken(w, r, "getRole", "")
	if err1 != nil {
		return
	}

	vars := mux.Vars(r)
	code, err1 := strconv.Atoi(vars["id"])
	if err1 != nil {
		common.SendError(log.Here(), w, err1)
		return
	}
	role, err2 := a.SqlContext.GetRole(code)
	if err2 != nil {
		common.SendError(log.Here(), w, err2)
		return
	}
	common.BuildHttp200Response(w, role)
}

//HTTP Get - /o/role
func (a *AppContext) getRoles(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "getRoles() : calling method -")

	err1 := a.CheckPermissionFromToken(w, r, "getRoles", "")
	if err1 != nil {
		return
	}

	roles, err := a.SqlContext.GetRoles()
	if err != nil {
		common.SendError(log.Here(), w, err)
		return
	}
	common.BuildHttp200Response(w, roles)
}
