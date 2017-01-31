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
	"github.com/pascallimeux/auth/modules/log"
	"net/http"
	"strings"
)

// Build and send http error
func sendError(from string, w http.ResponseWriter, err error) {
	log.Trace(log.Here(), "sendError() : calling method -")
	libelle := err.Error()
	libelle = strings.Replace(libelle, "\"", "'", -1)
	log.Error(from, "sendError: ", libelle)
	message := "{\"content\":\"" + libelle + "\"} "
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func buildHttp201Response(w http.ResponseWriter, data interface{}) {
	log.Trace(log.Here(), "buildHttp201Response() : calling method -")
	buildHttpResponse(w, data, http.StatusCreated)
}

func buildHttp200Response(w http.ResponseWriter, data interface{}) {
	log.Trace(log.Here(), "buildHttp200Response() : calling method -")
	buildHttpResponse(w, data, http.StatusOK)
}

func buildHttpResponse(w http.ResponseWriter, data interface{}, status int) {
	log.Trace(log.Here(), "buildHttpResponse() : calling method -")
	j, err := json.Marshal(data)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}

func buildHttp204Response(w http.ResponseWriter) {
	log.Trace(log.Here(), "buildHttp204Response() : calling method -")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func buildHttp401Response(w http.ResponseWriter) {
	log.Trace(log.Here(), "buildHttp401Response() : calling method -")
	error := "Not authorized request"
	message := "{\"content\":\"" + error + "\"} "
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(message))
}
