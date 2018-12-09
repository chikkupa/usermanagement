package db

import (
	"fmt"
	"time"
)

// User stores details of user and previlage level
type User struct {
	UserName string `json:"user_name"`
	Role     string `json:"role"`
	Age      string `json:"age"` // should be handled as int type
	// Extra    interface{}
	// Extra    map[string]string
	// Extra []Property
}

// Account describes the account details to be stored
type Account struct {
	AccountID          int       `json:"account_id" csv:"account_id"`
	AccStatus          string    `json:"acc_status" csv:"acc_status"`
	AccLastLogin       time.Time `json:"acc_last_login" csv:"acc_last_login"`
	AccLastFailedLogin time.Time `json:"acc_last_failed_login" csv:"acc_last_failed_login"`
	AccFailedLogins    int       `json:"acc_failed_logins" csv:"acc_failed_logins"`
	AccIP              string    `json:"acc_ip" csv:"acc_ip"`
	Mail               string    `json:"acc_mail" csv:"acc_mail"`
	Name               string    `json:"customer_name" csv:"customer_name"`
	Surname            string    `json:"customer_surname" csv:"customer_surname"`
	SecondName         string    `json:"customer_secondname" csv:"customer_secondname"`
	CompanyName        string    `json:"company_name" csv:"company_name"`
	Country            string    `json:"country" csv:"country"`
	State              string    `json:"state" csv:"state"`
	AddressSuffix      string    `json:"address_suffix" csv:"address_suffix"`
	Postcode           string    `json:"post_code" csv:"post_code"`
	Town               string    `json:"town" csv:"town"`
	Phone              string    `json:"phone" csv:"phone"`
	Applications       string    `json:"applications" csv:"applications"`
	// Extra              map[string]string `json:"extra" csv:"extra"`
}

type AppStruct struct {
	Applicationname string `json:"application_name" csv:"application_name"`
	ApplicationID   string `json:"application_id" csv:"application_id"`
	Namespace       string `json:"name_space" csv:"name_space"`
	Groups          string `json:"groups" csv:"groups"`
	Roles           string `json:"roles" csv:"roles"`
	// FieldName       string   `json:"fieldName" csv:"fieldName"`
	// FieldVal        []string `json:"fieldVal" csv:"fieldVal"`
	// Fields          []string `json:"fields" csv:"fields"` // changes into string so as to be able to use standard json unvarshal functionality
}

type MasterStruct struct {
	ID        string
	Account   Account
	AppStruct AppStruct
}

type ChildData struct {
	ApplicationID string
	FieldName     string
	FieldVal      string
}

// applicationFieldsMap := make(map[<ApplicationID>string]map[<fieldnr>int]<fieldname>string)

// Property contains rutime values in pairs
type Property struct {
	Name  string
	Value string
}

// ErrUserAlreadyExists sets standard error response
var ErrAccountAlreadyExists = fmt.Errorf("DB: Account already exists")

func checkStatus(status string) (err error) {

	if status == "premium" || status == "admin" || status == "default" {
	} else {
		err = fmt.Errorf("InsertDB err: Account Status definition unknown")
	}
	return
}
