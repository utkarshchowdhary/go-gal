package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

func Recoverer(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				http.Error(w, "Aaaah! Something went wrong", http.StatusInternalServerError)
			}
		}()
		next(w, r, p)
	}
}

func Authenticator(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if RequestUser(r) == nil {
			query := url.Values{}
			query.Add("redir", r.URL.String())
			http.Redirect(w, r, "/sign-in?"+query.Encode(), http.StatusFound)
			return
		}
		next(w, r, p)
	}
}
