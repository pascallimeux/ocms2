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

package consent

import (
	"encoding/json"
	"errors"
	"github.com/pascallimeux/ocms/hyperledger"
	"github.com/pascallimeux/ocms/utils/log"
	"strings"
)

type Consent_Helper struct {
	HP_helper     hyperledger.HP_Helper
	ChainCodePath string
	ChainCodeName string
	EnrollID      string
	EnrollSecret  string
}

func (c *Consent_Helper) DeployConsentSM(chainCodePath string) (string, error) {
	log.Trace(log.Here(), "DeployConsentSM() : calling method -")
	chaincodename := ""

	function := "init"
	args := make([]string, 0)
	response, err := c.HP_helper.DeployChainCode(chainCodePath, c.EnrollID, function, args)

	if err != nil {
		log.Error(log.Here(), "Deploy ConsentSM error : ", err.Error())
		return chaincodename, err
	}
	if !response.IsOK() {
		err = errors.New(response.GetError())
		log.Error(log.Here(), "Deploy ConsentSM error : ", err.Error())
		return chaincodename, err
	}
	chaincodename = response.GetMessage()
	c.ChainCodeName = chaincodename
	log.Info(log.Here(), "Retrieve a new chaincodename : ", c.ChainCodeName, " for ChainCodePath: ", c.ChainCodePath)
	return chaincodename, err
}

func (c *Consent_Helper) GetVersion() (hyperledger.Response, error) {
	log.Trace(log.Here(), "GetVersion() : calling method -")
	function := "GetVersion"
	args := make([]string, 0)
	response, err := c.HP_helper.Query(c.ChainCodeName, c.EnrollID, function, args)
	return response, err
}

func (c *Consent_Helper) Registar(enrollID, enrollSecret string) (bool, error) {
	log.Trace(log.Here(), "Registar() : calling method -")
	response, err := c.HP_helper.Registar(enrollID, enrollSecret)
	if err != nil {
		return false, err
	}
	if response.IsOK() {
		return true, nil
	} else {
		err = errors.New(response.GetError())
		return false, err
	}
}

func (c *Consent_Helper) IsRegistar(enrollID string) (bool, error) {
	log.Trace(log.Here(), "IsRegistar() : calling method -")
	response, err := c.HP_helper.IsRegistar(enrollID)
	if err != nil {
		return false, err
	}
	if response.IsOK() {
		return true, nil
	} else {
		err = errors.New(response.GetError())
		return false, err
	}
}

func (c *Consent_Helper) GetTRConsent(consentID string) (hyperledger.Transaction, error) {
	log.Trace(log.Here(), "GetTRConsent() : calling method -")
	response, err := c.HP_helper.GetTransaction(consentID)
	return response, err
}

func (c *Consent_Helper) CreateConsent(appID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end string) (string, error) {
	log.Trace(log.Here(), "CreateConsent() : calling method -")
	function := "PostConsent"
	args := []string{appID, ownerID, consumerID, datatype, dataaccess, dt_begin, dt_end}
	response, err := c.HP_helper.Invoke(c.ChainCodeName, c.EnrollID, function, args)
	return extractTransactionUUID(response, err)
}

func (c *Consent_Helper) GetConsent(appID, consentID string) (hyperledger.Consent, error) {
	log.Trace(log.Here(), "GetConsent(", consentID, ") : calling method -")
	function := "GetConsent"
	args := []string{appID, consentID}
	response, err := c.HP_helper.Query(c.ChainCodeName, c.EnrollID, function, args)
	return extractConsent(response, err)
}

func (c *Consent_Helper) GetAllConsents(appID string) ([]hyperledger.Consent, error) {
	log.Trace(log.Here(), "GetAllConsents() : calling method -")
	function := "GetConsents"
	args := []string{appID, "ALLM"}
	response, err := c.HP_helper.Query(c.ChainCodeName, c.EnrollID, function, args)
	return extractConsents(response, err)
}

func (c *Consent_Helper) GetActivesConsents(appID string) ([]hyperledger.Consent, error) {
	log.Trace(log.Here(), "GetActivesConsents() : calling method -")
	function := "GetConsents"
	args := []string{appID}
	response, err := c.HP_helper.Query(c.ChainCodeName, c.EnrollID, function, args)
	return extractConsents(response, err)
}

func (c *Consent_Helper) GetConsents4Consumer(appID, consumerID string) ([]hyperledger.Consent, error) {
	log.Trace(log.Here(), "GetConsents4Consumer() : calling method -")
	function := "GetConsumerConsents"
	args := []string{appID, consumerID}
	response, err := c.HP_helper.Query(c.ChainCodeName, c.EnrollID, function, args)
	return extractConsents(response, err)
}

func (c *Consent_Helper) GetConsents4Owner(appID, ownerID string) ([]hyperledger.Consent, error) {
	log.Trace(log.Here(), "GetConsents4Owner() : calling method -")
	function := "GetOwnerConsents"
	args := []string{appID, ownerID}
	response, err := c.HP_helper.Query(c.ChainCodeName, c.EnrollID, function, args)
	return extractConsents(response, err)
}

func (c *Consent_Helper) IsConsent(appID, ownerID, consumerID, datatype, dataaccess string) (bool, error) {
	log.Trace(log.Here(), "IsConsent() : calling method -")
	function := "IsConsent"
	args := []string{appID, ownerID, consumerID, datatype, dataaccess}
	response, err := c.HP_helper.Query(c.ChainCodeName, c.EnrollID, function, args)
	return extractIsConsent(response, err)
}

func (c *Consent_Helper) RemoveConsents() (bool, error) {
	log.Trace(log.Here(), "RemoveConsents() : calling method -")
	function := "Reset"
	args := []string{}
	response, err := c.HP_helper.Invoke(c.ChainCodeName, c.EnrollID, function, args)
	return extractStatusOK(response, err)
}

func (c *Consent_Helper) UnactivateConsent(appID, consentID string) (hyperledger.Response, error) {
	log.Trace(log.Here(), "UnactivateConsent() : calling method -")
	function := "RemoveConsent"
	args := []string{appID, consentID}
	response, err := c.HP_helper.Invoke(c.ChainCodeName, c.EnrollID, function, args)
	return response, err
}

func extractConsent(response hyperledger.Response, err error) (hyperledger.Consent, error) {
	var consent hyperledger.Consent
	if err != nil {
		return consent, err
	}
	dec := json.NewDecoder(strings.NewReader(response.GetMessage()))
	err = dec.Decode(&consent)
	return consent, err
}

func extractConsents(response hyperledger.Response, err error) ([]hyperledger.Consent, error) {
	var consents []hyperledger.Consent
	if err != nil {
		return consents, err
	}
	dec := json.NewDecoder(strings.NewReader(response.GetMessage()))
	err = dec.Decode(&consents)
	return consents, err
}

func extractTransactionUUID(response hyperledger.Response, err error) (string, error) {
	var uuid string
	if err != nil {
		return uuid, err
	}
	uuid = response.GetMessage()
	return uuid, err
}

func extractIsConsent(response hyperledger.Response, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	if response.GetMessage() == "True" {
		return true, nil
	} else {
		return false, nil
	}
}

func extractStatusOK(response hyperledger.Response, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	if response.IsOK() {
		return true, nil
	} else {
		err = errors.New(response.GetError())
		return false, err
	}
}
