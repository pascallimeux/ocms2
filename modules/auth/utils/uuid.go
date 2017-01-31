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
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/pascallimeux/auth/modules/log"
)

func Generate_uuid() string {
	log.Trace(log.Here(), "Generate_uuid() : calling method -")
	b := make([]byte, 16)
	rand.Read(b)
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	log.Trace(log.Here(), "generated UID: ", uuid)
	return uuid
}

func Generate_Token() string {
	log.Trace(log.Here(), "Generate_Token() : calling method -")
	var dictionary = "01234567890abcde"
	var bytes = make([]byte, 40)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func object2Bytes(obj interface{}) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(obj)
	return b, err
}
