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

package tests

import (
	"github.com/pascallimeux/ocms/api"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestGetAPIVersionNominal(t *testing.T) {
	res, err := http.Get(httpServerTest.URL + api.VERSIONURI)
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	body := string(data)
	if res.StatusCode != http.StatusOK {
		t.Fatal("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusCreated, res.StatusCode, body)
	}
	t.Log(body)
	if !strings.Contains(body, "{\"Version\"") {
		t.Fatalf("Non-expected body content: %v", body)
	}
}
