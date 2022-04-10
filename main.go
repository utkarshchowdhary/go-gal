package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func init() {
	userStore, err := NewFileUserStore("./data/users.json")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}
	globalUserStore = userStore

	sessionStore, err := NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}
	globalSessionStore = sessionStore

	db, err := NewPostgresDb(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	globalPostgresDb = db

	imageStore, err := NewDbImageStore("./data/images")
	if err != nil {
		panic(err)
	}
	globalImageStore = imageStore
}

func main() {
	router := httprouter.New()

	router.GET("/", Recoverer(HandleHome))
	router.GET("/sign-up", Recoverer(HandleUserNew))
	router.POST("/sign-up", Recoverer(HandleUserCreate))
	router.GET("/sign-in", Recoverer(HandleSessionNew))
	router.POST("/sign-in", Recoverer(HandleSessionCreate))
	router.GET("/image/:imageId", Recoverer(HandleImageShow))
	router.GET("/user/:userId", Recoverer(HandleUserShow))
	router.GET("/sign-out", Recoverer(Authenticator(HandleSessionDestroy)))
	router.GET("/account", Recoverer(Authenticator(HandleUserEdit)))
	router.POST("/account", Recoverer(Authenticator(HandleUserUpdate)))
	router.GET("/images/new", Recoverer(Authenticator(HandleImageNew)))
	router.POST("/images/new", Recoverer(Authenticator(HandleImageCreate)))

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	router.ServeFiles("/img/*filepath", http.Dir("data/images/"))

	log.Println("Server starting at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
