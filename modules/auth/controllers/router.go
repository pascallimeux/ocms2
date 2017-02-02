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
	"github.com/gorilla/mux"
	"github.com/pascallimeux/ocms2/modules/auth/model"
	"github.com/pascallimeux/ocms2/modules/auth/setting"
	"github.com/pascallimeux/ocms2/modules/log"
	"net/http"
)

const (
	USERURI     = "/o/user"
	REGISTERURI = "/o/register"
	ROLEURI     = "/o/role"
	AUTHURI     = "/o/auth"
	LOGURI      = "/o/log"
)

type AppContext struct {
	HttpServer *http.Server
	Settings   *setting.Settings
	SqlContext model.SqlContext
}

// Initialize API
func (appContext *AppContext) CreateAUTHRoutes() *mux.Router {
	log.Trace(log.Here(), "CreateAUTHRoutes() : calling method -")
	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc(REGISTERURI, appContext.registerUser).Methods("POST")     // create user
	router.HandleFunc(USERURI+"/{id}", appContext.getUser).Methods("GET")       // read a user
	router.HandleFunc(USERURI, appContext.updateUser).Methods("PUT")            // update user
	router.HandleFunc(USERURI+"/{id}", appContext.deleteUser).Methods("DELETE") // delete a user
	router.HandleFunc(USERURI, appContext.getUsers).Methods("GET")              // get liste of users

	router.HandleFunc(ROLEURI, appContext.postRole).Methods("POST")       // create a role
	router.HandleFunc(ROLEURI+"/{id}", appContext.getRole).Methods("GET") // read a role
	router.HandleFunc(ROLEURI, appContext.getRoles).Methods("GET")        // get liste of roles

	router.HandleFunc(AUTHURI, appContext.getToken).Methods("POST") // get a token

	router.HandleFunc(LOGURI+"/{from}/{to}", appContext.getLogs4dates).Methods("GET") // get logs for a periode
	router.HandleFunc(LOGURI, appContext.getLogs).Methods("GET")                      // get all logs
	return router
}
