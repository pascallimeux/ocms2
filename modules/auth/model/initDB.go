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
	_ "github.com/mattn/go-sqlite3"
	"github.com/pascallimeux/ocms2/modules/auth/setting"
	"github.com/pascallimeux/ocms2/modules/log"
)

func (s *SqlContext) InitBDD() {
	log.Trace(log.Here(), "InitBDD() : calling method -")
	s.DropTables()
	s.CreateTables()
	s.InitTables()
}

func (s *SqlContext) DropTables() {
	log.Trace(log.Here(), "DropTables() : calling method -")
	_, err := s.Db.Exec("DROP TABLE IF EXISTS users;")
	if err != nil {
		log.Fatal(log.Here(), "could not drop table:", err.Error())
	}
	_, err = s.Db.Exec("DROP TABLE IF EXISTS roles;")
	if err != nil {
		log.Fatal(log.Here(), "could not drop table:", err.Error())
	}
	_, err = s.Db.Exec("DROP TABLE IF EXISTS tokens;")
	if err != nil {
		log.Fatal(log.Here(), "could not drop table:", err.Error())
	}
	_, err = s.Db.Exec("DROP TABLE IF EXISTS permissions;")
	if err != nil {
		log.Fatal(log.Here(), "could not drop table:", err.Error())
	}
	_, err = s.Db.Exec("DROP TABLE IF EXISTS logs;")
	if err != nil {
		log.Fatal(log.Here(), "could not drop table:", err.Error())
	}
}

func (s *SqlContext) CreateTables() {
	log.Trace(log.Here(), "CreateTables() : calling method -")
	_, err := s.Db.Exec("create table if not exists users (Id varchar(255), Username varchar(50) not null UNIQUE, Firstname varchar(50), Lastname varchar(50), Email varchar(100), Password byte, CreatedAt datetime, UpdatedAt datetime, Activated boolean, Role_id integer, FOREIGN KEY(Role_id) REFERENCES roles(code), primary key (Id))")
	if err != nil {
		log.Fatal(log.Here(), err.Error())
	}
	_, err = s.Db.Exec("create table if not exists roles (Code integer, Label varchar(50), primary key (Code))")
	if err != nil {
		log.Fatal(log.Here(), err.Error())
	}
	_, err = s.Db.Exec("create table if not exists tokens (Token varchar(50), Expires_in datetime, User_id varchar(255), FOREIGN KEY(User_id) REFERENCES users(Id), primary key (Token))")
	if err != nil {
		log.Fatal(log.Here(), err.Error())
	}
	_, err = s.Db.Exec("create table if not exists permissions (Resource_name varchar(255), Role_code integer, Owner_only boolean, FOREIGN KEY(Role_code) REFERENCES roles(Code), primary key (Resource_name, Role_code))")
	if err != nil {
		log.Fatal(log.Here(), err.Error())
	}
	_, err = s.Db.Exec("create table if not exists logs (Timestamp datetime, Resource_name varchar(255), Resource_param varchar(255), User_id varchar(255), Username varchar(255), Access_granted boolean, FOREIGN KEY (User_id) REFERENCES users(Id) )")
	if err != nil {
		log.Fatal(log.Here(), err.Error())
	}
}

func (s *SqlContext) InitTables() {
	log.Trace(log.Here(), "InitTables() : calling method -")
	err := s.CreateInitialAdmin()
	if err != nil {
		log.Fatal(log.Here(), err.Error())
	}
	err = s.CreateInitialRoles()
	if err != nil {
		log.Fatal(log.Here(), err.Error())
	}
	err = s.CreateInitialPermissions()
	if err != nil {
		log.Fatal(log.Here(), err.Error())
	}
}

func (s *SqlContext) CreateInitialAdmin() error {
	log.Trace(log.Here(), "CreateInitialAdmin() : calling method -")
	_, err := s.CreateUser(setting.ADMINLOGIN, "", "", setting.ADMINEMAIL, setting.ADMINPWD, 1)
	if err != nil {
		log.Error(log.Here(), "initial admin already exist!!")
		return err
	}
	return nil
}

func (s *SqlContext) CreateInitialRoles() error {
	log.Trace(log.Here(), "CreateInitialRoles() : calling method -")
	none := Role{Code: 0, Label: "None"}
	admin := Role{Code: 1, Label: "Administrator"}
	superUser := Role{Code: 2, Label: "Superuser"}
	user := Role{Code: 3, Label: "user"}
	_, err := s.CreateRole(none)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return err
	}
	_, err = s.CreateRole(admin)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return err
	}
	_, err = s.CreateRole(superUser)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return err
	}
	_, err = s.CreateRole(user)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return err
	}
	return nil
}

func (s *SqlContext) CreateInitialPermissions() error {
	log.Trace(log.Here(), "CreateInitialPermissions() : calling method -")
	var perms []Permission

	perms = append(perms, Permission{Resource_name: "getToken", Role_code: 0, Owner_only: false})
	perms = append(perms, Permission{Resource_name: "getToken", Role_code: 1, Owner_only: false})
	perms = append(perms, Permission{Resource_name: "getToken", Role_code: 2, Owner_only: false})
	perms = append(perms, Permission{Resource_name: "getToken", Role_code: 3, Owner_only: false})

	perms = append(perms, Permission{Resource_name: "getRole", Role_code: 1, Owner_only: false})

	perms = append(perms, Permission{Resource_name: "getRoles", Role_code: 1, Owner_only: false})
	perms = append(perms, Permission{Resource_name: "postRole", Role_code: 1, Owner_only: false})

	perms = append(perms, Permission{Resource_name: "postRole", Role_code: 2, Owner_only: false})
	perms = append(perms, Permission{Resource_name: "postRole", Role_code: 3, Owner_only: false})

	perms = append(perms, Permission{Resource_name: "registerUser", Role_code: 1, Owner_only: false})

	perms = append(perms, Permission{Resource_name: "getUser", Role_code: 1, Owner_only: false})
	perms = append(perms, Permission{Resource_name: "getUser", Role_code: 2, Owner_only: true})
	perms = append(perms, Permission{Resource_name: "getUser", Role_code: 3, Owner_only: true})

	perms = append(perms, Permission{Resource_name: "getUsers", Role_code: 1, Owner_only: false})

	perms = append(perms, Permission{Resource_name: "updateUser", Role_code: 1, Owner_only: false})
	perms = append(perms, Permission{Resource_name: "updateUser", Role_code: 2, Owner_only: true})
	perms = append(perms, Permission{Resource_name: "updateUser", Role_code: 3, Owner_only: true})

	perms = append(perms, Permission{Resource_name: "deleteUser", Role_code: 1, Owner_only: false})

	perms = append(perms, Permission{Resource_name: "getLogs", Role_code: 1, Owner_only: false})

	for _, perm := range perms {
		_, err := s.CreatePermission(perm)
		if err != nil {
			log.Error(log.Here(), err.Error())
			return err
		}
	}
	return nil
}
