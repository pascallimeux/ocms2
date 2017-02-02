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
	"time"
)

func (a *SqlContext) CreateLog(logs Logg) (Logg, error) {
	log.Trace(log.Here(), "CreateLog() : calling method -")
	sql := "insert or replace into logs (Timestamp, Resource_name, Resource_param, User_id, Username, Access_granted) values (?, ?, ?, ?, ?, ?) "

	stmt, err1 := a.Db.Prepare(sql)
	if err1 != nil {
		return logs, err1
	}
	defer stmt.Close()

	result, err2 := stmt.Exec(logs.Timestamp, logs.Resource_name, logs.Resource_param, logs.User_id, logs.Username, logs.Access_granted)
	if err2 != nil {
		return logs, err2
	}
	rowAffected, err3 := result.RowsAffected()
	if err3 != nil {
		return logs, err3
	}
	if rowAffected != 1 {
		return logs, errors.New("row not created")
	}
	return logs, nil
}

func (a *SqlContext) GetLogs() ([]Logg, error) {
	log.Trace(log.Here(), "GetLogs() : calling method -")
	sql := "select * from logs"
	var result = make([]Logg, 0)
	rows, err1 := a.Db.Query(sql)
	if err1 != nil {
		return result, err1
	}
	defer rows.Close()
	for rows.Next() {
		logs := Logg{}
		err2 := rows.Scan(&logs.Timestamp, &logs.Resource_name, &logs.Resource_param, &logs.User_id, &logs.Username, &logs.Access_granted)
		if err2 != nil {
			return result, err2
		}
		result = append(result, logs)
	}
	return result, nil
}

func (a *SqlContext) GetLog4Period(date1, date2 time.Time) ([]Logg, error) {
	log.Trace(log.Here(), "GetLog4Period() : calling method -")
	sql := "select * from logs where Timestamp > ? and Timestamp < ?"
	var result = make([]Logg, 0)

	stmt, err := a.Db.Prepare(sql)
	if err != nil {
		return result, err
	}
	defer stmt.Close()

	rows, err1 := stmt.Query(date1, date2)
	if err1 != nil {
		return result, err1
	}
	defer rows.Close()
	for rows.Next() {
		logs := Logg{}
		err2 := rows.Scan(&logs.Timestamp, &logs.Resource_name, &logs.Resource_param, &logs.User_id, &logs.Username, &logs.Access_granted)
		if err2 != nil {
			return result, err2
		}
		result = append(result, logs)
	}
	return result, nil
}
