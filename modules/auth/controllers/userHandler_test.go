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
	"github.com/pascallimeux/ocms2/modules/auth/model"
	"github.com/pascallimeux/ocms2/modules/auth/setting"
	"github.com/pascallimeux/ocms2/modules/log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

var httpServerTest *httptest.Server
var logfile *os.File
var DeployTimeout time.Duration
var TransactionTimeout time.Duration
var sqlContext model.SqlContext

func setup() {
	var err error

	// Init config
	config := setting.Settings{DataSourceName: "/tmp/auth_test.db", LogFileName: "/tmp/test.log", LogMode: "Trace", ExpireInToken: 24}

	// Init logger
	logfile = log.Init_log(config.LogFileName, config.LogMode)

	// Init application context
	appContext := AppContext{Settings: &config}

	// Init sqliteDB
	sqlContext, err = model.GetSqlContext(config.DataSourceName, true)
	if err != nil {
		log.Fatal(log.Here(), err.Error())
	}
	appContext.SqlContext = sqlContext

	// Init http server
	router := appContext.CreateAUTHRoutes()
	httpServerTest = httptest.NewServer(router)
}

func shutdown() {
	log.Trace(log.Here(), "End of tests..")
	defer sqlContext.Db.Close()
	defer logfile.Close()
	defer httpServerTest.Close()
}

func TestCreateUserNominal(t *testing.T) {
	token, err0 := getToken(setting.ADMINLOGIN, setting.ADMINPWD)
	if err0 != nil {
		t.Error(err0)
	}
	user, statusCode, err := createUser(token.Token, HttpUser{Username: "username1", Lastname: "lastname1", Firstname: "firstname1", Email: "email1", Password: "password1", Role_id: 1})
	if err != nil {
		t.Error(err)
	}
	if user.Id == "" {
		t.Error("bad user ID")
	}
	if statusCode != http.StatusCreated {
		t.Error("Non-expected status code: ", http.StatusCreated, " ", statusCode)
	}
}

func TestGetUserNominal(t *testing.T) {
	token, err0 := getToken(setting.ADMINLOGIN, setting.ADMINPWD)
	if err0 != nil {
		t.Error(err0)
	}
	user, _, err1 := createUser(token.Token, HttpUser{Username: "username2", Lastname: "lastname2", Firstname: "firstname2", Email: "email2", Password: "password2", Role_id: 1})
	if err1 != nil {
		t.Error(err1)
	}
	user2, statusCode, err2 := getUser(token.Token, user.Id)
	if err2 != nil {
		t.Error(err2)
	}
	if user2.Id != user.Id {
		t.Error("bad user ID")
	}
	if statusCode != http.StatusOK {
		t.Error("Non-expected status code: %v\n\tbody: %v\n", http.StatusOK, statusCode)
	}
}

func TestGetUsersNominal(t *testing.T) {
	token, err0 := getToken(setting.ADMINLOGIN, setting.ADMINPWD)
	if err0 != nil {
		t.Error(err0)
	}
	createUser(token.Token, HttpUser{Username: "username3", Lastname: "lastname3", Firstname: "firstname3", Email: "email3", Password: "password3", Role_id: 1})
	createUser(token.Token, HttpUser{Username: "usename4", Lastname: "lastname4", Firstname: "firstname4", Email: "email4", Password: "password4", Role_id: 1})
	createUser(token.Token, HttpUser{Username: "username5", Lastname: "lastname5", Firstname: "firstname5", Email: "email5", Password: "password5", Role_id: 1})
	users, statusCode, err := getListOfUsers(token.Token)
	if err != nil {
		t.Error(err)
	}
	for _, user := range users {
		jsonUser, _ := json.Marshal(user)
		t.Log(string(jsonUser))
	}
	if statusCode != http.StatusOK {
		t.Error("Non-expected status code: %v\n\tbody: %v\n", http.StatusOK, statusCode)
	}
}

