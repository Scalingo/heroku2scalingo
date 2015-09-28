package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/heroku2scalingo/app"
	"github.com/Scalingo/heroku2scalingo/config"
	"github.com/Scalingo/heroku2scalingo/git"
	"github.com/Scalingo/heroku2scalingo/signals"
	"github.com/bgentry/heroku-go"
	"gopkg.in/errgo.v1"
)

var (
	ScalingoApp *scalingo.App
	HerokuApp   *heroku.App
)

func PushRepository() error {
	fmt.Println("Pushing to " + ScalingoApp.GitUrl + " ...")

	err := git.PushScalingoApp(HerokuApp.Name)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func CloneRepository() error {
	err := git.CloneHerokuApp(HerokuApp)
	if err != nil {
		return errgo.Mask(err)
	}

	err = git.AddRemotes(ScalingoApp, HerokuApp.Name)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func CreateScalingoApp() error {
	var err error

	fmt.Printf("Creating scalingo app %s ...\n", HerokuApp.Name)

	ScalingoApp, err = app.Create(HerokuApp.Name)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Println("Scalingo App '" + ScalingoApp.Name + "' created.")
	fmt.Println()
	fmt.Println("Importing Heroku environment to Scalingo ...")

	err = app.SetScalingoEnv(ScalingoApp.Name)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("<Usage>: " + os.Args[0] + " <appName>")
		return
	}

	go signals.Handle()

	fmt.Println("Heroku authentication ...")
	var err error
	HerokuApp, err = config.HerokuClient.AppInfo(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println()

	fmt.Println("Scalingo authentication ...")
	u, err := config.Authenticator.LoadAuth()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println()
	fmt.Printf("The migration will continue with the user %s / %s\n\n", u.Username, u.Email)

	err = CreateScalingoApp()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println()

	err = CloneRepository()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println()

	err = PushRepository()
	if err != nil {
		log.Fatal(err.Error())
	}
}
