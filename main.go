package main

import (
	"log"
	notion "nogo/api"
	"nogo/config"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

func main() {
	log.SetPrefix("[ nogo ERROR ]: ")
	log.SetFlags(0)
	app := &cli.App{
		Name:  "nogo",
		Usage: "do awesome stuff with notion from cli",
		Action: func(c *cli.Context) error {
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "configure",
				Aliases: []string{"c"},
				Usage:   "(⚙) configure nogo using prompt",
				Action: func(c *cli.Context) error {
					config.CreateOrReadLocalConfig(false)
					return nil
				},
			},
			{
				Name:    "ls",
				Aliases: []string{"l"},
				Usage:   "(🏠) show the main page",
				Action: func(c *cli.Context) error {
					local_config := config.CreateOrReadLocalConfig(true)
					token, err := local_config.GetSecret("api_token")
					if err != nil {
						return err
					}
					pageID, err := local_config.GetSecret("main_page_id")
					if err != nil {
						return err
					}
					return notion.ShowPage(token, pageID)
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
