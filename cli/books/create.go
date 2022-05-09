package books

import "github.com/urfave/cli/v2"

var create = &cli.Command{
	Name:        "new",
	Usage:       "create new book",
	Description: "this is longer text for create",
}
