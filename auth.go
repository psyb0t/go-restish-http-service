package simplehttp

import (
	"fmt"
	"math/rand"
)

type Auth struct {
	Service *HttpService
}

func NewAuth(service *HttpService) *Auth {
	return &Auth{
		Service: service,
	}
}

func (a *Auth) makeUserKeyName(user string) string {
	return "auth_user_" + user
}

func (a *Auth) makeAuthKeyKeyName(key string) string {
	return "auth_key_" + key
}

func (a *Auth) makeAuthKey(user string) string {
	auth_key := make([]byte, 55)
	rand.Read(auth_key)

	return fmt.Sprintf("%x", auth_key)
}

func (a *Auth) CreateAuthentication(user string) (string, error) {
	auth_key := a.makeAuthKey(user)

	err := a.Service.DB.Set(a.makeUserKeyName(user), []byte(auth_key))
	if err != nil {
		return "", err
	}

	err = a.Service.DB.Set(a.makeAuthKeyKeyName(auth_key), []byte(user))
	if err != nil {
		return "", err
	}

	return auth_key, nil
}

func (a *Auth) DeleteAuthentication(user string) error {
	auth_key, err := a.GetUserAuthKey(user)
	if err != nil {
		return err
	}

	err = a.Service.DB.Delete(a.makeUserKeyName(user))
	if err != nil {
		return err
	}

	err = a.Service.DB.Delete(a.makeAuthKeyKeyName(auth_key))
	if err != nil {
		return err
	}

	return nil
}

func (a *Auth) ResetAuthentication(user string) (string, error) {
	err := a.DeleteAuthentication(user)
	if err != nil {
		return "", err
	}

	return a.CreateAuthentication(user)
}

func (a *Auth) GetUserAuthKey(user string) (string, error) {
	val, err := a.Service.DB.Get(a.makeUserKeyName(user))
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (a *Auth) GetAuthKeyUser(key string) (string, error) {
	val, err := a.Service.DB.Get(a.makeAuthKeyKeyName(key))
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (a *Auth) IsValidKey(key string) bool {
	if _, err := a.GetAuthKeyUser(key); err != nil {
		return false
	}

	return true
}

func (a *Auth) IsValidUser(user string) bool {
	if _, err := a.GetUserAuthKey(user); err != nil {
		return false
	}

	return true
}
