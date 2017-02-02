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
	"github.com/pascallimeux/ocms2/modules/log"
	"strconv"
)

func (a *SqlContext) CreateRole(role Role) (Role, error) {
	log.Trace(log.Here(), "CreateRole(", role.Label, ") : calling method -")
	sql := "insert or replace into roles (Code, Label) values (?, ?) "

	stmt, err1 := a.Db.Prepare(sql)
	if err1 != nil {
		return role, err1
	}
	defer stmt.Close()

	result, err2 := stmt.Exec(role.Code, role.Label)
	if err2 != nil {
		return role, err2
	}
	rowAffected, err3 := result.RowsAffected()
	if err3 != nil {
		return role, err3
	}
	if rowAffected != 1 {
		return role, errors.New("row not created")
	}
	return role, nil
}

func (a *SqlContext) GetRole(code int) (Role, error) {
	log.Trace(log.Here(), "GetRole(", strconv.Itoa(code), ") : calling method -")
	sql := "select * from roles where code = ?"
	var role Role
	stmt, err := a.Db.Prepare(sql)
	if err != nil {
		return role, err
	}
	defer stmt.Close()
	err1 := stmt.QueryRow(code).Scan(&role.Code, &role.Label)
	if err1 != nil {
		return role, err1
	}
	return role, nil
}

func (a *SqlContext) GetRoles() ([]Role, error) {
	log.Trace(log.Here(), "GetRoles() : calling method -")
	sql := "select * from roles"
	var result = make([]Role, 0)
	rows, err1 := a.Db.Query(sql)
	if err1 != nil {
		return result, err1
	}
	defer rows.Close()
	for rows.Next() {
		role := Role{}
		err2 := rows.Scan(&role.Code, &role.Label)
		if err2 != nil {
			return result, err2
		}
		result = append(result, role)
	}
	return result, nil
}
