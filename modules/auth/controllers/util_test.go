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
	"github.com/pascallimeux/ocms2/modules/auth/model"
	"github.com/pascallimeux/ocms2/modules/common"
	"net/http"
)

type HttpUser struct {
	Username  string `json:"username"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role_id   int    `json:"role_id"`
}

func createUser(tokenValue string, httpUser HttpUser) (model.User, int, error) {
	var user model.User
	data, _ := json.Marshal(httpUser)
	request, err1 := common.BuildRequestWithToken("POST", httpServerTest.URL+REGISTERURI, string(data), tokenValue)
	if err1 != nil {
		return user, 0, err1
	}
	status, body_bytes, err2 := common.ExecuteRequest(request)
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
	request, err := common.BuildRequestWithToken("GET", httpServerTest.URL+USERURI+"/"+userID, " ", tokenValue)
	if err != nil {
		return user, 0, err
	}
	status, body_bytes, err2 := common.ExecuteRequest(request)
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
	request, err := common.BuildRequestWithToken("GET", httpServerTest.URL+USERURI, " ", tokenValue)
	if err != nil {
		return users, 0, err
	}
	status, body_bytes, err2 := common.ExecuteRequest(request)
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
	request, err := common.BuildRequestWithToken("PUT", httpServerTest.URL+USERURI, data, tokenValue)
	if err != nil {
		return user, 0, err
	}
	status, body_bytes, err2 := common.ExecuteRequest(request)
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
	request, err := common.BuildRequestWithToken("DELETE", httpServerTest.URL+USERURI+"/"+userID, " ", tokenValue)
	if err != nil {
		return err
	}
	_, _, err2 := common.ExecuteRequest(request)
	if err2 != nil {
		return err2
	}
	return nil
}

func getToken(username, password string) (model.Token, error) {
	token := model.Token{}
	credentials := "{\"Username\":\"" + username + "\",\"Password\":\"" + password + "\"}"

	request, err1 := common.BuildRequest("POST", httpServerTest.URL+AUTHURI, credentials)
	if err1 != nil {
		return token, err1
	}
	status, body_bytes, err2 := common.ExecuteRequest(request)
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
