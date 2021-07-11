package service

import (
	"errors"
	"io"
	"time"

	"github.com/eaciit/toolkit"
	"github.com/kharism/microservice_simple/model"
	"github.com/kharism/microservice_simple/repository"
	"github.com/kharism/microservice_simple/util"
	//"kano/simwas/pkg/module/toolkit"
)

const (
	// Salt salt for encrypting password
	Salt = "asdadszxccx"
)

// IAuth auth service interface
type IAuth interface {
	VerifyPassword(username string, password string) (model.User, error)
	RegisterUser(user model.User) error
}

type auth struct {
	user func() repository.IUser
}

// NewAuth create new service instance
func NewAuth() IAuth {
	return auth{
		user: repository.NewUser,
	}
}
func (s auth) RegisterUser(user model.User) error {
	_, err := s.user().FindByUsername(user.Username)
	if err == nil {
		return errors.New("User Sudah ada")
	}
	hashed, err := util.Hash(user.Password)
	user.PasswordHash = hashed
	if err != nil {
		return err
	}
	_, err = s.user().Save(user)
	return err
}
func (s auth) VerifyPassword(username string, password string) (model.User, error) {
	user, err := s.user().FindByUsername(username)
	if err != nil {
		if err == io.EOF {
			return user, errors.New("username atau password salah")
		}
		return user, err
	}

	if !user.VerifyPassword(password) {
		return user, toolkit.Error(ErrVerifyPasswordInvalid)
	}
	now := time.Now()
	user.LastLogin = &now
	user, err = s.user().Save(user)
	if err != nil {
		return user, err
	}

	return user, nil
}
