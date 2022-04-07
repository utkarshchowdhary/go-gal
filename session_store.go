package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type SessionStore interface {
	Find(string) (*Session, error)
	Save(*Session) error
	Delete(*Session) error
}

type FileSessionStore struct {
	filename string
	Sessions map[string]Session
}

var globalSessionStore SessionStore

func init() {
	store, err := NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}
	globalSessionStore = store
}

func NewFileSessionStore(filename string) (*FileSessionStore, error) {
	store := &FileSessionStore{
		filename,
		map[string]Session{},
	}
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}
	err = json.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}
	return store, err
}

func (store FileSessionStore) Find(id string) (*Session, error) {
	session, ok := store.Sessions[id]
	if ok {
		return &session, nil
	}
	return nil, nil
}

func (store FileSessionStore) Save(session *Session) error {
	store.Sessions[session.Id] = *session
	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(store.filename, contents, 0660)
}

func (store FileSessionStore) Delete(session *Session) error {
	delete(store.Sessions, session.Id)
	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(store.filename, contents, 0660)
}
