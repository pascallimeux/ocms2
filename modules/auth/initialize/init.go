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

package initialize

import (
	"github.com/gorilla/mux"
	"github.com/pascallimeux/ocms2/modules/auth/controllers"
	"github.com/pascallimeux/ocms2/modules/auth/model"
	"github.com/pascallimeux/ocms2/modules/auth/setting"
)

func Init(initDB bool, configuration *setting.Settings) (*mux.Router, controllers.AppContext, error) {

	if configuration == nil {
		var err error
		// Init settings
		configuration, err = setting.GetSettings(".", "authsettings")
		if err != nil {
			panic(err.Error())
		}
	}

	var router *mux.Router
	var appContext controllers.AppContext

	// Init sqliteDB
	sqlContext, err := model.GetSqlContext(configuration.DataSourceName, initDB)
	if err != nil {
		return router, appContext, err
	}

	// Init application context
	appContext = controllers.AppContext{Settings: configuration, SqlContext: sqlContext}

	// Init http server
	router = appContext.CreateAUTHRoutes()

	return router, appContext, nil
}
