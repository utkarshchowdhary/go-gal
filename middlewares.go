package main

import (
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

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
