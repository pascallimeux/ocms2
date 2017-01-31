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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/pascallimeux/ocms/api"
)

func TestCreateConsentFromAPINominal(t *testing.T) {
	consent := api.Consent{Ownerid: "1111", Consumerid: "2222"}
	consentID, err := createConsent(consent)
	if err != nil {
		t.Error(err)
	}
	if consentID == "" {
		t.Error("bad consent ID")
	}
}

func TestGetConsentDetailFromAPINominal(t *testing.T) {
	consent := api.Consent{Ownerid: "OOOO", Consumerid: "AAAA"}
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
	createConsent(api.Consent{Ownerid: "1111", Consumerid: "2222"})
	createConsent(api.Consent{Ownerid: "1111", Consumerid: "3333"})
	createConsent(api.Consent{Ownerid: "1111", Consumerid: "4444"})
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
	createConsent(api.Consent{Ownerid: "1111", Consumerid: "2222"})
	createConsent(api.Consent{Ownerid: "1111", Consumerid: "3333"})
	createConsent(api.Consent{Ownerid: "1111", Consumerid: "4444"})
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
	createConsent(api.Consent{Ownerid: "1111", Consumerid: consumerid})
	createConsent(api.Consent{Ownerid: "2222", Consumerid: consumerid})
	createConsent(api.Consent{Ownerid: "3333", Consumerid: consumerid})
	consents, err := getListOfConsents("", consumerid)
	if err != nil {
		t.Error(err)
	}
	for _, consent := range consents {
		t.Log(consent.Print())
	}
}

func TestGetTRconsentFromAPINominal(t *testing.T) {
	consent := api.Consent{Ownerid: "1111", Consumerid: "2222"}
	consentID, err := createConsent(consent)
	time.Sleep(TransactionTimeout)
	decConsent := api.DecodedConsent{}
	if err != nil {
		t.Error(err)
	}
	url := httpServerTest.URL + api.CONSENTTRAPI + "/" + consentID + "f"
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

func consent2Bytes(consent api.Consent) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(consent)
	return b, err
}

func httpResponse2Consent(response *http.Response) (int, string, api.Consent, error) {
	status := response.StatusCode
	body := ""
	consent := api.Consent{}
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
	return status, body, consent, err
}

func httpResponse2Consents(response *http.Response) (int, string, []api.Consent, error) {
	status := response.StatusCode
	body := ""
	consents := []api.Consent{}
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

func createConsent(consent api.Consent) (string, error) {
	consent.Action = "create"
	consent.Appid = AppContext.Configuration.ApplicationID
	//dt_begin := time.Now().Format("2006-01-02")
	//dt_end := dt_begin.Add(time.Hour * 24 )
	bytes, _ := consent2Bytes(consent)
	response, err := http.Post(httpServerTest.URL+api.CONSENTAPI, "application/json", bytes)
	if err != nil {
		return "", err
	}
	statusCode, body, consent, err := httpResponse2Consent(response)
	if err != nil {
		return "", err
	}
	if statusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusCreated, response.StatusCode, body))
	}
	return consent.Consentid, nil
}

func getConsent(consentID string) (api.Consent, error) {
	consent := api.Consent{Action: "get", Appid: AppContext.Configuration.ApplicationID, Consentid: consentID}
	consent2 := api.Consent{}
	bytes, _ := consent2Bytes(consent)
	response, err := http.Post(httpServerTest.URL+api.CONSENTAPI, "application/json", bytes)
	if err != nil {
		return consent2, err
	}
	statusCode, body, consent2, err2 := httpResponse2Consent(response)
	if err2 != nil {
		return consent2, err2
	}
	if statusCode != http.StatusOK {
		return consent2, errors.New(fmt.Sprintf("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusCreated, response.StatusCode, body))
	}
	return consent2, nil
}

func getListOfConsents(ownerID, consumerID string) ([]api.Consent, error) {
	consent := api.Consent{Action: "list", Appid: AppContext.Configuration.ApplicationID}
	consents := []api.Consent{}
	if ownerID != "" {
		consent.Ownerid = ownerID
		consent.Action = "list4owner"
	} else if consumerID != "" {
		consent.Consumerid = consumerID
		consent.Action = "list4consumer"
	}
	bytes, _ := consent2Bytes(consent)
	response, err := http.Post(httpServerTest.URL+api.CONSENTAPI, "application/json", bytes)
	if err != nil {
		return consents, err
	}
	statusCode, body, consents, err := httpResponse2Consents(response)
	if err != nil {
		return consents, err
	}
	if statusCode != http.StatusOK {
		return consents, errors.New(fmt.Sprintf("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusCreated, response.StatusCode, body))
	}
	return consents, nil
}
