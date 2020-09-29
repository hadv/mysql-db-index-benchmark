package utils

import (
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

// NewApp creates an app with sane defaults.
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = ""
	app.Email = ""
	app.Version = ""
	app.Usage = ""
	return app
}
