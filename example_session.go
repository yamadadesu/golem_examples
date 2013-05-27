package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/trevex/golem"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

const (
	secret      = "super-secret-key"
	sessionName = "golem.sid"
)

var store = sessions.NewCookieStore([]byte(secret))

func validateSession(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, sessionName)
	if v, ok := session.Values["isAuthorized"]; ok && v == true {
		fmt.Println("Authorized user identified!")
		return true
	} else {
		fmt.Println("Unauthorized user detected!")
		return false
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, sessionName)
	session.Values["isAuthorized"] = true
	session.Save(r, w)

	http.Redirect(w, r, "/example_session.html", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, sessionName)
	session.Values["isAuthorized"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/example_session.html", http.StatusFound)
}

func main() {
	flag.Parse()

	// Create a router
	myrouter := golem.NewRouter()
	myrouter.OnHandshake(validateSession)

	// Serve the public files
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Handle login
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	// Handle websockets using golems handler
	http.HandleFunc("/ws", myrouter.Handler())

	// Listen
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}