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
	"testing"
	"time"
)

func TestDeploySmartContractNominal(t *testing.T) {
	_, err := AppContext.Consent_helper.DeployConsentSM(AppContext.Configuration.ChainCodePath)
	if err != nil {
		t.Error(err)
	}
	t.Log("ChaincodeName: ", AppContext.Consent_helper.ChainCodeName)
}

func TestDeployBadSmartContractNominal(t *testing.T) {
	_, err := AppContext.Consent_helper.DeployConsentSM("Bad path")
	if err == nil {
		t.Error()
	}
	t.Log("Error:", err.Error())

}

func TestGetSmartContractVersionNominal(t *testing.T) {
	response, err := AppContext.Consent_helper.GetVersion()
	if err != nil {
		t.Error(err)
	}
	if !response.IsOK() {
		t.Error(response.GetError())
	}
	t.Log("ChaincodeVersion: ", response.GetMessage())
}

func TestRegistarNominal(t *testing.T) {
	response, err := AppContext.Consent_helper.Registar(AppContext.Configuration.EnrollID, AppContext.Configuration.EnrollSecret)
	if err != nil || !response {
		t.Error(err)
	}
	t.Log(response)
}

func TestBadRegistarNominal(t *testing.T) {
	_, err := AppContext.Consent_helper.Registar("bad user", AppContext.Configuration.EnrollSecret)
	if err == nil {
		t.Error()
	}
	t.Log(err.Error())
}

func TestIsRegistarTrueNominal(t *testing.T) {
	response, err := AppContext.Consent_helper.IsRegistar(AppContext.Configuration.EnrollID)
	if err != nil || !response {
		t.Error(err)
	}
	t.Log(response)
}

func TestIsRegistarFalseNominal(t *testing.T) {
	response, err := AppContext.Consent_helper.IsRegistar("badEnrollId")
	if response {
		t.Error(err)
	}
	t.Log(err.Error())
}

func TestCreateConsentNominal(t *testing.T) {
	ownerID := "1111"
	consumerID := "2222"
	datatype := "BP"
	dataaccess := "R"
	dt_begin := "2016-09-04"
	dt_end := "2016-12-24"
	tr_uuid, err := AppContext.Consent_helper.CreateConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	t.Log("Transaction_uuid: ", tr_uuid)
}

func TestGetAllConsentsNominal(t *testing.T) {
	consents, err := AppContext.Consent_helper.GetAllConsents(AppContext.Configuration.ApplicationID)
	if err != nil {
		t.Error(err)
	}
	for _, consent := range consents {
		t.Log(consent.ToString())
	}
	t.Log("nb consents: ", len(consents))
}

func TestGetActivesConsentsNominal(t *testing.T) {
	consents, err := AppContext.Consent_helper.GetActivesConsents(AppContext.Configuration.ApplicationID)
	if err != nil {
		t.Error(err)
	}
	for _, consent := range consents {
		t.Log(consent.ToString())
	}
	t.Log("nb consents: ", len(consents))
}

func TestGetAConsentNominal(t *testing.T) {
	ownerID := "3333"
	consumerID := "44444"
	datatype := "BP"
	dataaccess := "R"
	dt_begin := "2016-09-04"
	dt_end := "2016-12-24"
	tr_uuid, err := AppContext.Consent_helper.CreateConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	t.Log("transaction_uuid: ", tr_uuid)
	time.Sleep(TransactionTimeout)
	consent, err2 := AppContext.Consent_helper.GetConsent(AppContext.Configuration.ApplicationID, tr_uuid)
	if err2 != nil {
		t.Error(err2)
	}
	t.Log("Consent: ", consent.ToString())
}

func TestIsConsentTrueNominal(t *testing.T) {
	ownerID := "5555"
	consumerID := "6666"
	datatype := "BP"
	dataaccess := "R"
	dt_begin := "2016-09-04"
	dt_end := "2016-12-24"
	_, err := AppContext.Consent_helper.CreateConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(TransactionTimeout)
	isconsent, err2 := AppContext.Consent_helper.IsConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess)
	if err2 != nil {
		t.Error(err2)
	}
	if !isconsent {
		t.Error()
	}
}

