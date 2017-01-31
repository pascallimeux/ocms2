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
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Lastname  string    `json:"lastname"`
	Firstname string    `json:"firstname"`
	Email     string    `json:"email"`
	Password  []byte    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Activated bool      `json:"activated"`
	Role_id   int       `json:"role_id"` //foreign key with role table
}

type Token struct {
	Token      string    `json:"token"`
	Expires_in time.Time `json:"expire_in"`
	User_id    string    `json:"user_id"` // foreign key with user table
}

type Role struct {
	Code  int    `json:"code"`
	Label string `json:"label"`
}

type Permission struct {
	Resource_name string `json:"resource_name"`
	Role_code     int    `json:"role_code"` // foreign key with role table
	Owner_only    bool   `json:"owner_only"`
}

type Logg struct {
	CreatedAt      time.Time `json:"created_at"`
	Resource_name  string    `json:"resource_name"`
	Resource_param string    `json:"resource_param"`
	User_id        string    `json:"user_id"` // foreign key with user table
	Username       string    `json:"username"`
	Access_granted bool      `json:"access_granted"`
}

type SqlContext struct {
	dataSourceName string
	Db             *sql.DB
}

func GetSqlContext(dataSourceName string, initDB bool) (SqlContext, error) {
	var err error
	sqlContext := SqlContext{}
	sqlContext.dataSourceName = dataSourceName
	sqlContext.Db, err = GetSqlDB(dataSourceName)
	if err != nil {
		return sqlContext, err
	}
	if initDB {
		sqlContext.InitBDD()
	}
	return sqlContext, nil
}

func GetSqlDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return db, err
	}
	if db == nil {
		return db, errors.New("db nil")
	}
	return db, nil
}
