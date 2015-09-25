package git

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/go-scalingo"
	"github.com/bgentry/heroku-go"
)

func CloneHerokuApp(app *heroku.App) error {
	err := os.Chdir(os.TempDir())
	if err != nil {
		return errgo.Mask(err)
	}

	err = os.RemoveAll(path.Join(os.TempDir(), app.Name))
	if err != nil {
		return errgo.Mask(err)
	}

	cmd := exec.Command("git", "clone", app.GitURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func PushScalingoApp(herokuAppName string) error {
	err := os.Chdir(path.Join(os.TempDir(), herokuAppName))
	if err != nil {
		return errgo.Mask(err)
	}

	cmd := exec.Command("git", "push", "scalingo", "master")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func AddRemotes(app *scalingo.App, herokuAppName string) error {
	var remote string
	remote += "[remote \"scalingo\"]\n"
	remote += "\turl = " + app.GitUrl + "\n"
	remote += fmt.Sprintf("\tfetch = +refs/heads/*:refs/remotes/%s/*", app.Name) + "\n"

	filename := path.Join(os.TempDir(), herokuAppName, ".git", "config")

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		return errgo.Mask(err)
	}

	_, err = f.WriteString(remote)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
