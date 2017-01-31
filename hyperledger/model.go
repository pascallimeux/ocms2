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

type Error struct {
	Code    int
	Message string
}

type Result struct {
	Status  string
	Message string
}

type Invoke struct {
	Jsonrpc string
	Method  string
	Params  Params
	Id      int
}

type Query struct {
	Jsonrpc string
	Method  string
	Params  Params
	Id      int
}

type Deploy struct {
	Jsonrpc string
	Method  string
	Params  Params
	Id      int
}

type ChaincodeID struct {
	Path string
	Name string
}

type CtorMsg struct {
	Function string
	Args     []string
}

type Params struct {
	Type          int
	ChaincodeID   ChaincodeID
	SecureContext string
	CtorMsg       CtorMsg
}

type Timestamp struct {
	Seconds int
	Nanos   int
}

/***********************************************************/
type Consent struct {
	AppID      string
	State      string
	ConsentID  string
	OwnerID    string
	ConsumerID string
	Datatype   string
	Dataaccess string
	Dt_begin   string
	Dt_end     string
}

func (c *Consent) ToString() string {
	var ret string
	ret = "AppID:" + c.AppID + " State:" + c.State + " ID:" + c.ConsentID + " OwnerID:" + c.OwnerID + " ConsumerID:" + c.ConsumerID + " Datatype:" + c.Datatype + " Dataaccess:" + c.Dataaccess + " Dt_begin:" + c.Dt_begin + " Dt_end:" + c.Dt_end
	return ret
}

/***********************************************************/
type SimpleResponse struct {
	OK    string
	Error string
}

func (s *SimpleResponse) IsOK() bool {
	if s.OK != "" && s.Error == "" {
		return true
	} else {
		return false
	}
}

func (s *SimpleResponse) GetError() string {
	return s.Error
}

/***********************************************************/
type Response struct {
	Jsonrpc string
	Result  Result
	Id      int
	Error   Error
}

func (r *Response) IsOK() bool {
	if r.Result.Status == "OK" && r.Error.Message == "" {
		return true
	} else {
		return false
	}
}

func (r *Response) GetError() string {
	return r.Error.Message
}

func (r *Response) GetMessage() string {
	return r.Result.Message
}

/***********************************************************/
type Transaction struct {
	Type        int
	ChaincodeID string
	Payload     string
	Txid        string
	Timestamp   Timestamp
	Nonce       string
	Cert        string
	Signature   string
	Error       string
}

func (t *Transaction) GetError() string {
	return t.Error
}

func (t *Transaction) GetPayload() string {
	return t.Payload
}

func (t *Transaction) IsOK() bool {
	if t.Error == "" {
		return true
	} else {
		return false
	}
}
