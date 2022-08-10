package main

import (
	"log"
	notion "nogo/api"
	"nogo/config"
	"nogo/utils"
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
				Usage:   "(⚙) configure nogo using prompts",
				Action: func(c *cli.Context) error {
					config.CreateOrReadLocalConfig(false)
					return nil
				},
			},
			{
				Name:    "ls",
				Aliases: []string{"l"},
				Usage:   "(🏠) show the main page",
				Subcommands: []*cli.Command{
					{
						Name:    "main",
						Aliases: []string{"m"},
						Usage:   "show the main page",
						Action: func(cCtx *cli.Context) error {
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
					{
						Name:    "todo",
						Aliases: []string{"t"},
						Usage:   "show the todo list",
						Action: func(cCtx *cli.Context) error {
							local_config := config.CreateOrReadLocalConfig(true)
							token, err := local_config.GetSecret("api_token")
							if err != nil {
								return err
							}
							todo_id, err := local_config.GetSecret("todo_page_id")
							if err != nil {
								return err
							}
							return notion.ShowPage(token, todo_id)
						},
					},
				},
			},
			{
				Name:    "todo",
				Aliases: []string{"t"},
				Usage:   "(✔️) interact with the todo list",
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
					todo_id, err := local_config.GetSecret("todo_page_id")
					if err != nil {
						utils.Message("No to-do page found, creating one.", utils.Warning, true)
						newtodo_id, err := notion.CreatePage(token, pageID, "to-do", "🗒️")
						if err != nil {
							return err
						}
						err = local_config.SetSecret("todo_page_id", newtodo_id)
						if err != nil {
							return err
						}
						todo_id = newtodo_id
					}
					return notion.ShowPage(token, todo_id)
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
