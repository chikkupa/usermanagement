package usermanagement

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"appengine"
	"appengine/datastore"
	"appengine/user"

	"usermanagement/src/api"
	"usermanagement/src/db"
	"usermanagement/src/model"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var sqldb *sql.DB

var (
	tmpl = template.Must(template.New("").Funcs(fmap).ParseGlob("tmpl/*"))

	fmap = template.FuncMap{
		"isPremium": isPremium,
		"isAdmin":   isAdmin,
		"isDefault": isDefault,
	}

	i = 0
)

// Functions for template processing
func isAdmin(accStatus string) bool {
	return accStatus == "admin"
}

func isPremium(accStatus string) bool {
	return accStatus == "premium"
}

func isDefault(accStatus string) bool {
	return accStatus == "default"
}

/* URL definition from app.yaml */
func init() {
	// Api routes
	http.HandleFunc("/api/insert", api.InsertHandler)
	http.HandleFunc("/api/update", api.InsertHandler)
	http.HandleFunc("/api/validate", api.ValidateHandler)
	http.HandleFunc("/api/file", api.FileHandler)
	http.HandleFunc("/api/loginURL", loginURL)
	http.HandleFunc("/api/createApplication", api.CreateApplication)
	http.HandleFunc("/api/updateApplication", api.CreateApplication)
	http.HandleFunc("/api/insertApplicationIntoUser", api.InsertApplicationIntoUserHandler)
	http.HandleFunc("/api/addApplicationChildData", api.AddApplicationChildData)
	http.HandleFunc("/api/ReturnAppChildren", api.ReturnAppChildren)
	http.HandleFunc("/api/ReturnAppUsers", api.ReturnAppUsers)

	http.HandleFunc("/adminPage", pagePath)
	http.HandleFunc("/premiumPage", pagePath)
	http.HandleFunc("/defaultPage", pagePath)
	http.HandleFunc("/insertPage", pagePath)
	http.HandleFunc("/updatePage", pagePath)
	http.HandleFunc("/validatePage", pagePath)
	http.HandleFunc("/downloadPage", pagePath)
	http.HandleFunc("/uploadPage", pagePath)
	http.HandleFunc("/profilePage", pagePath)
	http.HandleFunc("/appAdmin", appAdminHandler)

	// virtual test applications
	http.HandleFunc("/createApplication", pagePath)
	http.HandleFunc("/updateApplication", pagePath)
	http.HandleFunc("/testapp01", appPage)
	http.HandleFunc("/testapp02", appPage)
	http.HandleFunc("/testapp03", appPage)

	http.HandleFunc("/", static)
	// http.HandleFunc("/", indexHandler)

	//http.HandleFunc("/_ah/stop", shutdownHandler)

	//	if appengine.IsDevAppServer() {
	//		fragEntries = readFile("./static/fragnewentries.html")
	//	} else {
	//		fragEntries = "PROD"
	//	}

}

// collects user and passes it to the template
func pagePath(w http.ResponseWriter, r *http.Request) {
	model.CreateAccountIfNotExixt()
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	// a, _, err := db.AccountExists(ctx, u.Email)
	a, err := model.GetAccountDetails(u.Email)
	if err != nil && err != db.ErrAccountAlreadyExists && err != datastore.ErrNoSuchEntity {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")

	// you access the cached templates with the defined name, not the filename
	err = tmpl.ExecuteTemplate(w, path, a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Application pages
func appPage(w http.ResponseWriter, r *http.Request) {
	model.CreateApplicationIfNotExixt()
	model.CreateApplicationChildIfNotExist()
	// ctx := appengine.NewContext(r)
	// u := user.Current(ctx)

	path := strings.TrimPrefix(r.URL.Path, "/")
	// as, appKey, err := db.AppExists(ctx, path)
	as, err := model.GetApplicationDetails(path)
	if err != nil && err != db.ErrAccountAlreadyExists {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cdSlc := []db.ChildData{}
	//if err == db.ErrAccountAlreadyExists {
	// datastore.NewQuery("applicationChild").Ancestor(appKey).GetAll(ctx, &cdSlc)

	cdSlc = model.GetApplicationChildList(as.ApplicationID)

	data := struct {
		Application db.AppStruct
		ChildData   []db.ChildData
	}{
		Application: as,
		ChildData:   cdSlc,
	}
	// you access the cached templates with the defined name, not the filename
	err = tmpl.ExecuteTemplate(w, path, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Application pages
func appAdminHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	// u := user.Current(ctx)
	u := user.Current(ctx)
	// a, _, err := db.AccountExists(ctx, u.Email)
	a, err := model.GetAccountDetails(u.Email)
	if err != nil && err != db.ErrAccountAlreadyExists && err != datastore.ErrNoSuchEntity {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var apps []db.AppStruct
	// q := datastore.NewQuery("application")
	apps, err = model.GetApplicationList()
	// _, err = q.GetAll(ctx, &apps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		AccStatus    string
		Applications []db.AppStruct
	}{
		AccStatus:    a.AccStatus,
		Applications: apps,
	}

	err = tmpl.ExecuteTemplate(w, "appAdmin", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getApplications returns a slice of Applications data
func getApplications() []db.AppStruct {

	return nil
}

// static serves the static content of the project like a webserver
func static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/"+r.URL.Path)
	// var (
	// 	connectionName = os.Getenv("CLOUDSQL_CONNECTION_NAME")
	// 	user           = os.Getenv("CLOUDSQL_USER")
	// 	password       = os.Getenv("CLOUDSQL_PASSWORD") // NOTE: password may be empty
	// )

	// var err error
	// sqldb, err = sql.Open("gae-postgres", fmt.Sprintf("%s:%s@cloudsql(%s)/", user, password, connectionName))
	// if err != nil {
	// 	log.Fatalf("Could not open db: %v", err)
	// }

	// model.CreateTable(w)

	// log.Print("Successfully connected to database")
}

func loginURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if u == nil {
		url, _ := user.LoginURL(ctx, "/")
		fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
		return
	}
	if i < 1 {
		api.InsertTestUsers(r)
		api.InsertTestApps(r)
		i++
	}
	url, _ := user.LogoutURL(ctx, "/")
	fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
}

func profilePage(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	a, _, err := db.AccountExists(ctx, u.Email)
	if err != nil && err != db.ErrAccountAlreadyExists {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// you access the cached templates with the defined name, not the filename
	err = tmpl.ExecuteTemplate(w, "profilePage", a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if u == nil {
		url, _ := user.LoginURL(ctx, "/")
		fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
		return
	}
	url, _ := user.LogoutURL(ctx, "/")
	fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
}
*/
