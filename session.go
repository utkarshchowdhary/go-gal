package main

import (
	"net/http"
	"time"
)

const (
	sessionCookieName = "GoGalSession"
	sessionIdLength   = 20
	sessionLength     = 24 * time.Hour
)

type Session struct {
	Id     string
	UserId string
	Expiry time.Time
}

func (session Session) Expired() bool {
	return session.Expiry.Before(time.Now())
}

func NewSession(w http.ResponseWriter) *Session {
	expiry := time.Now().Add(sessionLength)
	session := &Session{
		Id:     GenerateId("sess", sessionIdLength),
		Expiry: expiry,
	}
	cookie := http.Cookie{
		Name:    sessionCookieName,
		Value:   session.Id,
		Expires: expiry,
	}
	http.SetCookie(w, &cookie)
	return session
}

func RequestSession(r *http.Request) *Session {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil
	}
	session, err := globalSessionStore.Find(cookie.Value)
	if err != nil {
		panic(err)
	}
	if session == nil {
		return nil
	}
	if session.Expired() {
		globalSessionStore.Delete(session)
		return nil
	}
	return session
}

func RequestUser(r *http.Request) *User {
	session := RequestSession(r)
	if session == nil {
		return nil
	}
	user, err := globalUserStore.Find(session.UserId)
	if err != nil {
		panic(err)
	}
	return user
}

func FindOrCreateSession(w http.ResponseWriter, r *http.Request) *Session {
	session := RequestSession(r)
	if session == nil {
		session = NewSession(w)
	}

	return session
}
