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

package hyperledger

import (
	"encoding/json"
)

func Build_query_body(chaincode_name, hp_account, function string, args []string) ([]byte, error) {
	query := &Query{Jsonrpc: JSONRPC, Method: "query", Id: 1, Params: Params{Type: 1, ChaincodeID: ChaincodeID{Name: chaincode_name}, SecureContext: hp_account, CtorMsg: CtorMsg{Function: function, Args: args}}}
	bytes, err := json.Marshal(query)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

func Build_invoke_body(chaincode_name, hp_account, function string, args []string) ([]byte, error) {
	invoke := &Invoke{Jsonrpc: JSONRPC, Method: "invoke", Id: 1, Params: Params{Type: 1, ChaincodeID: ChaincodeID{Name: chaincode_name}, SecureContext: hp_account, CtorMsg: CtorMsg{Function: function, Args: args}}}
	bytes, err := json.Marshal(invoke)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

func Build_deploy_body(smartcontract_path, hp_account, function string, args []string) ([]byte, error) {
	deploy := &Deploy{Jsonrpc: JSONRPC, Method: "deploy", Id: 1, Params: Params{Type: 1, ChaincodeID: ChaincodeID{Path: smartcontract_path}, SecureContext: hp_account, CtorMsg: CtorMsg{Function: function, Args: args}}}
	bytes, err := json.Marshal(deploy)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

func Build_registar_body(enrollId, enrollSecret string) ([]byte, error) {
	bytes := []byte(`{"enrollId":` + "\"" + enrollId + "\"" + ` , "enrollSecret":` + "\"" + enrollSecret + "\"" + ` }`)
	return bytes, nil
}
