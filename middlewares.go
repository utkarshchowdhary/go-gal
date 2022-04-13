package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func Neuter(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func Authenticator(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if RequestUser(r) == nil {
			query := url.Values{}
			query.Add("redir", r.URL.Path)
			http.Redirect(w, r, "/sign-in?"+query.Encode(), http.StatusFound)
			return
		}
		next(w, r, p)
	}
}
