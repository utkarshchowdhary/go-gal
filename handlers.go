package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleHome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	images, err := globalImageStore.FindAll(0)
	if err != nil {
		panic(err)
	}
	RenderTemplate(w, r, "index/home", map[string]interface{}{
		"Images": images,
	})
}

func HandleUserNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	RenderTemplate(w, r, "users/new", nil)
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := NewUser(
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("password"),
	)

	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "users/new", map[string]interface{}{
				"Error": err,
				"User":  user,
			})
			return
		}
		panic(err)
	}

	err = globalUserStore.Save(user)
	if err != nil {
		panic(err)
	}

	session := NewSession(w)
	session.UserId = user.Id
	err = globalSessionStore.Save(session)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/?flash=Your+account+has+been+successfully+created", http.StatusFound)
}

func HandleSessionNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redir := r.FormValue("redir")
	RenderTemplate(w, r, "sessions/new", map[string]interface{}{
		"Redir": redir,
	})
}

func HandleSessionCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	redir := r.FormValue("redir")

	user, err := FindUser(username, password)
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "sessions/new", map[string]interface{}{
				"Error": err,
				"User":  user,
				"Redir": redir,
			})
			return
		}
		panic(err)
	}

	session := FindOrCreateSession(w, r)
	session.UserId = user.Id
	err = globalSessionStore.Save(session)
	if err != nil {
		panic(err)
	}

	if redir == "" {
		redir = "/"
	}

	http.Redirect(w, r, redir+"?flash=Signed+in", http.StatusFound)
}

func HandleSessionDestroy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := RequestSession(r)
	if session != nil {
		err := globalSessionStore.Delete(session)
		if err != nil {
			panic(err)
		}
	}

	http.Redirect(w, r, "/?flash=Your+have+been+successfully+signed+out", http.StatusFound)
}

func HandleUserEdit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := RequestUser(r)
	RenderTemplate(w, r, "users/edit", map[string]interface{}{
		"User": user,
	})
}

func HandleUserUpdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	currentUser := RequestUser(r)
	username := r.FormValue("username")
	email := r.FormValue("email")
	currentPassword := r.FormValue("currentPassword")
	newPassword := r.FormValue("newPassword")

	user, err := UpdateUser(currentUser, username, email, currentPassword, newPassword)
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "users/edit", map[string]interface{}{
				"Error": err,
				"User":  user,
			})
			return
		}
		panic(err)
	}

	err = globalUserStore.Save(*currentUser)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/account?flash=Your+changes+have+been+successfully+saved", http.StatusFound)
}

func HandleImageNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	RenderTemplate(w, r, "images/new", nil)
}

func HandleImageCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.FormValue("url") != "" {
		HandleImageCreateFromURL(w, r)
		return
	}

	HandleImageCreateFromFile(w, r)
}

func HandleImageCreateFromURL(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	image := NewImage(user)
	image.Description = r.FormValue("description")

	err := image.CreateFromUrl(r.FormValue("url"))
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "images/new", map[string]interface{}{
				"Error":    err,
				"ImageUrl": r.FormValue("url"),
				"Image":    image,
			})
			return
		}
		panic(err)
	}

	http.Redirect(w, r, "/?flash=Image+uploaded+successfully", http.StatusFound)
}

func HandleImageCreateFromFile(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	image := NewImage(user)
	image.Description = r.FormValue("description")

	file, headers, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			RenderTemplate(w, r, "images/new", map[string]interface{}{
				"Error": errNoImage,
				"Image": image,
			})
			return
		}
		panic(err)
	}
	defer file.Close()

	err = image.CreateFromFile(file, headers)
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "images/new", map[string]interface{}{
				"Error": err,
				"Image": image,
			})
			return
		}
		panic(err)
	}

	http.Redirect(w, r, "/?flash=Image+uploaded+successfully", http.StatusFound)
}

func HandleImageShow(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	image, err := globalImageStore.Find(p.ByName("imageId"))
	if err != nil {
		if err == sql.ErrNoRows {
			RenderTemplate(w, r, "images/show", map[string]interface{}{
				"Error": errors.New("The request resource does not exist"),
			})
			return
		}
		panic(err)
	}

	user, err := globalUserStore.Find(image.UserId)
	if err != nil {
		panic(err)
	}
	if user == nil {
		panic(fmt.Errorf("Could not find user %s", image.UserId))
	}

	RenderTemplate(w, r, "images/show", map[string]interface{}{
		"Image": image,
		"User":  user,
	})
}

func HandleUserShow(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := globalUserStore.Find(p.ByName("userId"))
	if err != nil {
		panic(err)
	}

	if user == nil {
		RenderTemplate(w, r, "users/show", map[string]interface{}{
			"Error": errors.New("The request resource does not exist"),
		})
		return
	}

	images, err := globalImageStore.FindAllByUser(user, 0)
	if err != nil {
		panic(err)
	}

	RenderTemplate(w, r, "users/show", map[string]interface{}{
		"Images": images,
		"User":   user,
	})
}
