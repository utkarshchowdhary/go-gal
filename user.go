package main

import "golang.org/x/crypto/bcrypt"

const (
	hashCost       = 10
	passwordLength = 8
	userIdLength   = 16
)

type User struct {
	Id             string
	Username       string
	Email          string
	HashedPassword string
}

func NewUser(username, email, password string) (User, error) {
	user := User{
		Email:    email,
		Username: username,
	}

	if username == "" {
		return user, errNoUsername
	}

	if email == "" {
		return user, errNoEmail
	}

	if password == "" {
		return user, errNoPassword
	}

	if len(password) < passwordLength {
		return user, errPasswordTooShort
	}

	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, errUsernameExists
	}

	existingUser, err = globalUserStore.FindByEmail(email)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, errEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)

	user.HashedPassword = string(hashedPassword)
	user.Id = GenerateID("usr", userIdLength)

	return user, err
}

func FindUser(username, password string) (*User, error) {
	user := &User{
		Username: username,
	}

	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return user, err
	}
	if existingUser == nil {
		return user, errCredentialsIncorrect
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(existingUser.HashedPassword),
		[]byte(password),
	) != nil {
		return user, errCredentialsIncorrect
	}

	return existingUser, nil
}

func UpdateUser(currentUser *User, username, email, currentPassword, newPassword string) (User, error) {
	user := *currentUser
	user.Username = username
	user.Email = email

	if username == "" {
		return user, errNoUsername
	}

	if email == "" {
		return user, errNoEmail
	}

	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return user, err
	}
	if existingUser != nil && existingUser.Id != currentUser.Id {
		return user, errUsernameExists
	}

	currentUser.Username = username

	existingUser, err = globalUserStore.FindByEmail(email)
	if err != nil {
		return user, err
	}
	if existingUser != nil && existingUser.Id != currentUser.Id {
		return user, errEmailExists
	}

	currentUser.Email = email

	if currentPassword == "" {
		return user, nil
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(currentUser.HashedPassword),
		[]byte(currentPassword),
	) != nil {
		return user, errPasswordIncorrect
	}

	if newPassword == "" {
		return user, errNoPassword
	}

	if len(newPassword) < passwordLength {
		return user, errPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), hashCost)
	currentUser.HashedPassword = string(hashedPassword)
	return user, err
}
