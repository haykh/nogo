package config

import (
	// "encoding/base64"
	"fmt"
	"os"
	// "os/exec"

	"github.com/haykh/nogo/utils"

	"github.com/BurntSushi/toml"
	// "github.com/haykh/goencode"
)

type Config struct {
	fname  string
	stored map[string]utils.Configuration
}

var MainConfig = Config{
	fname: os.Getenv("HOME") + "/.config/nogo/config.toml",
	stored: map[string]utils.Configuration{
		"legacy_secret": utils.NewParameter(
			"legacy_secret",
			"use legacy mode for encoding the Notion token when saving locally (not safe)",
			nil,
			false,
		),
		"stack_id": utils.NewConfiguration(
			"stack_id",
			"Notion ID for the stack (todo-list) block",
			nil,
		),
		"notion_token": utils.NewSecret(
			"notion_token",
			"Notion API token",
			nil,
		),
	},
}

func (c Config) WriteToFile() error {
	if _, exists := os.Stat(c.fname); os.IsNotExist(exists) {
		return fmt.Errorf("%sconfig file does not exist%s", utils.ColorRed, utils.ColorReset)
	} else {
		if err := os.Remove(c.fname); err != nil {
			return err
		}
		if err := utils.CreateFile(c.fname); err != nil {
			return err
		}
	}
	if f, err := os.OpenFile(g_fname, os.O_WRONLY, 0777); err != nil {
		return err
	} else {
		if err := toml.NewEncoder(f).Encode(cfg); err != nil {
			return err
		} else {
			if err := f.Close(); err != nil {
				return err
			} else {
				return nil
			}
		}
	}
}

func (c Config) EnsureFileExists(silent bool) error {
	if _, exists := os.Stat(c.fname); os.IsNotExist(exists) {
		if silent {
			return fmt.Errorf("%sconfig file does not exist%s", utils.ColorRed, utils.ColorReset)
		}
		utils.Message(
			fmt.Sprintf("config file does not exist. creating @ %s\n", c.fname),
			utils.Normal,
			true,
		)
		if err := utils.CreateFile(c.fname); err != nil {
			return err
		}
	} else if !silent {
		utils.Message(
			fmt.Sprintf("reading config file @ %s\n", c.fname),
			utils.Normal,
			true,
		)
	}
	return nil
}

func (c *Config) EnsureDefined(param string, silent bool) error {
	if _, exists := c.stored[param]; !exists {
		return fmt.Errorf(fmt.Sprintf("`%s` not defined\n", param))
	}
	if c.stored[param].Value() != nil {
		if !silent {
			utils.Message(
				fmt.Sprintf("`%s` found in config\n", param),
				utils.Normal,
				true,
			)
		}
		return nil
	} else if silent {
		return fmt.Errorf(fmt.Sprintf("no value defined for `%s`\n", param))
	}
	return c.Prompt()
}
