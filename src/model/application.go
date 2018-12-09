package model

import (
	"database/sql"
	"log"
	"usermanagement/src/config"
	"usermanagement/src/db"
)

// CreateApplicationIfNotExixt Create table if not exist
func CreateApplicationIfNotExixt() error {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	stmt := `CREATE TABLE IF NOT EXISTS application(
		id int primary key auto_increment,
		application_id varchar(50),
		application_name varchar(100),
		groups varchar(100),
		namespace varchar(100),
		roles varchar(100)
	)`

	_, err = dbconn.Exec(stmt)

	if err != nil {
		log.Fatal("Error Creating Account: ", err)
	}
	return err
}

//CreateApplication Add a new Application
func CreateApplication(applicationID string, applicationName string, groups string, namespace string, roles string) {
	db, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := "INSERT INTO application (application_id, application_name, groups, namespace, roles) VALUES (?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(sql)
	_, err = stmt.Exec(applicationID, applicationName, groups, namespace, roles)
	if err != nil {
		log.Fatal(err)
	}
}

// IsApplicationExists Checks whether the account exist
func IsApplicationExists(applicationID string) bool {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)

	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	results, err := dbconn.Query("select id from application where application_id=? ", applicationID)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	if results.Next() {
		return true
	}

	return false
}

//GetApplicationDetails Get Application details
func GetApplicationDetails(applicationID string) (db.AppStruct, error) {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	sql := "SELECT application_id, application_name, groups, namespace, roles FROM application WHERE application_id=?"
	results, err := dbconn.Query(sql, applicationID)
	if err != nil {
		panic(err.Error())
	}
	var item db.AppStruct
	for results.Next() {
		err = results.Scan(&item.ApplicationID, &item.Applicationname, &item.Groups, &item.Namespace, &item.Roles)
		if err != nil {
			panic(err.Error())
		}

		err = db.ErrAccountAlreadyExists
	}
	return item, err
}

//GetApplicationList Get Application list
func GetApplicationList() ([]db.AppStruct, error) {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	sql := "SELECT application_id, application_name, groups, namespace, roles FROM application"
	results, err := dbconn.Query(sql)
	if err != nil {
		panic(err.Error())
	}

	var list []db.AppStruct
	for results.Next() {
		var item db.AppStruct
		err = results.Scan(&item.ApplicationID, &item.Applicationname, &item.Groups, &item.Namespace, &item.Roles)
		if err != nil {
			panic(err.Error())
		}
		list = append(list, item)
	}
	return list, err
}

// CreateApplicationChildIfNotExist Create table if not exist
func CreateApplicationChildIfNotExist() error {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	stmt := `CREATE TABLE IF NOT EXISTS application_child(
			id int primary key auto_increment,
			application_id varchar(50),
			field_name varchar(50),
			field_val varchar(100)
		)`

	_, err = dbconn.Exec(stmt)

	if err != nil {
		log.Fatal("Error Creating Account: ", err)
	}
	return err
}

//CreateApplicationChild Add a new ApplicationChild
func CreateApplicationChild(applicationID string, fieldName string, fieldVal string) {
	CreateApplicationChildIfNotExist()
	db, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := "INSERT INTO application_child (application_id, field_name, field_val) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(sql)
	_, err = stmt.Exec(applicationID, fieldName, fieldVal)
	if err != nil {
		log.Fatal(err)
	}
}

//GetApplicationChildList Get ApplicationChild list
func GetApplicationChildList(applicationID string) []db.ChildData {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	sql := "SELECT application_id, field_name, field_val FROM application_child WHERE application_id=?"
	results, err := dbconn.Query(sql, applicationID)
	if err != nil {
		panic(err.Error())
	}

	var list []db.ChildData
	for results.Next() {
		var item db.ChildData
		err = results.Scan(&item.ApplicationID, &item.FieldName, &item.FieldVal)
		if err != nil {
			panic(err.Error())
		}
		list = append(list, item)
	}
	return list
}
