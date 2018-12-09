package api

import (
	"net/http"

	"usermanagement/src/db"
	"usermanagement/src/gocsv"
	"usermanagement/src/model"
)

// FileHandler /api/insert:
// collects POST fields like username and role. role can be premium or admin. default role if required to be set shall be default
func FileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// r.ParseForm()
		r.ParseMultipartForm(4096)
		if fileState, ok := r.PostForm["FileDownload"]; ok && fileState[0] == "Download" {
			// if err := db.InsertDB(w, r, r.FormValue("username"), r.FormValue("role")); err != nil {
			/*
				if err := db.InsertDB(w, r, r.FormValue("username"), r.FormValue("role")); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			*/
			// get data from db
			accounts, err := model.GetAccountList()
			if err != nil {
				returnJSON(w, err.Error(), http.StatusInternalServerError)
				return
			} else if len(accounts) < 1 {
				returnJSON(w, "DB: Empty database", http.StatusInternalServerError)
				return
			}

			// Set Trigger download of csv file
			w.Header().Set("Content-Description", "File Transfer")
			w.Header().Set("Content-Disposition", "attachment; filename=userdata.csv")

			// Marshal data into csv format and trigger download
			if err := gocsv.Marshal(accounts, w); err != nil {
				returnJSON(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Check for upload form
		} else if fileState, ok := r.PostForm["FileUpload"]; ok && fileState[0] == "Upload" {

			file, handle, err := r.FormFile("FileName")
			if err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				returnJSON(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			switch mimetype := handle.Header.Get("Content-Type"); mimetype {
			case "text/csv":
				var accounts []db.Account
				if err := gocsv.Unmarshal(file, &accounts); err != nil { // Load clients from file
					returnJSON(w, err.Error(), http.StatusInternalServerError)
				}

				for _, a := range accounts {
					if model.IsAccountExists(a.Mail) {
						returnJSON(w, "SQL: Account already exists", http.StatusInternalServerError)
					} else {
						model.CreateAccount(a.AccFailedLogins, a.AccIP, a.AccLastFailedLogin.Format("2006-01-02 15:04:05"), a.AccLastLogin.Format("2006-01-02 15:04:05"), a.AccStatus, a.AccountID, a.AddressSuffix, a.Applications, a.CompanyName, a.Country, a.Name, a.Mail, a.Phone, a.Postcode, a.SecondName, a.State, a.Surname, a.Town)
					}

					/*
						c := appengine.NewContext(r)
						_, k, err := db.AccountExists(c, a.Mail)
						if err != datastore.ErrNoSuchEntity && err != db.ErrAccountAlreadyExists {
							returnJSON(w, err.Error(), http.StatusInternalServerError)
						}

						if _, err := db.InsertDB(c, k, a); err != nil {
							returnJSON(w, err.Error(), http.StatusInternalServerError)
							// } else {
							// returnJSON(w, a, 0)
						}
					*/
				}
				returnJSON(w, "Success Inserted data into database", 0)
				// returnJSON(w, users, 0)
				return

			default:
				returnJSON(w, "Upload: file format not valid", http.StatusInternalServerError)
				return
			}

		} else {
			returnJSON(w, "Unable to parse form", http.StatusInternalServerError)
			return
		}
	}
}

/*
	b := &bytes.Buffer{}
	csvW := csv.NewWriter(b)

	var u db.User
	var record []string
	val := reflect.ValueOf(u)
	for i := 0; i < val.NumField(); i++ {
		record = append(record, val.Type().Field(i).Name)
	}
	if err := csvW.Write(record); err != nil {
		returnJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if users exist in db
	if len(users) > 1 {
		returnJSON(w, "Empty database", 1)
		return
	}

	// Add users in csv file
	for _, v := range users {
		var record []string
		val := reflect.Indirect(reflect.ValueOf(v))
		for i := 0; i < val.NumField(); i++ {
			record = append(record, val.Field(i).String())
		}
		// Write users to csv buffer
		if err := csvW.Write(record); err != nil {
			returnJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	csvW.Flush()
	if err := csvW.Error(); err != nil {
		returnJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
*/
