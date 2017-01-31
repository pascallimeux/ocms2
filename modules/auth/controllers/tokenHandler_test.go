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
	"github.com/pascallimeux/auth/common"
	"github.com/pascallimeux/auth/model"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestCreateTokenNominal(t *testing.T) {
	token, err0 := getToken(common.ADMINLOGIN, common.ADMINPWD)
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
	token, err0 := getToken(common.ADMINLOGIN, common.ADMINPWD)
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

func getToken(username, password string) (model.Token, error) {
	token := model.Token{}
	credentials := "{\"Username\":\"" + username + "\",\"Password\":\"" + password + "\"}"

	request, err1 := buildRequest("POST", AUTHURI, credentials)
	if err1 != nil {
		return token, err1
	}
	status, body_bytes, err2 := ExecuteRequest(request)
	if err2 != nil {
		return token, err2
	}

	if status != http.StatusCreated {
		return token, errors.New("bad http status")
	}

	err3 := json.Unmarshal(body_bytes, &token)
	if err3 != nil {
		return token, err3
	}
	return token, nil
}

func createUser(tokenValue string, httpUser HttpUser) (model.User, int, error) {
	var user model.User
	data, _ := json.Marshal(httpUser)
	request, err1 := buildRequestWithToken("POST", REGISTERURI, string(data), tokenValue)
	if err1 != nil {
		return user, 0, err1
	}
	status, body_bytes, err2 := ExecuteRequest(request)
	if err2 != nil {
		return user, status, err2
	}
	err := json.Unmarshal(body_bytes, &user)
	if err != nil {
		return user, status, err
	}
	return user, status, nil
}

func getUser(tokenValue, userID string) (model.User, int, error) {
	var user model.User
	request, err := buildRequestWithToken("GET", USERURI+"/"+userID, " ", tokenValue)
	if err != nil {
		return user, 0, err
	}
	status, body_bytes, err2 := ExecuteRequest(request)
	if err2 != nil {
		return user, status, err2
	}
	err = json.Unmarshal(body_bytes, &user)
	if err != nil {
		return user, status, err
	}
	return user, status, nil
}

func getListOfUsers(tokenValue string) ([]model.User, int, error) {
	users := []model.User{}
	request, err := buildRequestWithToken("GET", USERURI, " ", tokenValue)
	if err != nil {
		return users, 0, err
	}
	status, body_bytes, err2 := ExecuteRequest(request)
	if err2 != nil {
		return users, status, err2
	}
	err = json.Unmarshal(body_bytes, &users)
	if err != nil {
		return users, status, err
	}
	return users, status, nil
}

func updateUser(tokenValue, data string) (model.User, int, error) {
	user := model.User{}
	request, err := buildRequestWithToken("PUT", USERURI, data, tokenValue)
	if err != nil {
		return user, 0, err
	}
	status, body_bytes, err2 := ExecuteRequest(request)
	if err2 != nil {
		return user, status, err2
	}
	err = json.Unmarshal(body_bytes, &user)
	if err != nil {
		return user, status, err
	}

	return user, status, nil
}

func deleteUser(tokenValue string, userID string) error {
	request, err := buildRequestWithToken("DELETE", USERURI+"/"+userID, " ", tokenValue)
	if err != nil {
		return err
	}
	_, _, err2 := ExecuteRequest(request)
	if err2 != nil {
		return err2
	}
	return nil
}

func buildRequest(method, uri, data string) (*http.Request, error) {
	var requestData *strings.Reader
	if data != "" {
		requestData = strings.NewReader(data)
	} else {
		requestData = nil
	}
	request, err := http.NewRequest(method, httpServerTest.URL+uri, requestData)
	if err != nil {
		return request, err
	}
	return request, nil
}

func buildRequestWithToken(method, uri, data, tokenValue string) (*http.Request, error) {
	request, err := buildRequest(method, uri, data)
	if err != nil {
		return request, err
	}
	request.Header.Set("authorization", "bearer "+tokenValue)
	return request, nil
}

func ExecuteRequest(request *http.Request) (int, []byte, error) {
	status := 0
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return status, nil, err
	}
	status = response.StatusCode
	body_bytes, err2 := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err2 != nil {
		return status, body_bytes, err2
	}
	return status, body_bytes, nil
}
