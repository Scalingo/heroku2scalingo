package main

import (
	"fmt"
	"os"

	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/github.com/bgentry/heroku-go"
	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/heroku2scalingo/app"
	"github.com/Scalingo/heroku2scalingo/config"
	"github.com/Scalingo/heroku2scalingo/git"
	"github.com/Scalingo/heroku2scalingo/io"
	"github.com/Scalingo/heroku2scalingo/signals"
)

var (
	ScalingoApp *scalingo.App
	HerokuApp   *heroku.App
)

func PushRepository() error {
	fmt.Println()
	io.Info("Pushing to", ScalingoApp.GitUrl+"...\n")

	err := git.PushScalingoApp(HerokuApp.Name)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func CloneRepository() error {
	io.Info("Cloning Heroku GIT repository\n")
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

	io.Infof("Creating scalingo app %s...\n", HerokuApp.Name)

	ScalingoApp, err = app.Create(HerokuApp.Name)
	if err != nil {
		return errgo.Mask(err)
	}
	io.Print("Scalingo App '" + ScalingoApp.Name + "' created.\n")

	io.Info("Importing Heroku environment to Scalingo...")
	err = app.SetScalingoEnv(HerokuApp.Name, ScalingoApp.Name)
	if err != nil {
		return errgo.Mask(err)
	}
	fmt.Println()
	io.Print("Importation successful\n")

	return nil
}

func main() {
	if len(os.Args) <= 1 {
		io.Error("<Usage>:", os.Args[0], "<app-name>")
	}

	go signals.Handle()

	config.InitHerokuAuth()
	io.Infof("Heroku authentication... ")
	var err error
	HerokuApp, err = config.HerokuClient.AppInfo(os.Args[1])
	if err != nil {
		fmt.Println("ERR")
		io.Error(err)
	}
	fmt.Println("OK")

	io.Infof("Scalingo authentication... ")
	u, err := config.LoadAuthOrLogin()
	if err != nil {
		fmt.Println("ERR")
		io.Error(err)
	}
	fmt.Println("OK")

	io.Printf("You are now logged as %s/%s\n\n", u.Username, u.Email)
	err = CreateScalingoApp()
	if err != nil {
		io.Error(err)
	}

	err = CloneRepository()
	if err != nil {
		io.Error(err)
	}

	err = PushRepository()
	if err != nil {
		io.Error(err)
	}
}
