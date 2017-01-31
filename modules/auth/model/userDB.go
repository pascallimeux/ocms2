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
	"github.com/pascallimeux/auth/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (a *SqlContext) CreateUser(username, lastname, firstname, email, password string, role_id int) (User, error) {
	log.Trace(log.Here(), "CreateUser() : calling method -")
	sql := "insert into users (Id, Username, Lastname, Firstname, Email, Password, CreatedAt, UpdatedAt, Activated , Role_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	var user User
	user.Id = utils.Generate_uuid()
	user.Username = username
	user.Lastname = lastname
	user.Firstname = firstname
	user.Email = email
	user.Activated = true
	user.Role_id = role_id
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user.Password = hash
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	stmt, err1 := a.Db.Prepare(sql)
	if err1 != nil {
		return user, err1
	}
	defer stmt.Close()
	result, err2 := stmt.Exec(user.Id, user.Username, user.Lastname, user.Firstname, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.Activated, user.Role_id)
	if err2 != nil {
		return user, err2
	}
	rowAffected, err3 := result.RowsAffected()
	if err3 != nil {
		return user, err3
	}
	if rowAffected != 1 {
		return user, errors.New("row not created")
	}
	return user, nil
}

func (a *SqlContext) GetUserByCredentials(username, password string) (User, error) {
	log.Trace(log.Here(), "GetUserByCredentials(", username, ") : calling method -")
	sql := "select * from users where username = ?"
	var user User
	stmt, err := a.Db.Prepare(sql)
	if err != nil {
		return user, err
	}
	defer stmt.Close()
	err1 := stmt.QueryRow(username).Scan(&user.Id, &user.Username, &user.Lastname, &user.Firstname, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Activated, &user.Role_id)
	if err1 != nil {
		return user, err1
	}
	err2 := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err2 != nil {
		return user, err2
	}
	return user, nil
}

func (a *SqlContext) GetUsers() ([]User, error) {
	log.Trace(log.Here(), "GetUsers() : calling method -")
	sql := "select * from users"
	var result = make([]User, 0)

	rows, err1 := a.Db.Query(sql)
	if err1 != nil {
		return result, err1
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		err2 := rows.Scan(&user.Id, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Activated, &user.Role_id)
		if err2 != nil {
			return result, err2
		}
		user.Password = nil
		result = append(result, user)
	}
	return result, nil
}

func (a *SqlContext) GetUser(id string) (User, error) {
	log.Trace(log.Here(), "GetUser() : calling method -")
	sql := "select * from users where id = ?"
	var user User
	stmt, err := a.Db.Prepare(sql)
	if err != nil {
		return user, err
	}
	defer stmt.Close()
	err1 := stmt.QueryRow(id).Scan(&user.Id, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Activated, &user.Role_id)
	if err1 != nil {
		return user, err1
	}
	user.Password = nil
	return user, nil
}

func (a *SqlContext) PutUser(id, lastname, firstname, email, password string) (User, error) {
	log.Trace(log.Here(), "PutUser(id=", id, ") : calling method -")
	sql := "update users set Lastname = ?, Firstname = ?, Email = ?, Password = ?, UpdatedAt = ? where Id = ?"
	var newuser User
	stmt, err0 := a.Db.Prepare(sql)
	if err0 != nil {
		return newuser, err0
	}
	defer stmt.Close()

	newuser.Id = id
	newuser.Lastname = lastname
	newuser.Firstname = firstname
	newuser.Email = email
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return newuser, err
		}
		newuser.Password = hash
	}

	user, err1 := a.GetUser(newuser.Id)
	if err1 != nil {
		return user, err1
	}
	user.UpdatedAt = time.Now()

	if newuser.Lastname != "" {
		user.Lastname = newuser.Lastname
	}
	if newuser.Firstname != "" {
		user.Firstname = newuser.Firstname
	}
	if newuser.Email != "" {
		user.Email = newuser.Email
	}
	if newuser.Password != nil {
		user.Password = newuser.Password
	}

	_, err2 := stmt.Exec(user.Lastname, user.Firstname, user.Email, user.Password, user.UpdatedAt, user.Id)
	if err2 != nil {
		return user, err2
	}
	return user, nil
}

func (a *SqlContext) UnactivateUser(id string) error {
	log.Trace(log.Here(), "UnactivateUser() : calling method -")
	sql := "update users set Activated = ?, UpdatedAt = ?  where Id = ?"
	updatedAt := time.Now()
	activted := false
	stmt, err1 := a.Db.Prepare(sql)
	if err1 != nil {
		return err1
	}
	defer stmt.Close()
	_, err2 := stmt.Exec(activted, updatedAt, id)
	if err2 != nil {
		return err2
	}
	return nil
}
