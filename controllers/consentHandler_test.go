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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pascallimeux/ocms2/hyperledger"
	"github.com/pascallimeux/ocms2/model"
	authcontrollers "github.com/pascallimeux/ocms2/modules/auth/controllers"
	"github.com/pascallimeux/ocms2/modules/auth/initialize"
	authmodel "github.com/pascallimeux/ocms2/modules/auth/model"
	authsetting "github.com/pascallimeux/ocms2/modules/auth/setting"
	"github.com/pascallimeux/ocms2/modules/common"
	"github.com/pascallimeux/ocms2/modules/log"
	"github.com/pascallimeux/ocms2/setting"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var httpServerTest *httptest.Server
var logfile *os.File
var authContext authcontrollers.AppContext
var configuration setting.Settings
var tokenValue string

const TransactionTimeout = 5000000000

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	// Init configs
	authConfig := authsetting.Settings{DataSourceName: "/tmp/auth_test.db", LogFileName: "/tmp/test.log", LogMode: "Trace", ExpireInToken: 24}

	configuration = setting.Settings{Version: "1.0.1 (2017-02-01)", LogFileName: "/tmp/test.log", LogMode: "Trace", HttpHyperledger: "http://10.194.18.49:7050", ChainCodePath: "github.com/orangelabs/consent", ChainCodeName: "41149f5089e76b3b95fcf20e25a64fde7be07452b602dafe78f192bf826c14da25a3767c3a281255f4f355842a5a87a366aca65f44f5a02afc1784417079b76e", ApplicationID: "280399A20162908Z", EnrollID: "orange_user", EnrollSecret: "GtflmdhF6K32"}

	// Init logger
	logfile = log.Init_log(configuration.LogFileName, configuration.LogMode)

	// Init Auth mode
	var router *mux.Router
	var err error
	router, authContext, err = initialize.Init(true, &authConfig)
	if err != nil {
		panic(err.Error())
	}

	// Init Hyperledger helpers
	HP_helper := hyperledger.HP_Helper{HttpHyperledger: configuration.HttpHyperledger, HLTimeout: configuration.HLTimeout}
	Consent_Helper := hyperledger.Consent_Helper{HP_helper: HP_helper, ChainCodePath: configuration.ChainCodePath, ChainCodeName: configuration.ChainCodeName, EnrollID: configuration.EnrollID, EnrollSecret: configuration.EnrollSecret}

	// Init application context
	appContext := AppContext{Consent_helper: Consent_Helper, Configuration: configuration, AuthContext: authContext}

	// Init permissions for application
	err3 := appContext.InitPermissions()
	if err3 != nil {
		panic(err3.Error())
	}

	// Init routes for application
	appContext.CreateOCMSRoutes(router)

	// Init http server for tests
	httpServerTest = httptest.NewServer(router)

	// Get token
	token, err4 := getToken(authsetting.ADMINLOGIN, authsetting.ADMINPWD)
	if err4 != nil {
		panic(err4.Error())
	}
	tokenValue = token.Token
}

func shutdown() {
	log.Trace(log.Here(), "End of tests..")
	defer authContext.SqlContext.Db.Close()
	defer logfile.Close()
	defer httpServerTest.Close()
}

func TestCreateConsentFromAPINominal(t *testing.T) {
	consent := model.Consent{Ownerid: "1111", Consumerid: "2222"}
	consentID, err := createConsent(consent)
	if err != nil {
		t.Error(err)
	}
	if consentID == "" {
		t.Error("bad consent ID")
	}
}

func TestGetConsentDetailFromAPINominal(t *testing.T) {
	consent := model.Consent{Ownerid: "OOOO", Consumerid: "AAAA"}
	consentID, err := createConsent(consent)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(TransactionTimeout)
	consent2, err2 := getConsent(consentID)
	if err2 != nil {
		t.Error(err2)
	}
	if consent2.Consentid != consentID || consent2.Consumerid != consent.Consumerid {
		t.Error(err)
	}

}

func TestGetConsentsFromAPINominal(t *testing.T) {
	createConsent(model.Consent{Ownerid: "1111", Consumerid: "2222"})
	createConsent(model.Consent{Ownerid: "1111", Consumerid: "3333"})
	createConsent(model.Consent{Ownerid: "1111", Consumerid: "4444"})
	consents, err := getListOfConsents("", "")
	if err != nil {
		t.Error(err)
	}
	for _, consent := range consents {
		t.Log(consent.Print())
	}
}

func TestGetConsents4OwnerFromAPINominal(t *testing.T) {
	ownerid := "1111"
	createConsent(model.Consent{Ownerid: "1111", Consumerid: "2222"})
	createConsent(model.Consent{Ownerid: "1111", Consumerid: "3333"})
	createConsent(model.Consent{Ownerid: "1111", Consumerid: "4444"})
	consents, err := getListOfConsents(ownerid, "")
	if err != nil {
		t.Error(err)
	}
	for _, consent := range consents {
		t.Log(consent.Print())
	}
}

