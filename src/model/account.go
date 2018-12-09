package model

import (
	"database/sql"
	"log"
	"os"
	"time"
	"usermanagement/src/config"
	"usermanagement/src/db"

	_ "github.com/go-sql-driver/mysql"
)

var dbconn *sql.DB

var (
	connectionName = os.Getenv("CLOUDSQL_CONNECTION_NAME")
	user           = os.Getenv("CLOUDSQL_USER")
	password       = os.Getenv("CLOUDSQL_PASSWORD") // NOTE: password may be empty
	database       = os.Getenv("CLOUDSQL_DATABASE")
)

// CreateAccountIfNotExixt Create table if not exist
func CreateAccountIfNotExixt() error {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	stmt := `CREATE TABLE IF NOT EXISTS account (
		id int not null primary key auto_increment,
		acc_failed_logins int,
		acc_ip varchar(30),
		acc_last_failed_logins datetime,
		acc_last_login datetime,
		acc_status varchar(30),
		account_id int,
		address_suffix text,
		applications text,
		company_name varchar(255),
		country varchar(100),
		mail varchar(255),
		name varchar(50),
		phone varchar(30),
		postcode varchar(30),
		second_name varchar(50),
		state varchar(50),
		surname varchar(50),
		town varchar(50)
	)`
	_, err = dbconn.Exec(stmt)

	if err != nil {
		log.Fatal("Error Creating Account: ", err)
	}
	return err
}

//CreateAccount Add a new Account
func CreateAccount(accFailedLogins int, accIP string, accLastFailedLogins string, accLastLogin string, accStatus string, accountID int, addressSuffix string, applications string, companyName string, country string, name string, mail string, phone string, postcode string, secondName string, state string, surname string, town string) {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	sql := "INSERT INTO account (acc_failed_logins, acc_ip, acc_last_failed_logins, acc_last_login, acc_status, account_id, address_suffix, applications, company_name, country, name, mail, phone, postcode, second_name, state, surname, town) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := dbconn.Prepare(sql)
	_, err = stmt.Exec(accFailedLogins, accIP, accLastFailedLogins, accLastLogin, accStatus, accountID, addressSuffix, applications, companyName, country, name, mail, phone, postcode, secondName, state, surname, town)
	if err != nil {
		log.Fatal(err)
	}
}

// IsAccountExists Checks whether the account exist
func IsAccountExists(mail string) bool {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)

	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	results, err := dbconn.Query("select id from account where mail=? ", mail)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	if results.Next() {
		return true
	}

	return false
}

//UpdateAccount Update a new Account
func UpdateAccount(accFailedLogins int, accIP string, accLastFailedLogins string, accLastLogin string, accStatus string, accountID int, addressSuffix string, applications string, companyName string, country string, name string, mail string, phone string, postcode string, secondName string, state string, surname string, town string) {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	sql := "UPDATE account SET acc_failed_logins=?, acc_ip=?, acc_last_failed_logins=?, acc_last_login=?, acc_status=?, account_id=?, address_suffix=?, applications=?, company_name=?, country=?, name=?, mail=?, phone=?, postcode=?, second_name=?, state=?, surname=?, town=? WHERE mail=?"
	stmt, err := dbconn.Prepare(sql)
	_, err = stmt.Exec(accFailedLogins, accIP, accLastFailedLogins, accLastLogin, accStatus, accountID, addressSuffix, applications, companyName, country, name, mail, phone, postcode, secondName, state, surname, town, mail)
	if err != nil {
		log.Fatal(err)
	}
}

// UpdateAccountApplications Update applications field of account
func UpdateAccountApplications(applications string, mail string) {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	sql := "UPDATE account SET applications=? WHERE mail=?"
	stmt, err := dbconn.Prepare(sql)
	_, err = stmt.Exec(applications, mail)
	if err != nil {
		log.Fatal(err)
	}
}

