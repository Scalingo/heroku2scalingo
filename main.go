package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/heroku2scalingo/app"
	"github.com/Scalingo/heroku2scalingo/config"
	"github.com/Scalingo/heroku2scalingo/git"
	"github.com/bgentry/heroku-go"
	"gopkg.in/errgo.v1"
)

var (
	ScalingoApp *scalingo.App
	HerokuApp   *heroku.App
)

func PushRepository() error {
	fmt.Println("Pushing to   " + ScalingoApp.GitUrl + " ...")

	err := git.PushScalingoApp(HerokuApp.Name)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func CloneRepository() error {
	fmt.Println("Cloning from " + HerokuApp.GitURL + " ...")

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

	var env map[string]string
	env, err = config.HerokuClient.ConfigVarInfo(HerokuApp.Name)
	if err != nil {
		return errgo.Mask(err)
	}

	for k := range env {
		fmt.Printf("Add %s=%s\n", k, env[k])
		scalingo.VariableSet(ScalingoApp.Name, k, env[k])
	}

	return nil
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("<Usage>: " + os.Args[0] + " <appName>")
		return
	}

	fmt.Println("Scalingo authentication ...")
	u, err := config.Authenticator.LoadAuth()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("\nThe migration will continue with the user %s / %s\n\n", u.Username, u.Email)

	fmt.Println("Heroku authentication ...")
	HerokuApp, err = config.HerokuClient.AppInfo(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	err = CreateScalingoApp()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = CloneRepository()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = PushRepository()
	if err != nil {
		log.Fatal(err.Error())
	}
}