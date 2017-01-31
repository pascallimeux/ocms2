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

package model

import (
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pascallimeux/auth/modules/log"
	"strconv"
	"time"
)

func (a *SqlContext) CreatePermission(permission Permission) (Permission, error) {
	log.Trace(log.Here(), "CreatePermission() : calling method -")
	sql := "insert or replace into permissions (Resource_name, Role_code, Owner_only) values (?, ?, ?) "

	stmt, err1 := a.Db.Prepare(sql)
	if err1 != nil {
		return permission, err1
	}
	defer stmt.Close()

	result, err2 := stmt.Exec(permission.Resource_name, permission.Role_code, permission.Owner_only)
	if err2 != nil {
		return permission, err2
	}
	rowAffected, err3 := result.RowsAffected()
	if err3 != nil {
		return permission, err3
	}
	if rowAffected != 1 {
		return permission, errors.New("row not created")
	}
	return permission, nil
}

func (a *SqlContext) GetPermission(resource_name string, role_code int) (Permission, error) {
	log.Trace(log.Here(), "GetPermission(", strconv.Itoa(role_code), "  ", resource_name, ") : calling method -")
	sql := "select * from permissions where resource_name = ? and role_code = ?"
	var permission Permission
	stmt, err := a.Db.Prepare(sql)
	if err != nil {
		return permission, err
	}
	defer stmt.Close()
	err1 := stmt.QueryRow(resource_name, role_code).Scan(&permission.Resource_name, &permission.Role_code, &permission.Owner_only)
	if err1 != nil {
		return permission, err1
	}
	return permission, nil
}

func (a *SqlContext) GetPermissions() ([]Permission, error) {
	log.Trace(log.Here(), "GetPermissions() : calling method -")
	sql := "select * from permissions"
	var result = make([]Permission, 0)
	rows, err1 := a.Db.Query(sql)
	if err1 != nil {
		return result, err1
	}
	defer rows.Close()
	for rows.Next() {
		permission := Permission{}
		err2 := rows.Scan(&permission.Resource_name, &permission.Role_code, &permission.Owner_only)
		if err2 != nil {
			return result, err2
		}
		result = append(result, permission)
	}
	return result, nil
}

func (a *SqlContext) IsAuthorized4Token(tokenValue, resourceName, resourceId string) error {
	log.Trace(log.Here(), "IsAuthorized4Token() : calling method -")

	token, err2 := a.GetToken(tokenValue)
	if err2 != nil {
		log.Trace(log.Here(), "get token (", tokenValue, ") ", err2.Error())
		return err2
	}
	if token.IsValid() == false {
		log.Trace(log.Here(), "invalide token ")
		return errors.New("Invalid token!")
	}
	user, err3 := a.GetUser(token.User_id)
	if err3 != nil {
		log.Trace(log.Here(), "get user ", err3.Error())
		return err3
	}
	if user.Activated == false {
		log.Trace(log.Here(), "user not activated ")
		return errors.New("Unactivated user!")
	}
	userRole, _ := a.GetRole(user.Role_id)
	authorized := a.IsPermitted4User(userRole, user.Id, user.Username, resourceName, resourceId)
	if authorized {
		log.Trace(log.Here(), "The user: ", user.Username, " is granted to access to the resource: ", resourceName)
		return nil
	} else {
		log.Trace(log.Here(), "The user: ", user.Username, " is not authorized to access to the resource: ", resourceName)
		return errors.New("User not authorized for this resource!")
	}
}

func (a *SqlContext) IsPermitted4User(userRole Role, userId, username, resourceName, resourceParam string) bool {
	log.Trace(log.Here(), "IsPermitted4User() : calling method -")
	permission, err := a.GetPermission(resourceName, userRole.Code)
	if resourceParam == "" {
		resourceParam = "None"
	}
	var granted = false
	if err != nil {
		log.Error(log.Here(), err.Error())
	} else {
		if !permission.Owner_only {
			granted = true
		} else {
			if resourceParam == userId {
				granted = true
			}
		}
	}
	logs := Logg{CreatedAt: time.Now(), Resource_name: resourceName, Resource_param: resourceParam, User_id: userId, Username: username, Access_granted: granted}
	_, err = a.CreateLog(logs)
	if err != nil {
		log.Error(log.Here(), err.Error())
	}
	return granted
}
