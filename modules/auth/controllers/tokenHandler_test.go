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
	"net/http"
	"testing"

	"github.com/pascallimeux/ocms2/modules/auth/setting"
)

func TestCreateTokenNominal(t *testing.T) {
	token, err0 := getToken(setting.ADMINLOGIN, setting.ADMINPWD)
	if err0 != nil {
		t.Error(err0)
	}
	username := "usernameA"
	password := "passwordA"
	_, _, err := createUser(token.Token, HttpUser{Username: username, Lastname: "lastnameA", Firstname: "firstnameA", Email: "emailA", Password: password, Role_id: 1})
	if err != nil {
		t.Error(err)
	}
	token, err2 := getToken(username, password)
	if err2 != nil {
		t.Error(err2)
	}
	jsonToken, _ := json.Marshal(token)
	t.Log(string(jsonToken))
}

func TestUserTokenNominal(t *testing.T) {
	username := "user1"
	password := "user_pwd1"
	email := "user1@orange.fr"
	role := 3
	token, err0 := getToken(setting.ADMINLOGIN, setting.ADMINPWD)
	if err0 != nil {
		t.Error(err0)
	}
	user, _, err := createUser(token.Token, HttpUser{Username: username, Lastname: "", Firstname: "", Email: email, Password: password, Role_id: role})
	if err != nil {
		t.Error(err)
	}

	token, err2 := getToken(username, password)
	if err2 != nil {
		t.Error(err2)
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
