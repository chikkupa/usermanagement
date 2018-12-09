package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"usermanagement/src/db"
	"usermanagement/src/model"
)

// ValidateHandler /api/insert:
// collects POST fields like username and role. role can be premium or admin. default role if required to be set shall be default
func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var a db.Account
		err = json.Unmarshal(b, &a)
		if err != nil {
			http.Error(w, "\nUnmarshal err:\t"+err.Error(), http.StatusInternalServerError)
			return
		}

		if !model.IsAccountExists(a.Mail) {
			returnJSON(w, "SQL: No such entry", http.StatusInternalServerError)
			return
		}

		returnJSON(w, a, 0)

	} else {
		returnJSON(w, "Incorrect Method", http.StatusInternalServerError)
	}
	return
}
