package app

import (
	"fmt"

	"github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/heroku-go/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/heroku2scalingo/config"
	"github.com/bgentry/heroku-go"
)

func SetScalingoEnv(app string) error {
	var variables scalingo.Variables

	env, err := getHerokuEnv()
	if err != nil {
		return errgo.Mask(err)
	}

	for k := range env {
		variables = append(variables, &scalingo.Variable{
			Name:  k,
			Value: env[k],
		})
	}

	_, _, err = scalingo.VariableMultipleSet(app, variables)
	if err != nil {
		return errgo.Mask(err)
	}

	for k := range env {
		fmt.Printf("-----> %s has been set to %s\n", k, env[k])
	}

	return nil
}

func getHerokuEnv() (map[string]string, error) {
	env, err := config.HerokuClient.ConfigVarInfo(HerokuApp.Name)
	if err != nil {
		return errgo.Mask(err)
	}

	return env, nil
}
