package main

import (
  "fmt"
  "log"
  "os"
  "sort"

  "nogo/config"
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
          fmt.Println("- - - - - ( nogo configure wizard ⚙ ) - - - - -")
          a := config.CreateOrReadGlobalConfig()
          fmt.Println(a)
          return nil
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
