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

package api

import (
	//"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pascallimeux/ocms/hyperledger"
	"github.com/pascallimeux/ocms/utils/log"
	"net/http"
	"strconv"
	"time"
)

//HTTP Post - /ocms/v1/api/consent
func (a *AppContext) processConsent(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "processConsent() : calling method -")
	var bytes []byte
	var consent Consent
	err := json.NewDecoder(r.Body).Decode(&consent)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	switch action := consent.Action; action {
	case "create":
		bytes, err = a.createConsent(consent)
	case "list":
		bytes, err = a.listConsents(consent.Appid)
	case "get":
		bytes, err = a.getConsent(consent.Appid, consent.Consentid)
	case "remove":
		bytes, err = a.unactivateConsent(consent.Appid, consent.Consentid)
	case "list4owner":
		bytes, err = a.getConsents4Owner(consent.Appid, consent.Ownerid)
	case "list4consumer":
		bytes, err = a.getConsents4Consumer(consent.Appid, consent.Consumerid)
	case "isconsent":
		bytes, err = a.isConsent(consent)
	default:
		log.Error(log.Here(), "bad action request")
		sendError(log.Here(), w, err)
		return
	}
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

//HTTP Get - /ocms/v1/api/hyperledger/consenttr/{truuid}
func (a *AppContext) processConsentTR(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tr_uuid := vars["truuid"]
	message := fmt.Sprintf("processConsentTR(tr_uuid=%s) : calling method -", tr_uuid)
	log.Info(log.Here(), message)
	var bytes []byte
	transaction, err := a.Consent_helper.HP_helper.GetTransaction(tr_uuid)
	if transaction.Txid == "" {
		sendError(log.Here(), w, errors.New("This is not a consent creation transaction"))
		return
	}
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	bytes, err1 := buildDecodedConsent(transaction)
	if err1 != nil {
		sendError(log.Here(), w, err1)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (a *AppContext) createConsent(consent Consent) ([]byte, error) {
	err := check_args(&consent)
	var message string
	if err != nil {
		message = fmt.Sprintf("createConsent(%s) : calling method -", err.Error())
	} else {
		message = fmt.Sprintf("createConsent(%s) : calling method -", consent.Print())
	}
	log.Info(log.Here(), message)
	if err != nil {
		return nil, err
	}
	consentID, err := a.Consent_helper.CreateConsent(a.Configuration.ApplicationID, consent.Ownerid, consent.Consumerid, consent.Datatype, consent.Dataaccess, consent.Dt_begin, consent.Dt_end)
	if err != nil {
		return nil, err
	}
	consent.Consentid = consentID
	return consent2Bytes(consent)
}

func (a *AppContext) listConsents(applicationID string) ([]byte, error) {
	message := fmt.Sprintf("listConsents(applicationID=%s) : calling method -", applicationID)
	log.Info(log.Here(), message)
	consents, err := a.Consent_helper.GetActivesConsents(a.Configuration.ApplicationID)
	if err != nil {
		return nil, err
	}
	return HPconsents2ConsentsBytes(consents)
}

func (a *AppContext) getConsent(applicationID, consentID string) ([]byte, error) {
	message := fmt.Sprintf("getConsent(applicationID=%s, consentID=%s) : calling method -", applicationID, consentID)
	log.Info(log.Here(), message)
	consent, err := a.Consent_helper.GetConsent(a.Configuration.ApplicationID, consentID)
	if err != nil {
		return nil, err
	}
	return HPconsent2ConsentBytes(consent)
}

func (a *AppContext) unactivateConsent(applicationID, consentID string) ([]byte, error) {
	message := fmt.Sprintf("unactivateConsent(applicationID=%s, consentID=%s) : calling method -", applicationID, consentID)
	log.Info(log.Here(), message)
	_, err := a.Consent_helper.UnactivateConsent(a.Configuration.ApplicationID, consentID)
	if err != nil {
		return nil, err
	}
	consent, err := a.Consent_helper.GetConsent(a.Configuration.ApplicationID, consentID)
	if err != nil {
		return nil, err
	}
	return HPconsent2ConsentBytes(consent)
}

func (a *AppContext) getConsents4Consumer(applicationID, consumerID string) ([]byte, error) {
	message := fmt.Sprintf("getConsents4Consumer(applicationID=%s, consumerID=%s) : calling method -", applicationID, consumerID)
	log.Info(log.Here(), message)
	consents, err := a.Consent_helper.GetConsents4Consumer(a.Configuration.ApplicationID, consumerID)
	if err != nil {
		return nil, err
	}
	return HPconsents2ConsentsBytes(consents)
}

func (a *AppContext) getConsents4Owner(applicationID, ownerID string) ([]byte, error) {
	message := fmt.Sprintf("getConsents4Owner(applicationID=%s, ownerID=%s) : calling method -", applicationID, ownerID)
	log.Info(log.Here(), message)
	consents, err := a.Consent_helper.GetConsents4Owner(a.Configuration.ApplicationID, ownerID)
	if err != nil {
		return nil, err
	}
	return HPconsents2ConsentsBytes(consents)
}

func (a *AppContext) isConsent(consent Consent) ([]byte, error) {
	message := fmt.Sprintf("isConsent(consent=%s) : calling method -", consent.Print())
	log.Info(log.Here(), message)
	isconsent, err := a.Consent_helper.IsConsent(a.Configuration.ApplicationID, consent.Ownerid, consent.Consumerid, consent.Datatype, consent.Dataaccess)
	if err != nil {
		return nil, err
	}
	response := IsConsent{}
	if isconsent {
		response.Consent = "True"
	} else {
		response.Consent = "False"
	}
	content, _ := json.Marshal(response)
	return content, nil
}

func convertHPConsents2APIConsents(HPconsents []hyperledger.Consent) []Consent {
	log.Trace(log.Here(), "convertHPConsents2APIConsents() : calling method -")
	consents := make([]Consent, len(HPconsents))
	for i, HPconsent := range HPconsents {
		consents[i] = convertHPConsent2APIConsent(HPconsent)
	}
	return consents
}

func convertHPConsent2APIConsent(HPconsent hyperledger.Consent) Consent {
	log.Trace(log.Here(), "convertHPConsent2APIConsent() : calling method -")
	consent := Consent{}
	consent.Consentid = HPconsent.ConsentID
	consent.Ownerid = HPconsent.OwnerID
	consent.Consumerid = HPconsent.ConsumerID
	consent.Dataaccess = HPconsent.Dataaccess
	consent.Datatype = HPconsent.Datatype
	consent.Dt_begin = HPconsent.Dt_begin
	consent.Dt_end = HPconsent.Dt_end
	return consent
}

func HPconsents2ConsentsBytes(HPconsents []hyperledger.Consent) ([]byte, error) {
	log.Trace(log.Here(), "HPconsents2ConsentsBytes() : calling method -")
	consents := convertHPConsents2APIConsents(HPconsents)
	j, err := json.Marshal(consents)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func HPconsent2ConsentBytes(HPconsent hyperledger.Consent) ([]byte, error) {
	log.Trace(log.Here(), "HPconsent2ConsentBytes() : calling method -")
	consent := convertHPConsent2APIConsent(HPconsent)
	return consent2Bytes(consent)
}

func consent2Bytes(consent Consent) ([]byte, error) {
	log.Trace(log.Here(), "consent2Bytes() : calling method -")
	j, err := json.Marshal(consent)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func check_args(consent *Consent) error {
	log.Trace(log.Here(), "check_args() : calling method -")
	if consent.Appid == "" {
		return errors.New("appID is mandatory!")
	}
	if consent.Ownerid == "" {
		return errors.New("ownerID is mandatory!")
	}
	if consent.Consumerid == "" {
		return errors.New("consumerID is mandatory!")
	}
	if consent.Dataaccess == "" {
		consent.Dataaccess = "A"
	}
	if consent.Datatype == "" {
		consent.Datatype = "All"
	}
	if consent.Dt_begin == "" {
		consent.Dt_begin = time.Now().Format("2006-01-02")
	}
	if consent.Dt_end == "" {
		consent.Dt_end = "2099-01-01"
	}
	return nil
}

func buildDecodedConsent(transaction hyperledger.Transaction) ([]byte, error) {
	log.Trace(log.Here(), "buildDecodedConsent() : calling method -")
	decConsent := DecodedConsent{}
	decConsent.Txuuid = transaction.Txid
	decConsent.Ccid = transaction.ChaincodeID
	decPayload, _ := b64.StdEncoding.DecodeString(transaction.Payload)
	//ccid := string(decPayload[11:139])
	args := make([]string, 9)
	j := 1
	i := 1
	decPayload = decPayload[140:]
	for i < len(decPayload) {
		if decPayload[i] == 10 {
			size, _ := strconv.Atoi(fmt.Sprintf("%v", decPayload[i+1]))
			args[j] = string(decPayload[i+2 : i+2+size])
			j++
			i = i + 2 + size
		} else {
			i++
		}
	}
	decConsent.Appid = string(fmt.Sprintf("%v", args[2]))
	decConsent.Ownerid = string(fmt.Sprintf("%v", args[3]))
	decConsent.Consumerid = string(fmt.Sprintf("%v", args[4]))
	decConsent.Datatype = string(fmt.Sprintf("%v", args[5]))
	decConsent.Dataaccess = string(fmt.Sprintf("%v", args[6]))
	decConsent.Dt_begin = string(fmt.Sprintf("%v", args[7]))
	decConsent.Dt_end = string(fmt.Sprintf("%v", args[8]))
	js, err := json.Marshal(decConsent)
	if err != nil {
		return nil, err
	}
	return js, nil
}
