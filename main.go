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
		panic(fmt.Errorf("Error establishing connection to database: %s", err))
	}
	globalPostgresDb = db

	err = DbInitSchema()
	if err != nil {
		panic(fmt.Errorf("Error initializing database: %s", err))
	}

	imageStore, err := NewDbImageStore("./data/images")
	if err != nil {
		panic(fmt.Errorf("Error creating image store: %s", err))
	}
	globalImageStore = imageStore
}

func main() {
	router := httprouter.New()
	router.PanicHandler = HandlePanic

	router.GET("/", HandleHome)
	router.GET("/sign-up", HandleUserNew)
	router.POST("/sign-up", HandleUserCreate)
	router.GET("/sign-in", HandleSessionNew)
	router.POST("/sign-in", HandleSessionCreate)
	router.GET("/image/:imageId", HandleImageShow)
	router.GET("/user/:userId", HandleUserShow)
	router.GET("/sign-out", Authenticator(HandleSessionDestroy))
	router.GET("/account", Authenticator(HandleUserEdit))
	router.POST("/account", Authenticator(HandleUserUpdate))
	router.GET("/images/new", Authenticator(HandleImageNew))
	router.POST("/images/new", Authenticator(HandleImageCreate))

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	router.ServeFiles("/img/*filepath", http.Dir("data/images/"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	log.Println("Server starting at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
