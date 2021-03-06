package main 

import (
	"net/http"
	"strings"
	"log"
	"fmt"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		//	not authenticated
		w.Header().Set("location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		//	some other error
		panic(err.Error())
	} else {
		//	Success - call the next handler
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler {next: handler}
}

//	loginHandler handles the third-party login process.
//	format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("Error when trying to get provider", provider, "-", err)
		}
		//	GetBeginAuthURL(1st arg, 2nd arg)
		//	1st arg - is a "state map of data" sending the user to a specific page after authantication.
		//	2nd arg - data provided to the auth provider that modifies the behavior of the process.
		loginURL, err := provider.GetBeginAuthURL(nil, nil)  
		if err != nil {
			log.Fatalln("Error when trying to GetBeginAuthURL for", provider, "-", err)
		}
		w.Header.Set("Location", loginURL)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
		//log.Println("TODO handle login for", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}