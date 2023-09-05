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
					&cli.StringFlag{Name: "rm", Aliases: []string{"r"}},
					&cli.StringFlag{Name: "add", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "done", Aliases: []string{"d"}},
					&cli.StringFlag{Name: "undone", Aliases: []string{"u"}},
				},
				Action: func(cCtx *cli.Context) error {
					show := cCtx.Bool("ls")
					rm := cCtx.String("rm")
					add := cCtx.String("add")
					done := cCtx.String("done")
					undone := cCtx.String("undone")
					if !show && rm == "" && add == "" && done == "" {
						cli.ShowSubcommandHelp(cCtx)
						return nil
					}
					local_config := config.CreateOrReadLocalConfig(true)
					if token, err := local_config.GetSecret("api_token"); err != nil {
						return err
					} else {
						if stackID, err := local_config.GetSecret("stack_page_id"); err != nil {
							return err
						} else {
							err = nil
							if rm != "" {
								err = notion.RmFromStack(token, stackID, rm)
								fmt.Println()
							}
							if add != "" && err == nil {
								err = notion.AddToStack(token, stackID, add)
								fmt.Println()
							}
							if done != "" && err == nil {
								err = notion.MarkAs(token, stackID, done, true)
								fmt.Println()
							}
							if undone != "" && err == nil {
								err = notion.MarkAs(token, stackID, undone, false)
								fmt.Println()
							}
							if show && err == nil {
								err = notion.ShowPage(token, stackID)
							}
							return err
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
