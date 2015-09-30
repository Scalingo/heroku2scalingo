package app

import (
	"fmt"
	"strings"

	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/heroku2scalingo/io"
)

func Create(appName string) (*scalingo.App, error) {
	for len(appName) < 6 || len(appName) > 32 {
		fmt.Println("Your app '" + appName + "' should contain between 6 and 32 characters")
		appName = getNewAppName()
	}

	app, err := scalingo.AppsCreate(appName)
	if err != nil {
		if strings.Contains(err.Error(), "is already taken") {
			io.Warnf("The name '%s' is already taken.\n\n", appName)
			return Create(getNewAppName())
		}
		return nil, errgo.Mask(err)
	}

	return app, nil
}

func getNewAppName() string {
	input := "a"
	inputConfirm := ""

	for input != inputConfirm {
		fmt.Print("New app name: ")
		fmt.Scanln(&input)
		fmt.Print("Confirmation: ")
		fmt.Scanln(&inputConfirm)
		fmt.Println()
	}

	return inputConfirm
}
