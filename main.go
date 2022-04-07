package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.GET("/", Recoverer(HandleHome))
	router.GET("/sign-up", Recoverer(HandleUserNew))
	router.POST("/sign-up", Recoverer(HandleUserCreate))
	router.GET("/sign-in", Recoverer(HandleSessionNew))
	router.POST("/sign-in", Recoverer(HandleSessionCreate))
	router.GET("/sign-out", Recoverer(Authenticator(HandleSessionDestroy)))
	router.GET("/account", Recoverer(Authenticator(HandleUserEdit)))
	router.POST("/account", Recoverer(Authenticator(HandleUserUpdate)))

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	log.Println("Server starting at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
