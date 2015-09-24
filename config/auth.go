package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/go-scalingo/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/go-scalingo/users"
	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/code.google.com/p/gopass"
)

type HtoSAuthenticator struct{}

type AuthConfigData struct {
	LastUpdate        time.Time              `json:"last_update"`
	AuthConfigPerHost map[string]*users.User `json:"auth_config_data"`
}

var Authenticator = &HtoSAuthenticator{}

func Auth() (*users.User, error) {
	var user *users.User
	var err error
	for i := 0; i < 3; i++ {
		user, err = tryAuth()
		if err == nil {
			break
		} else if errgo.Cause(err) == io.EOF {
			return nil, errors.New("canceled by user")
		} else {
			fmt.Printf("Fail to login (%d/3): %v\n", i+1, err)
		}
	}
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	err = Authenticator.StoreAuth(user)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return user, nil
}

func (a *HtoSAuthenticator) StoreAuth(user *users.User) error {
	return nil
}

func (a *HtoSAuthenticator) LoadAuth() (*users.User, error) {
	file, err := os.OpenFile(C.AuthFile, os.O_RDONLY, 0644)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer file.Close()

	var authConfig AuthConfigData
	dec := json.NewDecoder(file)
	if err := dec.Decode(&authConfig); err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	if user, ok := authConfig.AuthConfigPerHost[C.apiHost]; !ok {
		return Auth()
	} else {
		if user == nil {
			return Auth()
		}
		return user, nil
	}
}

func (a *HtoSAuthenticator) RemoveAuth() error {
	return nil
}

func tryAuth() (*users.User, error) {
	var login string
	for login == "" {
		fmt.Print("Username or email: ")
		_, err := fmt.Scanln(&login)
		if err != nil {
			if strings.Contains(err.Error(), "unexpected newline") {
				continue
			}
			return nil, errgo.Mask(err, errgo.Any)
		}
	}

	password, err := gopass.GetPass("Password: ")
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	user, err := scalingo.AuthUser(login, password)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return user, nil
}

func existingAuth() (*AuthConfigData, error) {
	authConfig := &AuthConfigData{}
	content, err := ioutil.ReadFile(C.AuthFile)
	if err == nil {
		// We don't care of the error
		json.Unmarshal(content, &authConfig)
	} else if os.IsNotExist(err) {
		authConfig.AuthConfigPerHost = make(map[string]*users.User)
	} else {
		return nil, errgo.Mask(err)
	}
	return authConfig, nil
}
