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

package controllers

import (
	"github.com/pascallimeux/ocms2/modules/auth/model"
	"github.com/pascallimeux/ocms2/modules/log"
)

func (a *AppContext) InitPermissions() error {
	log.Trace(log.Here(), "InitPermissions() : calling method -")
	var perms []model.Permission

	perms = append(perms, model.Permission{Resource_name: "processConsent", Role_code: 1, Owner_only: false})
	perms = append(perms, model.Permission{Resource_name: "processConsent", Role_code: 2, Owner_only: false})
	perms = append(perms, model.Permission{Resource_name: "processConsent", Role_code: 3, Owner_only: false})

	for _, perm := range perms {
		_, err := a.AuthContext.SqlContext.CreatePermission(perm)
		if err != nil {
			log.Error(log.Here(), err.Error())
			return err
		}
	}
	return nil
}
