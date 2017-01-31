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
	"github.com/pascallimeux/ocms/common"
	"github.com/pascallimeux/ocms/hyperledger"
	"github.com/pascallimeux/ocms/hyperledger/consent"
	"github.com/pascallimeux/ocms/utils"
	"github.com/pascallimeux/ocms/utils/log"

	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup(true)
	code := m.Run()
	shutdown()
	os.Exit(code)
}

var AppContext api.AppContext
var httpServerTest *httptest.Server
var logfile *os.File
var DeployTimeout time.Duration
var TransactionTimeout time.Duration

func setup(isDropDB bool) {

	const (
		LOGFILENAME = "/var/log/ocms/ocmstest.log"
		LOGMODE     = "Trace"
	)

	// Read configuration file
	config_file := "../config/configtest.json"
	var configuration common.Configuration
	err := utils.Read_Conf(config_file, &configuration)
	if err != nil {
		panic(err.Error())
	}
	DeployTimeout = configuration.DeployTimeout
	TransactionTimeout = configuration.TransactionTimeout

	// Init logger
	logfile = log.Init_log(configuration.LogFileName, configuration.Logger)

	// Write configuration in log
	log.Info(log.Here(), utils.Get_fields(configuration))

	// Init Hyperledger helpers
	HP_helper := hyperledger.HP_Helper{HttpHyperledger: configuration.HttpHyperledger, HLTimeout: configuration.HLTimeout}
	Consent_Helper := consent.Consent_Helper{HP_helper: HP_helper, ChainCodePath: configuration.ChainCodePath, ChainCodeName: configuration.ChainCodeName, EnrollID: configuration.EnrollID, EnrollSecret: configuration.EnrollSecret}

	// Init application context
	AppContext = api.AppContext{Consent_helper: Consent_Helper, Configuration: configuration}

	// Init http server
	router := AppContext.CreateRoutes()
	httpServerTest = httptest.NewServer(router)

}

func shutdown() {
	log.Trace(log.Here(), "End of tests..")
	defer logfile.Close()
	defer httpServerTest.Close()
}
