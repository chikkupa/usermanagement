package pocTest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/urlfetch"
	"appengine/user"
)

// ipAddr nececary during testing so as to send data to service
var ipAddr = "192.168.43.66"

/* URL definition from app.yaml */
func init() {
	// Api routes
	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {
		serviceRequest(ipAddr, "api/insert", w, r)
	})
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		serviceRequest(ipAddr, "api/update", w, r)
	})
	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		serviceRequest(ipAddr, "api/validate", w, r)
	})
	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		serviceRequest(ipAddr, "api/file", w, r)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		u := user.Current(ctx)
		if u == nil {
			url, _ := user.LoginURL(ctx, "/")
			fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
			return
		}
		url, _ := user.LogoutURL(ctx, "/")
		fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
	})
	/*
			ctx := appengine.NewContext(r)
			// unable to make request via appengine as making outbound requests is only supported via urlfetch. and url fetch does not support creating a reverse proxy
			target, err := url.Parse(ipAddr + "/login")
			if err != nil {
				fmt.Fprint(w, err)
				return
			}
			rproxy := httputil.NewSingleHostReverseProxy(target)
			rproxy.ServeHTTP(w, r)
		})
	*/
	http.HandleFunc("/", static)
}

// TestHandler Sends data from form to service
func serviceRequest(ipAddr string, path string, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Send data to usermanagement service
	client := urlfetch.Client(ctx)
	resp, _ := client.Post(fmt.Sprintf("http://%v:8080/%v", ipAddr, path), "application/json", r.Body)

	// Collect response and send to client
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(bodyBytes))
}

// static serves the static content of the project like a webserver
func static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/"+r.URL.Path)
}