func TestGetConsents4ConsumerFromAPINominal(t *testing.T) {
	consumerid := "3333"
	createConsent(model.Consent{Ownerid: "1111", Consumerid: consumerid})
	createConsent(model.Consent{Ownerid: "2222", Consumerid: consumerid})
	createConsent(model.Consent{Ownerid: "3333", Consumerid: consumerid})
	consents, err := getListOfConsents("", consumerid)
	if err != nil {
		t.Error(err)
	}
	for _, consent := range consents {
		t.Log(consent.Print())
	}
}

func TestGetTRconsentFromAPINominal(t *testing.T) {
	consent := model.Consent{Ownerid: "1111", Consumerid: "2222"}
	consentID, err := createConsent(consent)
	time.Sleep(TransactionTimeout)
	decConsent := model.DecodedConsent{}
	if err != nil {
		t.Error(err)
	}
	url := httpServerTest.URL + CONSENTTRAPI + "/" + consentID + "f"
	t.Log("TR consent for ", consentID, "whith url:", url)
	response, err2 := http.Get(url)
	if err2 != nil {
		t.Error(err2)
	}
	rec_bytes, err3 := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err3 != nil {
		t.Error(err3)
	}
	t.Log(string(rec_bytes))
	err4 := json.Unmarshal(rec_bytes, &decConsent)
	if err4 != nil {
		t.Error(err4)
	}
	t.Log(decConsent.Ccid)
}

func consent2BufferBytes(consent model.Consent) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(consent)
	return b, err
}

func httpResponse2Consent(response *http.Response) (int, string, model.Consent, error) {
	status := response.StatusCode
	body := ""
	consent := model.Consent{}
	rec_bytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return status, body, consent, err
	}
	body = string(rec_bytes)
	err = json.Unmarshal(rec_bytes, &consent)
	if err != nil {
		return status, body, consent, err
	}
	fmt.Println("Status: ", status)
	fmt.Println("Body: ", body)

	return status, body, consent, err
}

func httpResponse2Consents(response *http.Response) (int, string, []model.Consent, error) {
	status := response.StatusCode
	body := ""
	consents := []model.Consent{}
	rec_bytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return status, body, consents, err
	}
	body = string(rec_bytes)
	err = json.Unmarshal(rec_bytes, &consents)
	if err != nil {
		return status, body, consents, err
	}
	return status, body, consents, err
}

func createConsent(consent model.Consent) (string, error) {
	var responseConsent model.Consent
	consent.Action = "create"
	consent.Appid = configuration.ApplicationID
	data, _ := json.Marshal(consent)
	request, err1 := common.BuildRequestWithToken("POST", httpServerTest.URL+CONSENTAPI, string(data), tokenValue)
	if err1 != nil {
		return "", err1
	}
	status, body_bytes, err2 := common.ExecuteRequest(request)
	if err2 != nil {
		return "", err2
	}
	err3 := json.Unmarshal(body_bytes, &responseConsent)
	if err3 != nil {
		return "", err3
	}

	if status != http.StatusOK {
		return "", errors.New("bad status")
	}
	return responseConsent.Consentid, nil
}

func getConsent(consentID string) (model.Consent, error) {
	consent := model.Consent{Action: "get", Appid: configuration.ApplicationID, Consentid: consentID}
	responseConsent := model.Consent{}
	data, _ := json.Marshal(consent)
	request, err1 := common.BuildRequestWithToken("POST", httpServerTest.URL+CONSENTAPI, string(data), tokenValue)
	if err1 != nil {
		return responseConsent, err1
	}
	status, body_bytes, err2 := common.ExecuteRequest(request)
	if err2 != nil {
		return responseConsent, err2
	}
	err3 := json.Unmarshal(body_bytes, &responseConsent)
	if err3 != nil {
		return responseConsent, err3
	}
	if status != http.StatusOK {
		return responseConsent, errors.New("bad status")
	}
	return responseConsent, nil
}

func getListOfConsents(ownerID, consumerID string) ([]model.Consent, error) {
	consent := model.Consent{Action: "list", Appid: configuration.ApplicationID}
	consents := []model.Consent{}
	if ownerID != "" {
		consent.Ownerid = ownerID
		consent.Action = "list4owner"
	} else if consumerID != "" {
		consent.Consumerid = consumerID
		consent.Action = "list4consumer"
	}
	data, _ := json.Marshal(consent)
	request, err1 := common.BuildRequestWithToken("POST", httpServerTest.URL+CONSENTAPI, string(data), tokenValue)
	if err1 != nil {
		return consents, err1
	}
	status, body_bytes, err2 := common.ExecuteRequest(request)
	if err2 != nil {
		return consents, err2
	}
	err3 := json.Unmarshal(body_bytes, &consents)
	if err3 != nil {
		return consents, err3
	}
	if status != http.StatusOK {
		return consents, errors.New("bad status")
	}
	return consents, nil
}

func getToken(username, password string) (authmodel.Token, error) {
	token := authmodel.Token{}
	credentials := "{\"Username\":\"" + username + "\",\"Password\":\"" + password + "\"}"

	request, err1 := common.BuildRequest("POST", httpServerTest.URL+authcontrollers.AUTHURI, credentials)
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
