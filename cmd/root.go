package cmd

import (
	"github.com/imiMoisesEducation/beat-cookie-discbot/beater"

	cmd "github.com/elastic/beats/libbeat/cmd"
)

// Name of this beat
var Name = "beat-cookie-discbot"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmd(Name, "", beater.New)
