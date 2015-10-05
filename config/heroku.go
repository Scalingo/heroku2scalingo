package config

import (
	"log"
	"os/user"
	"path"

	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/github.com/bgentry/go-netrc/netrc"
	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/github.com/bgentry/heroku-go"
	"github.com/Scalingo/heroku2scalingo/input"
)

var (
	HerokuClient heroku.Client
	machine      *netrc.Machine
	herokuApiUrl = "api.heroku.com"
)

func InitHerokuAuth() {
	apiKey := ""
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	answer := true
	machine, err = netrc.FindMachine(path.Join(usr.HomeDir, ".netrc"), herokuApiUrl)
	if err == nil {
		answer = input.AskForConfirmation("A Heroku authentication token has been found, do you allow us to use it? [y/n] ")
		apiKey = machine.Password
	}
	if err != nil || !answer {
		apiKey = input.AskForString("We need an api token in order to get the environment of your heroku app (https://dashboard.heroku.com/account)\nApi key: ")
	}
	HerokuClient = heroku.Client{Password: apiKey}
}
