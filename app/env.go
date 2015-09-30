package app

import (
	"strings"

	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/heroku2scalingo/config"
	"github.com/Scalingo/heroku2scalingo/io"
)

func SetScalingoEnv(herokuAppName, scalingoAppName string) error {
	var variables scalingo.Variables

	env, err := getHerokuEnv(herokuAppName)
	if err != nil {
		return errgo.Mask(err)
	}

	for k := range env {
		if strings.TrimSpace(env[k]) == "" {

			io.Warnf("%s has an empty value and couldn't be set\n", k)
		} else {
			variables = append(variables, &scalingo.Variable{
				Name:  k,
				Value: env[k],
			})
		}
	}

	_, _, err = scalingo.VariableMultipleSet(scalingoAppName, variables)
	if err != nil {
		return errgo.Mask(err)
	}

	for _, v := range variables {
		io.Printf("%s has been set to %s\n", v.Name, v.Value)
	}

	return nil
}

func getHerokuEnv(herokuAppName string) (map[string]string, error) {
	env, err := config.HerokuClient.ConfigVarInfo(herokuAppName)
	if err != nil {
		return env, errgo.Mask(err)
	}

	return env, nil
}
