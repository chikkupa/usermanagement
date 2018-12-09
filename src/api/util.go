package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"

	"usermanagement/src/db"
	"usermanagement/src/model"
)

// ResponseStruct returns a result with an exit cody
type ResponseStruct struct {
	Body    interface{}
	ErrCode int
}

var testUser = db.Account{
	AccStatus:    "default",
	AccLastLogin: time.Now(),
	Mail:         "test@example.com",
	Name:         "testUser",
	CompanyName:  "testCompany ltd",
	Country:      "Germany",
}

var testPremiumUser = db.Account{
	AccStatus:    "premium",
	AccLastLogin: time.Now(),
	Mail:         "premium@example.com",
	Name:         "premiumUser",
	CompanyName:  "premiumCompany ltd",
	Country:      "Germany",
}

var testAdminUser = db.Account{
	AccStatus:    "admin",
	AccLastLogin: time.Now(),
	Mail:         "admin@example.com",
	Name:         "adminUser",
	CompanyName:  "adminCompany ltd",
	Country:      "Germany",
}
var testUserSlice = []db.Account{testUser, testPremiumUser, testAdminUser}

var testapp01 = db.AppStruct{
	Applicationname: "testapp01",
	ApplicationID:   "testapp01",
	Namespace:       "testNameSpace",
	Groups:          "testGroups",
	Roles:           "testRoles",
	// Fields:          []string{"testFirlds"},
}
var testAppSlice = []db.AppStruct{testapp01}

// InsertTestUsers sets initial users
func InsertTestUsers(r *http.Request) {
	model.CreateAccountIfNotExixt()

	for _, a := range testUserSlice {
		if model.IsAccountExists(a.Mail) {
			return
		}

		model.CreateAccount(a.AccFailedLogins, a.AccIP, a.AccLastFailedLogin.Format("2006-01-02 15:04:05"), a.AccLastLogin.Format("2006-01-02 15:04:05"), a.AccStatus, a.AccountID, a.AddressSuffix, a.Applications, a.CompanyName, a.Country, a.Name, a.Mail, a.Phone, a.Postcode, a.SecondName, a.State, a.Surname, a.Town)
	}
}

// InsertTestApps sets initial users
func InsertTestApps(r *http.Request) {
	model.CreateApplicationIfNotExixt()

	for _, as := range testAppSlice {
		c := appengine.NewContext(r)
		// u, k, err := db.AccountExists(c, r.FormValue("UserName"))
		_, k, err := db.AppExists(c, as.ApplicationID)
		if err != datastore.ErrNoSuchEntity {
			return
		}
		datastore.Put(c, k, &as)
	}
}

func returnJSON(w http.ResponseWriter, body interface{}, errCode int) {
	responseStruct := ResponseStruct{
		Body:    body,
		ErrCode: errCode,
	}

	json.NewEncoder(w).Encode(responseStruct)
}

func insertIntoStruct(r *http.Request, a interface{}) (err error) {
	el := reflect.ValueOf(&a).Elem()
	// field, ok := reflect.TypeOf(u).Elem().FieldByName("UserName")
	for k, v := range r.Form {
		if el.Kind() == reflect.Struct {
			// exported field
			f := el.FieldByName(k)
			if f.IsValid() {
				if f.CanSet() {
					/*
						if f == "extra" {
							log.Println("extra value", reflect.ValueOf(v[0]))
						}
					*/
					f.Set(reflect.ValueOf(v[0]))
					//			fmt.Fprintf(w, "set: %v<br>", k)
				}
			} else {
				err = fmt.Errorf("InsertDB error: Unexpected field %v", k)
				//err = fmt.Errorf("InsertDB error: Unexpected field %v", k)
				// return
			}
		}

	}
	return err
}

// StringToSlice creates a slice from a string with comma as the delimiter
func StringToSlice(s string) []string {
	return strings.Split(s, ",")
}

func ReturnAppChildren(w http.ResponseWriter, r *http.Request) {
	// ctx := appengine.NewContext(r)
	// u := user.Current(ctx)

	appID := struct {
		ApplicationID string
	}{}

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(b, &appID)
	if err != nil {
		http.Error(w, "\nUnmarshal err:\t"+err.Error(), http.StatusInternalServerError)
		return
	}

	// _, appKey, err := db.AppExists(ctx, appID.ApplicationID)
	_, err = model.GetApplicationDetails(appID.ApplicationID)
	if err != nil && err != db.ErrAccountAlreadyExists {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cdSlc := []db.ChildData{}
	// query := datastore.NewQuery("applicationChild").Ancestor(appKey)
	// _, err = query.GetAll(ctx, &cdSlc)
	cdSlc = model.GetApplicationChildList(appID.ApplicationID)

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	data := struct {
		ChildData []db.ChildData
	}{
		ChildData: cdSlc,
	}

	returnJSON(w, data, 0)
}

func ReturnAppUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}

		data := struct {
			ApplicationID string
		}{}

		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(b, &data)
		if err != nil {
			http.Error(w, "\nUnmarshal err:\t"+err.Error(), http.StatusInternalServerError)
			return
		}

		// c := appengine.NewContext(r)
		// u, k, err := db.AccountExists(c, r.FormValue("UserName"))
		var accounts []db.Account
		accounts, err = model.GetAccountsConnectedToApplication(data.ApplicationID)
		// datastore.NewQuery("account").Filter("Application =", data.ApplicationID).GetAll(c, &accounts)
		returnJSON(w, accounts, 0)
	}

}
