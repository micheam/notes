package books

import "github.com/urfave/cli/v2"

var Cmd = &cli.Command{
	Name: "books",
	Subcommands: []*cli.Command{
		create,
	},
}
