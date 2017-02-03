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
	"github.com/pascallimeux/ocms2/modules/common"
	"github.com/pascallimeux/ocms2/modules/log"
	"github.com/pascallimeux/ocms2/setting"
	"os"
	"testing"
	"time"
)

var logfile *os.File
var consent_helper Consent_Helper
var config setting.Settings

const TransactionTimeout = 5000000000

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	// Init configs
	config = setting.Settings{Version: "1.0.1 (2017-02-01)", LogFileName: "/tmp/test.log", LogMode: "Trace", HttpHyperledger: "http://10.194.18.49:7050", ChainCodePath: "github.com/orangelabs/consent", ChainCodeName: "41149f5089e76b3b95fcf20e25a64fde7be07452b602dafe78f192bf826c14da25a3767c3a281255f4f355842a5a87a366aca65f44f5a02afc1784417079b76e", ApplicationID: "280399A20162908Z", EnrollID: "orange_user", EnrollSecret: "GtflmdhF6K32"}

	// Init logger
	logfile = log.Init_log(config.LogFileName, config.LogMode)

	// Init Hyperledger helpers
	HP_helper := HP_Helper{HttpHyperledger: config.HttpHyperledger, HLTimeout: config.HLTimeout}
	consent_helper = Consent_Helper{HP_helper: HP_helper, ChainCodePath: config.ChainCodePath, ChainCodeName: config.ChainCodeName, EnrollID: config.EnrollID, EnrollSecret: config.EnrollSecret}

}

func shutdown() {
	log.Trace(log.Here(), "End of tests..")
	defer logfile.Close()
}

func TestDeploySmartContractNominal(t *testing.T) {
	_, err := consent_helper.DeployConsentSM(config.ChainCodePath)
	if err != nil {
		t.Error(err)
	}
	t.Log("ChaincodeName: ", consent_helper.ChainCodeName)
}

func TestDeployBadSmartContractNominal(t *testing.T) {
	_, err := consent_helper.DeployConsentSM("Bad path")
	if err == nil {
		t.Error()
	}
	t.Log("Error:", err.Error())

}

func TestGetSmartContractVersionNominal(t *testing.T) {
	response, err := consent_helper.GetVersion()
	if err != nil {
		t.Error(err)
	}
	if !response.IsOK() {
		t.Error(response.GetError())
	}
	t.Log("ChaincodeVersion: ", response.GetMessage())
}

func TestRegistarNominal(t *testing.T) {
	response, err := consent_helper.Registar(config.EnrollID, config.EnrollSecret)
	if err != nil || !response {
		t.Error(err)
	}
	t.Log(response)
}

func TestBadRegistarNominal(t *testing.T) {
	_, err := consent_helper.Registar("bad user", config.EnrollSecret)
	if err == nil {
		t.Error()
	}
	t.Log(err.Error())
}

func TestIsRegistarTrueNominal(t *testing.T) {
	response, err := consent_helper.IsRegistar(config.EnrollID)
	if err != nil || !response {
		t.Error(err)
	}
	t.Log(response)
}

func TestIsRegistarFalseNominal(t *testing.T) {
	response, err := consent_helper.IsRegistar("badEnrollId")
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
	tr_uuid, err := consent_helper.CreateConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	t.Log("Transaction_uuid: ", tr_uuid)
}

func TestGetAllConsentsNominal(t *testing.T) {
	consents, err := consent_helper.GetAllConsents(config.ApplicationID)
	if err != nil {
		t.Error(err)
	}
	for _, consent := range consents {
		t.Log(consent.ToString())
	}
	t.Log("nb consents: ", len(consents))
}

func TestGetActivesConsentsNominal(t *testing.T) {
	consents, err := consent_helper.GetActivesConsents(config.ApplicationID)
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
	tr_uuid, err := consent_helper.CreateConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	t.Log("transaction_uuid: ", tr_uuid)
	time.Sleep(TransactionTimeout)
	consent, err2 := consent_helper.GetConsent(config.ApplicationID, tr_uuid)
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
	dt_begin := common.GetStringDateNow(0)
	dt_end := common.GetStringDateNow(1)
	_, err := consent_helper.CreateConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(TransactionTimeout)
	isconsent, err2 := consent_helper.IsConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess)
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
	_, err := consent_helper.CreateConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	isconsent, err2 := consent_helper.IsConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess)
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
	tr_uuid, err := consent_helper.CreateConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	t.Log("transaction_uuid: ", tr_uuid)
	time.Sleep(TransactionTimeout)
	response2, err2 := consent_helper.GetTRConsent(tr_uuid)
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
	dt_begin := common.GetStringDateNow(0)
	dt_end := common.GetStringDateNow(1)
	tr_uuid, err := consent_helper.CreateConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	t.Log("transaction_uuid: ", tr_uuid)
	time.Sleep(TransactionTimeout)
	consent1, err1 := consent_helper.GetConsent(config.ApplicationID, tr_uuid)
	if err1 != nil {
		t.Error(err1)
	}
	if consent1.State == "False" {
		t.Error()
	}
	_, err2 := consent_helper.UnactivateConsent(config.ApplicationID, tr_uuid)
	if err2 != nil {
		t.Error(err2)
	}
	time.Sleep(TransactionTimeout)
	consent2, err3 := consent_helper.GetConsent(config.ApplicationID, tr_uuid)
	if err3 == nil {
		t.Error(err3)
	}
}

func TestRemoveConsentsNominal(t *testing.T) {
	isOK, err := consent_helper.RemoveConsents()
	if err != nil {
		t.Error(err)
	}
	if !isOK {
		t.Error(err)
	}
	time.Sleep(TransactionTimeout)
	consents, err := consent_helper.GetAllConsents(config.ApplicationID)
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
	_, err := consent_helper.CreateConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	consumerID = "2222"
	_, err = consent_helper.CreateConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	consumerID = "3333"
	_, err = consent_helper.CreateConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(TransactionTimeout)
	consents, err2 := consent_helper.GetConsents4Owner(config.ApplicationID, ownerID)
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
	_, err := consent_helper.CreateConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	ownerID = "2222"
	_, err = consent_helper.CreateConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	ownerID = "3333"
	_, err = consent_helper.CreateConsent(config.ApplicationID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(TransactionTimeout)
	consents, err2 := consent_helper.GetConsents4Consumer(config.ApplicationID, consumerID)
	if err2 != nil {
		t.Error(err2)
	}
	if len(consents) != 3 {
		t.Error("3 expected but ", len(consents))
	}
}
