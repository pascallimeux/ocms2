/*
Copyright Pascal Limeux. 2017 All Rights Reserved.
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

package main

import (
	"github.com/pascallimeux/ocms2/controllers"
	"github.com/pascallimeux/ocms2/hyperledger"
	auth "github.com/pascallimeux/ocms2/modules/auth/initialize"
	"github.com/pascallimeux/ocms2/modules/log"
	"github.com/pascallimeux/ocms2/setting"
	"net/http"
	"os"
	"time"
)

func main() {

	// Check command line parameters
	initDB := false
	args := os.Args[1:]
	if len(args) >= 1 && args[0] == "init" {
		initDB = true
	}

	// Init settings
	configuration, err := setting.GetSettings(".", "settings")
	if err != nil {
		panic(err.Error())
	}

	// Init logger
	f := log.Init_log(configuration.LogFileName, configuration.LogMode)
	defer f.Close()

	// Init Auth mode
	router, authContext, err2 := auth.Init(initDB, nil)
	if err2 != nil {
		panic(err2.Error())
	}
	defer authContext.SqlContext.Db.Close()

	// Init Hyperledger helpers
	HP_helper := hyperledger.HP_Helper{HttpHyperledger: configuration.HttpHyperledger, HLTimeout: configuration.HLTimeout}
	Consent_Helper := hyperledger.Consent_Helper{HP_helper: HP_helper, ChainCodePath: configuration.ChainCodePath, ChainCodeName: configuration.ChainCodeName, EnrollID: configuration.EnrollID, EnrollSecret: configuration.EnrollSecret}

	// Init application context
	appContext := controllers.AppContext{Consent_helper: Consent_Helper, Configuration: configuration, AuthContext: authContext}

	// Init permissions for application
	err3 := appContext.InitPermissions()
	if err3 != nil {
		panic(err3.Error())
	}

	// Init routes for application
	appContext.CreateOCMSRoutes(router)

	// Start http server
	log.Info(log.Here(), "Listening on: ", configuration.HttpHostUrl)
	s := &http.Server{
		Addr:         configuration.HttpHostUrl,
		Handler:      router,
		ReadTimeout:  configuration.ReadTimeout * time.Nanosecond,
		WriteTimeout: configuration.WriteTimeout * time.Nanosecond,
	}
	log.Fatal(log.Here(), s.ListenAndServe().Error())
}
