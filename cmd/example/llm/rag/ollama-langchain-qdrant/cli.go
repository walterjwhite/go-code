package main

import (
	"github.com/urfave/cli/v3"
)

func cliCommand() *cli.Command {
	return &cli.Command{
		Name:  "Simple CLI RAG Tool",
		Usage: "A simple CLI RAG Tool with index and query commands",
		Commands: []*cli.Command{
			{
				Name:  "index",
				Usage: "file to index",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "File to be indexed",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "collection",
						Aliases: []string{"c"},
						Usage:   "Collection Name",
						Value:   "default",
					},
				},
				Action: index,
			},
			{
				Name:  "query",
				Usage: "Prints the query string",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "collection",
						Aliases: []string{"c"},
						Usage:   "Collection name",
						Value:   "default",
					},
					&cli.StringFlag{
						Name:     "query",
						Aliases:  []string{"q"},
						Usage:    "Query",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "docs",
						Aliases: []string{"d"},
						Usage:   "docs",
						Value:   5,
					},
					&cli.FloatFlag{
						Name:    "threshold",
						Aliases: []string{"t"},
						Usage:   "threshold",
						Value:   0.6,
					},
				},

				Action: query,
			},
		},
	}
}
