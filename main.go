package main

import (
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/haykh/nogo/config"

	notion "github.com/haykh/nogo/api"

	"github.com/urfave/cli/v2"
)

func main() {
	log.SetPrefix("[ nogo ERROR ]: ")
	log.SetFlags(0)
	cli.AppHelpTemplate = `{{.Name}} v{{.Version}} [by {{range .Authors}}{{ . }}{{end}}]

	 {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
	tmpls := []*string{
		&cli.CommandHelpTemplate,
		&cli.SubcommandHelpTemplate,
	}
	old := []string{
		"NAME:",
		"USAGE:",
		"COMMANDS:",
		"OPTIONS:",
		"[arguments...]",
		"options",
	}
	new := []string{
		"name:",
		"usage:",
		"commands:",
		"options:",
		"[args...]",
		"opts",
	}
	for _, tmpl := range tmpls {
		for i, o := range old {
			*tmpl = strings.ReplaceAll(*tmpl, o, new[i])
		}
	}

	app := &cli.App{
		Name:     "nogo",
		Version:  "1.6.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{{
			Name: "@haykh",
		}},
		Usage: "do awesome stuff with notion from a cli",
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
						return err
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
								return err
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
								return err
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
								return err
							} else {
								return notion.ToggleStack(client, sID)
							}
						},
					},
					{
						Name:    "rnd",
						Usage:   "select a random unfinished task from the stack",
						Action: func(cCtx *cli.Context) error {
							if client, sID, err := notion.InitAPI(); err != nil {
								return err
							} else {
                return notion.RandomStackEntry(client, sID)
							}
						},
					},
					{
						Name:    "rm",
						Aliases: []string{"r"},
						Usage:   "remove stack entries",
						Action: func(cCtx *cli.Context) error {
							if client, sID, err := notion.InitAPI(); err != nil {
								return err
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
