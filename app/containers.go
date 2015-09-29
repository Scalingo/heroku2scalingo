package app

import (
	"fmt"
	"strings"

	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/github.com/bgentry/heroku-go"
	"github.com/Scalingo/heroku2scalingo/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/heroku2scalingo/config"
)

func ScaleContainers(herokuAppName, scalingoAppName string) error {
	dynos, err := config.HerokuClient.DynoList(herokuAppName, nil)
	if err != nil {
		return errgo.Mask(err)
	}

	var containers []scalingo.Container
	for _, d := range dynos {
		isAlreadyContainer := false
		for _, c := range containers {
			if d.Name == c.Name {
				c.Amount += 1
				isAlreadyContainer = true
				break
			}
		}
		if !isAlreadyContainer {
			containers = append(containers, getContainerFromDyno(d))
		}
	}

	if c := getContainerByName(containers, "web"); c != nil {
		_, err = scalingo.AppsScale(scalingoAppName, &scalingo.AppsScaleParams{containers})
		if err == nil {
			fmt.Printf("%s has been scaled to:\n", scalingoAppName)
			fmt.Printf("-------> %s: %d - %s\n", c.Name, c.Amount, c.Size)
		}

		return nil
	}

	fmt.Printf("[Scale]: Couldn't find any 'web' dyno on your Heroku app, please check your Procfile or scale it manually.\n")

	return nil
}

func getContainerFromDyno(dyno heroku.Dyno) scalingo.Container {
	if dyno.Size == "Free" {
		dyno.Size = "M"
	}

	dotIndex := strings.LastIndex(dyno.Name, ".")
	if dotIndex > -1 {
		dyno.Name = dyno.Name[0:dotIndex]
	}

	return scalingo.Container{
		Name:    dyno.Name,
		Amount:  1,
		Command: dyno.Command,
		Size:    dyno.Size,
	}
}

func getContainerByName(containers []scalingo.Container, name string) *scalingo.Container {
	for _, c := range containers {
		if c.Name == name {
			return &c
		}
	}
	return nil
}
