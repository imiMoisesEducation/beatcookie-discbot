package main

import (
	"os"

	"github.com/imiMoisesEducation/beat-cookie-discbot/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