func TestUpdateUserNominal(t *testing.T) {
	token, err0 := getToken(setting.ADMINLOGIN, setting.ADMINPWD)
	if err0 != nil {
		t.Error(err0)
	}
	httpuser := HttpUser{Username: "username6", Lastname: "lastname6", Firstname: "firstname6", Email: "email6", Password: "password6", Role_id: 1}
	user, _, _ := createUser(token.Token, httpuser)
	newUser := "{\"Id\":\"" + user.Id + "\",\"Lastname\":\"lastname66\",\"Firstname\":\"firstname66\",\"Email\":\"email66\",\"Password\":\"password66\"}"
	user2, _, err := updateUser(token.Token, newUser)
	if err != nil {
		t.Error(err)
	}
	user2, _, err2 := getUser(token.Token, user2.Id)
	if err2 != nil {
		t.Error(err2)
	}
	jsonUser, _ := json.Marshal(user2)
	t.Log(string(jsonUser))

	if user2.Username != "username6" || user2.Lastname != "lastname66" {
		t.Error("bad username or bad lastname")
	}
}

func TestDeleteUserNominal(t *testing.T) {
	token, err0 := getToken(setting.ADMINLOGIN, setting.ADMINPWD)
	if err0 != nil {
		t.Error(err0)
	}
	user, _, _ := createUser(token.Token, HttpUser{Username: "username7", Lastname: "lastname7", Firstname: "firstname7", Email: "email7", Password: "password7", Role_id: 1})
	err := deleteUser(token.Token, user.Id)
	if err != nil {
		t.Error(err)
	}
	user2, _, err2 := getUser(token.Token, user.Id)
	if err2 != nil {
		t.Error(err2)
	}
	if user2.Activated == true {
		t.Error("user is not deleted")
	}
}

func TestCreateUserWithBadToken(t *testing.T) {
	_, statusCode, err := createUser("1111111111111111111111111111111111111111", HttpUser{Username: "username1", Lastname: "lastname1", Firstname: "firstname1", Email: "email1", Password: "password1", Role_id: 1})
	if err != nil {
		t.Error(err)
	}
	if statusCode != http.StatusUnauthorized {
		t.Error("Non-expected status code: ", http.StatusUnauthorized, " ", statusCode)
	}
}

func TestCreateUserWithBadRole(t *testing.T) {
	username := "user2"
	password := "user_pwd2"
	email := "user2@orange.fr"
	role := 3

	token, err0 := getToken(setting.ADMINLOGIN, setting.ADMINPWD)
	if err0 != nil {
		t.Error(err0)
	}
	_, _, err := createUser(token.Token, HttpUser{Username: username, Lastname: "", Firstname: "", Email: email, Password: password, Role_id: role})
	if err != nil {
		t.Error(err)
	}

	token, err2 := getToken(username, password)
	if err2 != nil {
		t.Error(err2)
	}

	_, statusCode, err := createUser(token.Token, HttpUser{Username: "username1", Lastname: "lastname1", Firstname: "firstname1", Email: "email1", Password: "password1", Role_id: 1})
	if err != nil {
		t.Error(err)
	}
	if statusCode != http.StatusUnauthorized {
		t.Error("Non-expected status code: ", http.StatusUnauthorized, " ", statusCode)
	}
}

func TestCreateUserWithBadParams(t *testing.T) {
	token, err0 := getToken(setting.ADMINLOGIN, setting.ADMINPWD)
	if err0 != nil {
		t.Error(err0)
	}
	_, statusCode, err := createUser(token.Token, HttpUser{Lastname: "lastname1", Firstname: "firstname1", Email: "email1", Password: "password1", Role_id: 1})
	if err != nil {
		t.Error(err)
	}
	if statusCode != http.StatusBadRequest {
		t.Error("Non-expected status code: ", http.StatusBadRequest, " ", statusCode)
	}
}

func TestUpdateUserWithBadToken(t *testing.T) {
	token, err0 := getToken(setting.ADMINLOGIN, setting.ADMINPWD)
	if err0 != nil {
		t.Error(err0)
	}
	httpuser := HttpUser{Username: "username8", Lastname: "lastname8", Firstname: "firstname8", Email: "email8", Password: "password8", Role_id: 1}
	user, _, _ := createUser(token.Token, httpuser)
	newUser := "{\"Id\":\"" + user.Id + "\",\"Lastname\":\"lastname66\",\"Firstname\":\"firstname66\",\"Email\":\"email66\",\"Password\":\"password66\"}"
	_, statusCode, err := updateUser("1111111111111111111111111111111111111111", newUser)
	if err != nil {
		t.Error(err)
	}
	if statusCode != http.StatusUnauthorized {
		t.Error("Non-expected status code: ", http.StatusUnauthorized, " ", statusCode)
	}
}
