package app

import (
	"fmt"

	"github.com/Scalingo/go-scalingo"
)

func Create(appName string) error {
	fmt.Println(appName + "created")
	return nil
}
