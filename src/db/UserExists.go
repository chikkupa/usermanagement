package db

import (
	"appengine"
	"appengine/datastore"
)

// AccountExists checks if a user exists
func AccountExists(c appengine.Context, aName string) (a Account, key *datastore.Key, err error) {
	key = datastore.NewKey(c, "account", aName, 0, nil)
	if err = datastore.Get(c, key, &a); err == nil {
		err = ErrAccountAlreadyExists
	}
	return
}

// Exists checks if a user exists
func AppExists(c appengine.Context, aName string) (a AppStruct, key *datastore.Key, err error) {

	key = datastore.NewKey(c, "application", aName, 0, nil)
	if err = datastore.Get(c, key, &a); err == nil {
		err = ErrAccountAlreadyExists
	}
	return
}
