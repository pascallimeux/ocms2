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

package main

import (
	"github.com/pascallimeux/auth/controllers"
	"github.com/pascallimeux/auth/model"
	"github.com/pascallimeux/auth/modules/log"
	"github.com/pascallimeux/auth/modules/setting"
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

	// Init application context
	appContext := controllers.AppContext{Settings: configuration}

	// Init sqliteDB
	sqlContext, err := model.GetSqlContext(configuration.DataSourceName, initDB)
	if err != nil {
		log.Fatal(log.Here(), err.Error())
	}
	defer sqlContext.Db.Close()

	// Init and Start http server
	router := appContext.CreateRoutes()
	log.Info(log.Here(), "Listening on: ", configuration.HttpHostUrl)
	s := &http.Server{
		Addr:         configuration.HttpHostUrl,
		Handler:      router,
		ReadTimeout:  configuration.ReadTimeout * time.Nanosecond,
		WriteTimeout: configuration.WriteTimeout * time.Nanosecond,
	}
	log.Fatal(log.Here(), s.ListenAndServe().Error())
}