//GetAccountDetails Get Account details
func GetAccountDetails(mail string) (db.Account, error) {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	sql := "SELECT id, acc_failed_logins, acc_ip, acc_last_failed_logins, acc_last_login, acc_status, account_id, address_suffix, applications, company_name, country, name, mail, phone, postcode, second_name, state, surname, town FROM account WHERE mail=?"
	results, err := dbconn.Query(sql, mail)
	if err != nil {
		panic(err.Error())
	}
	var item db.Account
	for results.Next() {
		var accLastLogin string
		var accLastFailedLogin string
		err = results.Scan(&item.AccountID, &item.AccFailedLogins, &item.AccIP, &accLastFailedLogin, &accLastLogin, &item.AccStatus, &item.AccountID, &item.AddressSuffix, &item.Applications, &item.CompanyName, &item.Country, &item.Name, &item.Mail, &item.Phone, &item.Postcode, &item.SecondName, &item.State, &item.Surname, &item.Town)
		if err != nil {
			panic(err.Error())
		}

		layout := "2006-01-02 15:04:05"
		item.AccLastFailedLogin, err = time.Parse(layout, accLastFailedLogin)
		item.AccLastLogin, err = time.Parse(layout, accLastLogin)

		err = db.ErrAccountAlreadyExists
	}
	return item, err
}

//GetAccountList Get Account list
func GetAccountList() ([]db.Account, error) {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	sql := "SELECT id, acc_failed_logins, acc_ip, acc_last_failed_logins, acc_last_login, acc_status, account_id, address_suffix, applications, company_name, country, name, mail, phone, postcode, second_name, state, surname, town FROM account"
	results, err := dbconn.Query(sql)
	if err != nil {
		panic(err.Error())
	}

	var list []db.Account
	for results.Next() {
		var item db.Account
		var accLastLogin string
		var accLastFailedLogin string
		err = results.Scan(&item.AccountID, &item.AccFailedLogins, &item.AccIP, &accLastFailedLogin, &accLastLogin, &item.AccStatus, &item.AccountID, &item.AddressSuffix, &item.Applications, &item.CompanyName, &item.Country, &item.Name, &item.Mail, &item.Phone, &item.Postcode, &item.SecondName, &item.State, &item.Surname, &item.Town)
		if err != nil {
			panic(err.Error())
		}

		layout := "2006-01-02 15:04:05"
		item.AccLastFailedLogin, err = time.Parse(layout, accLastFailedLogin)
		item.AccLastLogin, err = time.Parse(layout, accLastLogin)

		list = append(list, item)
	}
	return list, err
}

// GetAccountsConnectedToApplication Get Accounts connected to application
func GetAccountsConnectedToApplication(applicationID string) ([]db.Account, error) {
	dbconn, err := sql.Open(config.Mysql, config.Dbconnection)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	sql := "SELECT id, acc_failed_logins, acc_ip, acc_last_failed_logins, acc_last_login, acc_status, account_id, address_suffix, applications, company_name, country, name, mail, phone, postcode, second_name, state, surname, town FROM account WHERE applications like '%" + applicationID + "%'"
	results, err := dbconn.Query(sql)
	if err != nil {
		panic(err.Error())
	}

	var list []db.Account
	for results.Next() {
		var item db.Account
		var accLastLogin string
		var accLastFailedLogin string
		err = results.Scan(&item.AccountID, &item.AccFailedLogins, &item.AccIP, &accLastFailedLogin, &accLastLogin, &item.AccStatus, &item.AccountID, &item.AddressSuffix, &item.Applications, &item.CompanyName, &item.Country, &item.Name, &item.Mail, &item.Phone, &item.Postcode, &item.SecondName, &item.State, &item.Surname, &item.Town)
		if err != nil {
			panic(err.Error())
		}

		layout := "2006-01-02 15:04:05"
		item.AccLastFailedLogin, err = time.Parse(layout, accLastFailedLogin)
		item.AccLastLogin, err = time.Parse(layout, accLastLogin)

		list = append(list, item)
	}
	return list, err
}
