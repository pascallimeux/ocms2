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
	"github.com/pascallimeux/ocms2/modules/auth/initialize"
	"github.com/pascallimeux/ocms2/modules/auth/setting"
	"github.com/pascallimeux/ocms2/modules/log"
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

	configuration, err := setting.GetSettings(".", "auth")
	if err != nil {
		panic(err.Error())
	}

	// Init application
	router, authContext, err2 := initialize.Init(initDB, configuration)
	if err2 != nil {
		panic(err2.Error())
	}
	defer authContext.SqlContext.Db.Close()

	log.Info(log.Here(), "Listening on: ", configuration.HttpHostUrl)
	s := &http.Server{
		Addr:         configuration.HttpHostUrl,
		Handler:      router,
		ReadTimeout:  configuration.ReadTimeout * time.Nanosecond,
		WriteTimeout: configuration.WriteTimeout * time.Nanosecond,
	}
	log.Fatal(log.Here(), s.ListenAndServe().Error())
}
