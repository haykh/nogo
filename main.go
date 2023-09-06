package main

import (
	"fmt"
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
		Action: func(cCtx *cli.Context) error {
			cli.ShowAppHelp(cCtx)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "configure nogo",
				Action: func(cCtx *cli.Context) error {
					config.CreateOrReadLocalConfig(false)
					return nil
				},
			},
			{
				Name:                   "stack",
				Aliases:                []string{"s"},
				Usage:                  "show the stack",
				UseShortOptionHandling: true,
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "ls", Aliases: []string{"l"}},
					&cli.BoolFlag{Name: "rm", Aliases: []string{"r"}},
					&cli.StringFlag{Name: "add", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "do", Aliases: []string{"d"}},
					&cli.StringFlag{Name: "undo", Aliases: []string{"u"}},
					&cli.StringFlag{Name: "mod", Aliases: []string{"m"}},
				},
				Action: func(cCtx *cli.Context) error {
					show := cCtx.Bool("ls")
					rm := cCtx.Bool("rm")
					add := cCtx.String("add")
					do := cCtx.String("do")
					undo := cCtx.String("undo")
					mod := cCtx.String("mod")
					if !show && !rm && add == "" && do == "" && mod == "" {
						cli.ShowSubcommandHelp(cCtx)
						return nil
					}
					if local_config, err := config.CreateOrReadLocalConfig(true); err != nil {
						return err
					} else {
						if token, err := local_config.GetSecret("api_token"); err != nil {
							return err
						} else {
							if stackID, err := local_config.GetSecret("stack_page_id"); err != nil {
								return err
							} else {
								client := notion.NewClient(token)
								err = nil
								if rm {
									err = notion.RmFromStack(client, stackID)
									fmt.Println()
								}
								if mod != "" && err == nil {
									err = notion.ModifyStack(client, stackID, mod)
									fmt.Println()
								}
								if add != "" && err == nil {
									err = notion.AddToStack(client, stackID, add)
									fmt.Println()
								}
								if do != "" && err == nil {
									err = notion.MarkAs(client, stackID, do, true)
									fmt.Println()
								}
								if undo != "" && err == nil {
									err = notion.MarkAs(client, stackID, undo, false)
									fmt.Println()
								}
								if show && err == nil {
									err = notion.ShowPage(client, stackID)
								}
								return err
							}
						}
					}
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
