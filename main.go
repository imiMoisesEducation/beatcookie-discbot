package main

import (
	"os"

	_ "github.com/imiMoisesEducation/beat-cookie-discbot/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
