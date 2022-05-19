package cli

import (
	"github.com/micheam/notes/cli/books"
	"github.com/urfave/cli/v2"
)

const Version = "devel"

func New() *cli.App {
	return &cli.App{
		Name:    "notes",
		Version: Version,
		Commands: []*cli.Command{
			books.Cmd,
		},
	}
}
