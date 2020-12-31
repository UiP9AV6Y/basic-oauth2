package main

import (
	"fmt"
	"os"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/command"
)

func main() {
	if err := command.BasicOauth2(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
