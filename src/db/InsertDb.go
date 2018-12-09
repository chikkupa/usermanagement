package db

import (
	"appengine/datastore"

	"appengine"
)

// InsertDB Stores user details into database
func InsertDB(c appengine.Context, key *datastore.Key, a Account) (success string, err error) {
	success = "DB: Processed Sucessfully"

	// Validateion before insert
	// more comprehensive setup is needed after finalizing this functionality
	err = checkStatus(a.AccStatus)
	if err != nil {
		return
	}

	_, err = datastore.Put(c, key, &a)
	return
}
