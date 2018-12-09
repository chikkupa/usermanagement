package db

import (
	"fmt"
	"net/http"

	"appengine/datastore"

	"appengine"
)

// RetreveDB Stores user details into database
func RetreveDB(r *http.Request) (accounts []Account, err error) {
	c := appengine.NewContext(r)

	query := datastore.NewQuery("account")
	for it := query.Run(c); ; {
		// for {
		var account Account
		_, err := it.Next(&account)
		if err == datastore.Done {
			break
		}
		if err != nil {
			err = fmt.Errorf("Error fetching next user: %v", err)
			return accounts, err
		}
		accounts = append(accounts, account)
	}

	return
}
