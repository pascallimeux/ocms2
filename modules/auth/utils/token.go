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

package utils

import (
	"encoding/json"
	"github.com/pascallimeux/ocms2/modules/log"
	"net/http"
	"strings"
)

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
