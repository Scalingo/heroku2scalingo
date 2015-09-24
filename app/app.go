package app

import (
	"fmt"

	"github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/go-scalingo/Godeps/_workspace/src/gopkg.in/errgo.v1"
)

func Create(appName string) (*scalingo.App, error) {
	for len(appName) < 6 || len(appName) > 32 {
		appName = getNewAppName(appName)
	}

	app, err := scalingo.AppsCreate(appName)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return app, nil
}

func getNewAppName(appName string) string {
	input := "a"
	inputConfirm := ""

	fmt.Println("Your app '" + appName + "' should contain between 6 and 32 characters")
	for input != inputConfirm {
		fmt.Print("New app name: ")
		fmt.Scanln(&input)
		fmt.Print("Confirmation: ")
		fmt.Scanln(&inputConfirm)
		fmt.Println()
	}

	return inputConfirm
}
