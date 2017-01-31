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

package model

import (
	"encoding/json"
	"github.com/pascallimeux/auth/common"
	"github.com/pascallimeux/auth/utils/log"
	"golang.org/x/crypto/bcrypt"
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

var logfile *os.File
var DeployTimeout time.Duration
var TransactionTimeout time.Duration
var sqlContext SqlContext

func setup() {
	var err error

	// Init config
	config := common.Configuration{DataSourceName: "/tmp/auth_test.db", LogFileName: "/tmp/test.log", Logger: "Trace"}

	// Init logger
	logfile = log.Init_log(config.LogFileName, config.Logger)

	// Init sqliteDB
	sqlContext, err = GetSqlContext(config.DataSourceName)
	if err != nil {
		log.Fatal(log.Here(), err.Error())
	}
	sqlContext.InitBDD()

}

func shutdown() {
	log.Trace(log.Here(), "End of tests..")
	defer sqlContext.Db.Close()
	defer logfile.Close()
}

func TestCreateUserNominal(t *testing.T) {
	username := "usernamedb1"
	lastname := "lastnamedb1"
	firstname := "firstnamedb1"
	email := "emaildb1"
	password := "passworddb1"
	role_id := 3
	user, err := sqlContext.CreateUser(username, lastname, firstname, email, password, role_id)
	if err != nil {
		t.Error(err)
	}
	if user.Activated != true {
		t.Error("bad value for activated parameter")
	}
	if user.Id == "" {
		t.Error("bad user ID")
	}
	if user.CreatedAt.After(time.Now()) {
		t.Error("bad created date")
	}
	if user.UpdatedAt.After(time.Now()) {
		t.Error("bad updated date")
	}
	if user.Username != username {
		t.Error("bad username")
	}
	if user.Firstname != firstname {
		t.Error("bad firstname")
	}
	if user.Lastname != lastname {
		t.Error("bad lastname")
	}
	if user.Email != email {
		t.Error("bad email")
	}
	if user.Role_id != role_id {
		t.Error("bad role_id")
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		t.Error("bad password")
	}
	userString, _ := json.Marshal(user)
	t.Log("user:", string(userString))
}

func TestGetUserByCredentialsNominal(t *testing.T) {
	username := "usernamedb2"
	lastname := "lastnamedb2"
	firstname := "firstnamedb2"
	email := "emaildb2"
	password := "passworddb2"
	role_id := 3
	user, err := sqlContext.CreateUser(username, lastname, firstname, email, password, role_id)
	if err != nil {
		t.Error(err)
	}
	user2, err2 := sqlContext.GetUserByCredentials(username, password)
	if err2 != nil {
		t.Error(err2)
	}
	if !user.UpdatedAt.Equal(user2.UpdatedAt) {
		t.Error("bad updated date")
	}
	if !user.CreatedAt.Equal(user2.CreatedAt) {
		t.Error("bad CreatedAt date")
	}
	userString, _ := json.Marshal(user)
	user2String, _ := json.Marshal(user2)
	t.Log("user:", string(userString))
	t.Log("user:", string(user2String))
}

func TestGetUserNominal(t *testing.T) {
	username := "usernamedb3"
	lastname := "lastnamedb3"
	firstname := "firstnamedb3"
	email := "emaildb3"
	password := "passworddb3"
	role_id := 3
	user, err := sqlContext.CreateUser(username, lastname, firstname, email, password, role_id)
	if err != nil {
		t.Error(err)
	}
	userID := user.Id
	user2, err2 := sqlContext.GetUser(userID)
	if err2 != nil {
		t.Error(err2)
	}
	if !user.UpdatedAt.Equal(user2.UpdatedAt) {
		t.Error("bad updated date")
	}
	if !user.CreatedAt.Equal(user2.CreatedAt) {
		t.Error("bad CreatedAt date")
	}
	userString, _ := json.Marshal(user)
	user2String, _ := json.Marshal(user2)
	t.Log("user:", string(userString))
	t.Log("user:", string(user2String))
}

func TestPutUserNominal(t *testing.T) {
	username := "usernamedb4"
	lastname := "lastnamedb4"
	firstname := "firstnamedb4"
	email := "emaildb4"
	password := "passworddb4"
	role_id := 3
	user, err := sqlContext.CreateUser(username, lastname, firstname, email, password, role_id)
	if err != nil {
		t.Error(err)
	}
	userID := user.Id
	lastname2 := "lastnamedb44"
	firstname2 := "firstnamedb4"
	email2 := "emaildb4"
	password2 := "passworddb4"
	user2, err2 := sqlContext.PutUser(userID, lastname2, firstname2, email2, password2)
	if err2 != nil {
		t.Error(err2)
	}
	if user2.Activated != true {
		t.Error("bad activated value")
	}
	if user2.Id != userID {
		t.Error("bad user ID")
	}
	if !user2.CreatedAt.Equal(user.CreatedAt) {
		t.Error("bad created date")
	}
	if user2.UpdatedAt.Equal(user.UpdatedAt) {
		t.Error("bad updated date")
	}
	if user2.Username != username {
		t.Error("bad username")
	}
	if user2.Firstname != firstname2 {
		t.Error("bad firstname")
	}
	if user2.Lastname != lastname2 {
		t.Error("bad lastname")
	}
	if user2.Email != email2 {
		t.Error("bad email")
	}
	if user2.Role_id != role_id {
		t.Error("bad role_id")
	}
	err = bcrypt.CompareHashAndPassword(user2.Password, []byte(password2))
	if err != nil {
		t.Error("bad password")
	}
}

func TestDeleteUserNominal(t *testing.T) {
	username := "usernamedb5"
	lastname := "lastnamedb5"
	firstname := "firstnamedb5"
	email := "emaildb5"
	password := "passworddb5"
	role_id := 3
	user, err := sqlContext.CreateUser(username, lastname, firstname, email, password, role_id)
	if err != nil {
		t.Error(err)
	}
	userID := user.Id
	err2 := sqlContext.UnactivateUser(userID)
	if err2 != nil {
		t.Error(err2)
	}
	user2, err3 := sqlContext.GetUser(userID)
	if err3 != nil {
		t.Error(err3)
	}
	if user2.Activated {
		t.Error("the user is not delete")
	}
}
