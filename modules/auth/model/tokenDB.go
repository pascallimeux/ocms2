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
	"github.com/pascallimeux/ocms2/modules/common"
	"github.com/pascallimeux/ocms2/modules/log"
	"time"
)

func (a *SqlContext) CreateToken(user User, expire_in time.Duration) (interface{}, error) {
	log.Trace(log.Here(), "CreateToken() : calling method -")
	sql := "insert into tokens (Token, Expires_in, User_id) values (?, ?, ?)"
	type Token struct {
		Token      string
		Expires_in time.Time
	}
	var token Token
	token.Token = common.Generate_Token()
	token.Expires_in = time.Now().Add(expire_in * time.Hour)
	stmt, err1 := a.Db.Prepare(sql)
	if err1 != nil {
		return token, err1
	}
	defer stmt.Close()
	result, err2 := stmt.Exec(token.Token, token.Expires_in, user.Id)
	if err2 != nil {
		return token, err2
	}
	rowAffected, err3 := result.RowsAffected()
	if err3 != nil {
		return token, err3
	}
	if rowAffected != 1 {
		return token, errors.New("row not created")
	}
	return token, nil
}

func (a *SqlContext) GetToken(tokenvalue string) (Token, error) {
	log.Trace(log.Here(), "GetToken() : calling method -")
	sql := "select * from tokens where token = ?"
	var token Token
	stmt, err := a.Db.Prepare(sql)
	if err != nil {
		return token, err
	}
	defer stmt.Close()
	err1 := stmt.QueryRow(tokenvalue).Scan(&token.Token, &token.Expires_in, &token.User_id)
	if err1 != nil {
		return token, err1
	}
	return token, nil
}

func (a *Token) IsValid() bool {
	log.Trace(log.Here(), "IsValid() : calling method -")
	log.Trace(log.Here(), "token Expire_in:", a.Expires_in.String(), "  now:", time.Now().String())
	return a.Expires_in.After(time.Now())
}
