package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Zyko0/scalingo-from-heroku/app"
	"github.com/Zyko0/scalingo-from-heroku/config"
	"gopkg.in/errgo.v1"
)

func GetHerokuApp(appName string) error {
	env, err := config.HerokuClient.ConfigVarInfo(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

}

func CloneRepository(appName string) error {
	err := config.CloneHerokuApp(appName)
	if err != nil {
		return errgo.Mask(err)
	}

	err = config.AddRemotes(appName)
	if err != nil {
		return errgo.Mask(err)
	}

	err = config.AddRemotes(appName)
	if err != nil {
		return errgo.Mask(err)
	}

	err = app.Create(appName)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func main() {
	if len(os.Args) <= 1 {
		log.Println("<Usage>: " + os.Args[0] + " <appName>")
		return
	}

	app.Login()

	err := CloneRepository(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}
}
