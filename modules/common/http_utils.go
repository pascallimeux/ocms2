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

package common

import (
	"encoding/json"
	"github.com/pascallimeux/ocms2/modules/log"
	"io/ioutil"
	"net/http"
	"strings"
)

// Build and send http error
func SendError(from string, w http.ResponseWriter, err error) {
	log.Trace(log.Here(), "sendError() : calling method -")
	libelle := err.Error()
	libelle = strings.Replace(libelle, "\"", "'", -1)
	log.Error(from, "sendError: ", libelle)
	message := "{\"content\":\"" + libelle + "\"} "
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func BuildHttp201Response(w http.ResponseWriter, data interface{}) {
	log.Trace(log.Here(), "buildHttp201Response() : calling method -")
	buildHttpResponse(w, data, http.StatusCreated)
}

func BuildHttp200Response(w http.ResponseWriter, data interface{}) {
	log.Trace(log.Here(), "buildHttp200Response() : calling method -")
	buildHttpResponse(w, data, http.StatusOK)
}

func buildHttpResponse(w http.ResponseWriter, data interface{}, status int) {
	log.Trace(log.Here(), "buildHttpResponse() : calling method -")
	j, err := json.Marshal(data)
	if err != nil {
		SendError(log.Here(), w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}

func BuildHttp204Response(w http.ResponseWriter) {
	log.Trace(log.Here(), "buildHttp204Response() : calling method -")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func BuildHttp401Response(w http.ResponseWriter) {
	log.Trace(log.Here(), "buildHttp401Response() : calling method -")
	error := "Not authorized request"
	message := "{\"content\":\"" + error + "\"} "
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(message))
}

func BuildRequest(method, uri, data string) (*http.Request, error) {
	var requestData *strings.Reader
	if data != "" {
		requestData = strings.NewReader(data)
	} else {
		requestData = nil
	}
	request, err := http.NewRequest(method, uri, requestData)
	if err != nil {
		return request, err
	}
	return request, nil
}

func BuildRequestWithToken(method, uri, data, tokenValue string) (*http.Request, error) {
	request, err := BuildRequest(method, uri, data)
	if err != nil {
		return request, err
	}
	request.Header.Set("authorization", "Bearer "+tokenValue)
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
