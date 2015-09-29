package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AskForString(question string) string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	fmt.Print(question)
	scanner.Scan()
	return scanner.Text()
}

func AskForConfirmation(confirmation string) bool {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	fmt.Print(confirmation)
	scanner.Scan()
	response := strings.ToLower(scanner.Text())
	if response == "yes" || response == "y" {
		return true
	} else if response == "no" || response == "n" {
		return false
	} else {
		return AskForConfirmation(confirmation)
	}
}
