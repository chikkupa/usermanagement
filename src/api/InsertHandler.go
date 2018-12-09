package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"usermanagement/src/db"
	"usermanagement/src/model"
)

// InsertHandler /api/insert:
// collects POST fields like username and role. role can be premium or admin. default role if required to be set shall be default
func InsertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		var a db.Account

		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(b, &a)
		if err != nil {
			http.Error(w, "\nUnmarshal err:\t"+err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Account: ", a)

		err = model.CreateAccountIfNotExixt()
		if err != nil {
			log.Println(err.Error())
		}

		if r.URL.Path == "/api/insert" {
			if model.IsAccountExists(a.Mail) {
				returnJSON(w, "SQL: Account already exists", http.StatusInternalServerError)
				return
			}

			model.CreateAccount(a.AccFailedLogins, a.AccIP, a.AccLastFailedLogin.Format("2006-01-02 15:04:05"), a.AccLastLogin.Format("2006-01-02 15:04:05"), a.AccStatus, a.AccountID, a.AddressSuffix, a.Applications, a.CompanyName, a.Country, a.Name, a.Mail, a.Phone, a.Postcode, a.SecondName, a.State, a.Surname, a.Town)
		}

		if r.URL.Path == "/api/update" {
			if !model.IsAccountExists(a.Mail) {
				returnJSON(w, "SQL: No such user", http.StatusInternalServerError)
				return
			}

			model.UpdateAccount(a.AccFailedLogins, a.AccIP, a.AccLastFailedLogin.Format("2006-01-02 15:04:05"), a.AccLastLogin.Format("2006-01-02 15:04:05"), a.AccStatus, a.AccountID, a.AddressSuffix, a.Applications, a.CompanyName, a.Country, a.Name, a.Mail, a.Phone, a.Postcode, a.SecondName, a.State, a.Surname, a.Town)
		}

		returnJSON(w, a, 0)

	}
}

// CreateApplication creates and stores new applications into the datastore
func CreateApplication(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		var as db.AppStruct

		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(b, &as)
		if err != nil {
			http.Error(w, "\nUnmarshal err:\t"+err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(as)
		// c := appengine.NewContext(r)

		err = model.CreateApplicationIfNotExixt()
		if err != nil {
			log.Println(err.Error())
		}

		if model.IsApplicationExists(as.ApplicationID) {
			returnJSON(w, "SQL: Applicaton already exists", http.StatusInternalServerError)
			return
		}

		model.CreateApplication(as.ApplicationID, as.Applicationname, as.Groups, as.Namespace, as.Roles)
		returnJSON(w, as, 0)
		/*_, key, err := db.AppExists(c, as.ApplicationID)
		if err != nil && err != datastore.ErrNoSuchEntity {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err = datastore.Put(c, key, &as); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
			//		} else {
			//			returnJSON(w, as, 0)
		} */
	}
}

// AddExtraApplicationField adds extra children fields to the application
func AddApplicationChildData(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	childData := db.ChildData{}

	json.Unmarshal(b, &childData)

	// Check if application exists in db
	// c := appengine.NewContext(r)
	// _, key, err := db.AppExists(c, childData.ApplicationID)
	_, err = model.GetAccountDetails(childData.ApplicationID)
	if err != nil && err != db.ErrAccountAlreadyExists {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// adding extra fields as childeren of applicationStruct
	if childData.FieldName != "" && childData.FieldVal != "" {
		model.CreateApplicationChild(childData.ApplicationID, childData.FieldName, childData.FieldVal)
		// insert into db as child
		/*childKey := datastore.NewKey(c, "applicationChild", childData.FieldName, 0, key)
		if _, err = datastore.Put(c, childKey, &childData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(childData)*/

	}

	returnJSON(w, childData, 0)

}

// InsertApplicationIntoUserHandler adds allowed application details to a user
func InsertApplicationIntoUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}

		data := struct {
			ApplicationID string
			Mail          string
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
		// a, k, err := db.AccountExists(c, data.Mail)
		a, err := model.GetAccountDetails(data.Mail)
		if err != db.ErrAccountAlreadyExists {
			returnJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s := StringToSlice(a.Applications)
		for _, v := range s {
			if v == data.ApplicationID {
				returnJSON(w, "User Already has access to Application", 1)
				return
			}
		}

		if len(a.Applications) > 1 {
			a.Applications += ","
		}
		a.Applications += data.ApplicationID

		model.UpdateAccountApplications(a.Applications, a.Mail)
		returnJSON(w, a, 0)

		/*if _, err := db.InsertDB(c, k, a); err != nil {
			returnJSON(w, err.Error(), http.StatusInternalServerError)
		} else {
			returnJSON(w, a, 0)
		}*/
	}
}
