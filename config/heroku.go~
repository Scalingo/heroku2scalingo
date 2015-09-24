package config

import (
	"log"
	"os/user"
	"path"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/bgentry/heroku-go"
)

var (
	HerokuClient heroku.Client
	machine      *netrc.Machine
	herokuApiUrl = "api.heroku.com"
)

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	machine, err = netrc.FindMachine(path.Join(usr.HomeDir, ".netrc"), herokuApiUrl)
	if err != nil {
		log.Fatal(err)
	}

	HerokuClient = heroku.Client{Username: machine.Login, Password: machine.Password}
}
