package config

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
)

func CloneHerokuApp(appName string) error {
	herokuRepo := fmt.Sprintf("git@heroku.com:%s.git", appName)

	os.Chdir(os.TempDir())
	os.RemoveAll(path.Join(os.TempDir(), appName))
	fmt.Println("Cloning " + herokuRepo + " ...")
	b, err := exec.Command("git", "clone", herokuRepo).Output()
	if err != nil {
		return errgo.Mask(err)
	}

	lines := strings.Split(string(b), "\n")
	for i := range lines {
		fmt.Println(lines[i])
	}

	return nil
}

func AddRemotes(appName string) error {
	var remote string
	remote += "[remote \"scalingo\"]\n"
	remote += "\turl = " + fmt.Sprintf("git@scalingo.com:%s.git", appName) + "\n"
	remote += "\tfetch = " + fmt.Sprintf("+refs/heads/*:refs/remotes/%s/*", appName) + "\n"

	filename := path.Join(os.TempDir(), appName, ".git", "config")

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
