package auth

import (
	"database/sql"
	"time"

	"github.com/Pashakrut94/SwiftChat/users"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func validateSignUpRequest(username, password, phone string) error {
	if len(username) < 2 {
		return errors.Wrap(ErrValidationError, "username is too short")
	}
	if len(password) < 8 {
		return errors.Wrap(ErrValidationError, "password can't be shorter then 8 symbols")
	}
	if len(phone) < 12 {
		return errors.Wrap(ErrValidationError, "incorrect phone number")
	}
	return nil
}

func HandleSignUp(userRepo *users.UserRepo, username, password, phone string) (*users.User, error) {
	if err := validateSignUpRequest(username, password, phone); err != nil {
		return nil, err
	}
	_, err := userRepo.GetByPhone(phone)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "error getting user by phone")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate password hash")
	}
	user := users.User{Name: username, Phone: phone, Password: string(hashedPassword)}
	if err := userRepo.Create(&user); err != nil {
		return nil, errors.Wrap(err, "unable to create user")
	}
	return &user, nil
}

func validateSignIn(phone, password string) error {
	if len(phone) < 12 {
		return errors.Wrap(ErrValidationError, "incorrect input of phone number")
	}
	if len(password) < 8 {
		return errors.Wrap(ErrValidationError, "password can't be shorter then 8 symbols")
	}
	return nil
}

func HandleSignIn(userRepo *users.UserRepo, phone, password string) (*users.User, error) {
	if err := validateSignIn(phone, password); err != nil {
		return nil, err
	}
	user, err := userRepo.GetByPhone(phone)
	if err != nil {
		return nil, errors.Wrap(err, "error getting user by phone from DB")
	}
	if err == sql.ErrNoRows {
		return nil, errors.Wrap(err, "error getting user from DB")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.Wrap(err, "mismatch converted password and hash")
	}
	return user, nil
}

func HandleAuthentication(sessRepo SessionRepo, cookieValue string) (Session, error) {
	session, err := sessRepo.Get(cookieValue)
	if err == sql.ErrNoRows {
		return Session{}, errors.Wrap(err, "error session not found")
	}
	if err != nil {
		return Session{}, errors.Wrap(err, "error getting session from DB")
	}
	if time.Now().After(session.ExpiresAt) {
		if err := sessRepo.Delete(cookieValue); err != nil {
			return Session{}, errors.Wrap(err, "error updating deleted_at field")
		}
		return Session{}, errors.Wrap(err, "session expired")
	}
	return *session, nil
}

func SetNewSession(sessRepo SessionRepo, userID int) (sid string, expires time.Time, err error) {
	sessID, err := uuid.NewV4()
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "error getting hash of the session")
	}
	sid = sessID.String()
	expires = time.Now().Add(time.Minute * 72000)
	if err = sessRepo.Create(sid, userID, expires); err != nil {
		return sid, expires, errors.Wrap(err, "error creating new session")
	}
	return
}