func TestIsConsentFalseNominal(t *testing.T) {
	ownerID := "7777"
	consumerID := "8888"
	datatype := "BP"
	dataaccess := "R"
	dt_begin := "2016-09-04"
	dt_end := "2016-12-24"
	_, err := AppContext.Consent_helper.CreateConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	isconsent, err2 := AppContext.Consent_helper.IsConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess)
	if err2 != nil {
		t.Error(err2)
	}
	if isconsent {
		t.Error()
	}
}

func TestGetTR4ConsentNominal(t *testing.T) {
	ownerID := "9999"
	consumerID := "1010"
	datatype := "BP"
	dataaccess := "R"
	dt_begin := "2016-09-04"
	dt_end := "2016-12-24"
	tr_uuid, err := AppContext.Consent_helper.CreateConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	t.Log("transaction_uuid: ", tr_uuid)
	time.Sleep(TransactionTimeout)
	response2, err2 := AppContext.Consent_helper.GetTRConsent(tr_uuid)
	if err2 != nil {
		t.Error(err2)
	}
	if !response2.IsOK() {
		t.Error(response2.GetError())
	}
	t.Log("Payload: ", response2.GetPayload())
}

func TestUnctivateConsentNominal(t *testing.T) {
	ownerID := "AAAA"
	consumerID := "BBBB"
	datatype := "BP"
	dataaccess := "R"
	dt_begin := "2016-09-04"
	dt_end := "2016-12-24"
	tr_uuid, err := AppContext.Consent_helper.CreateConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	t.Log("transaction_uuid: ", tr_uuid)
	time.Sleep(TransactionTimeout)
	consent1, err1 := AppContext.Consent_helper.GetConsent(AppContext.Configuration.ApplicationID, tr_uuid)
	if err1 != nil {
		t.Error(err1)
	}
	if consent1.State == "False" {
		t.Error()
	}
	_, err2 := AppContext.Consent_helper.UnactivateConsent(AppContext.Configuration.ApplicationID, tr_uuid)
	if err2 != nil {
		t.Error(err2)
	}
	time.Sleep(TransactionTimeout)
	consent2, err3 := AppContext.Consent_helper.GetConsent(AppContext.Configuration.ApplicationID, tr_uuid)
	if err3 != nil {
		t.Error(err3)
	}
	if consent2.State == "True" {
		t.Error()
	}
}

func TestRemoveConsentsNominal(t *testing.T) {
	isOK, err := AppContext.Consent_helper.RemoveConsents()
	if err != nil {
		t.Error(err)
	}
	if !isOK {
		t.Error(err)
	}
	time.Sleep(TransactionTimeout)
	consents, err := AppContext.Consent_helper.GetAllConsents(AppContext.Configuration.ApplicationID)
	if err != nil {
		t.Error(err)
	}
	if len(consents) != 0 {
		t.Error(err)
	}
}

func TestGetConsents4OwnerNominal(t *testing.T) {
	ownerID := "DDDD"
	consumerID := "1111"
	datatype := "BP"
	dataaccess := "R"
	dt_begin := "2016-09-04"
	dt_end := "2016-12-24"
	_, err := AppContext.Consent_helper.CreateConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	consumerID = "2222"
	_, err = AppContext.Consent_helper.CreateConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	consumerID = "3333"
	_, err = AppContext.Consent_helper.CreateConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(TransactionTimeout)
	consents, err2 := AppContext.Consent_helper.GetConsents4Owner(AppContext.Configuration.ApplicationID, ownerID)
	if err2 != nil {
		t.Error(err2)
	}
	if len(consents) != 3 {
		t.Error("3 expected but ", len(consents))
	}
}

func TestGetConsents4ConsumerNominal(t *testing.T) {
	ownerID := "1111"
	consumerID := "EEEE"
	datatype := "BP"
	dataaccess := "R"
	dt_begin := "2016-09-04"
	dt_end := "2016-12-24"
	_, err := AppContext.Consent_helper.CreateConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	ownerID = "2222"
	_, err = AppContext.Consent_helper.CreateConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	ownerID = "3333"
	_, err = AppContext.Consent_helper.CreateConsent(AppContext.Configuration.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(TransactionTimeout)
	consents, err2 := AppContext.Consent_helper.GetConsents4Consumer(AppContext.Configuration.ApplicationID, consumerID)
	if err2 != nil {
		t.Error(err2)
	}
	if len(consents) != 3 {
		t.Error("3 expected but ", len(consents))
	}
}
