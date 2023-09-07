package main

import (
	"log"
	"os"
	"sort"

	"github.com/haykh/nogo/config"

	notion "github.com/haykh/nogo/api"

	"github.com/urfave/cli/v2"
)

func main() {
	log.SetPrefix("[ nogo ERROR ]: ")
	log.SetFlags(0)
	app := &cli.App{
		Name:  "nogo",
		Usage: "do awesome stuff with notion from cli",
		Action: func(cCtx *cli.Context) error {
			return cli.ShowAppHelp(cCtx)
		},
		Commands: []*cli.Command{
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "configure nogo",
				Action: func(cCtx *cli.Context) error {
					_, err := config.CreateOrReadLocalConfig(false)
					return err
				},
			},
			{
				Name:                   "stack",
				Aliases:                []string{"s"},
				Usage:                  "interact with the stack (todo list)",
				UseShortOptionHandling: true,
				Action: func(cCtx *cli.Context) error {
					if client, sID, err := notion.InitAPI(); err != nil {
						return nil
					} else {
						return notion.ShowPage(client, sID)
					}
				},
				Subcommands: []*cli.Command{
					{
						Name:    "add",
						Aliases: []string{"a"},
						Usage:   "add a new entry to the stack",
						Action: func(cCtx *cli.Context) error {
							if client, sID, err := notion.InitAPI(); err != nil {
								return nil
							} else {
								return notion.AddToStack(client, sID)
							}
						},
					},
					{
						Name:    "mod",
						Aliases: []string{"m"},
						Usage:   "modify a stack entry",
						Action: func(cCtx *cli.Context) error {
							if client, sID, err := notion.InitAPI(); err != nil {
								return nil
							} else {
								return notion.ModifyStack(client, sID)
							}
						},
					},
					{
						Name:    "toggle",
						Aliases: []string{"t"},
						Usage:   "toggle stack entries",
						Action: func(cCtx *cli.Context) error {
							if client, sID, err := notion.InitAPI(); err != nil {
								return nil
							} else {
								return notion.ToggleStack(client, sID)
							}
						},
					},
					{
						Name:    "rm",
						Aliases: []string{"r"},
						Usage:   "remove stack entries",
						Action: func(cCtx *cli.Context) error {
							if client, sID, err := notion.InitAPI(); err != nil {
								return nil
							} else {
								return notion.RmFromStack(client, sID)
							}
						},
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
