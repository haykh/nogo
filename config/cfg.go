package config2

import (
	// "encoding/base64"
	"fmt"
	"os"
	// "os/exec"

	"github.com/haykh/nogo/utils"

	"github.com/BurntSushi/toml"
	// "github.com/haykh/goencode"
)

// type Config struct {
// 	configPath string
// 	configFile string
// }
//
// func (c *Config) Fname() string {
// 	return fmt.Sprintf("%s%s", c.configPath, c.configFile)
// }
//
// type ParseTemplate struct {
// 	config_file Config
// 	configs     map[string]interface{}
// }

//	func (c *ParseTemplate) GetSecret(param string) (string, error) {
//		// encoding_key := base64.StdEncoding.EncodeToString([]byte(os.Getenv("USER")))
//		// api_fname, ok := c.configs["nogo_vault"].(string)
//		// if !ok {
//		// 	return "", fmt.Errorf("`nogo_vault` not found in config file")
//		// }
//		// v := goencode.File(encoding_key, api_fname)
//		// return v.Get(param)
//		key, err := exec.Command("pass", "nogo-token").Output()
//		if err != nil {
//			return "", fmt.Errorf(fmt.Sprintf("error: %s", err))
//		}
//		return string(key), nil
//	}
//
//	func (c *ParseTemplate) SetSecret(param, newvalue string) error {
//		cmd := exec.Command("pass", "insert", "nogo-token")
//		stdin, err := cmd.StdinPipe()
//		if err != nil {
//			return err
//		}
//
//		go func() {
//			defer stdin.Close()
//		}()
//
//		out, err := cmd.CombinedOutput()
//		if err != nil {
//			return err
//		}
//
//		fmt.Printf("%s\n", out)
//		return nil
//		// encoding_key := base64.StdEncoding.EncodeToString([]byte(os.Getenv("USER")))
//		// api_fname, ok := c.configs["nogo_vault"].(string)
//		// if !ok {
//		// 	return fmt.Errorf("`nogo_vault` not found in config file")
//		// }
//		// v := goencode.File(encoding_key, api_fname)
//		// if err := v.Set(param, newvalue); err != nil {
//		// 	return err
//		// }
//		// return nil
//	}
func WriteToFile(cfg Config) error {
	g_fname := os.Getenv("HOME") + "/.config/nogo/config.toml"
	if _, exists := os.Stat(g_fname); os.IsNotExist(exists) {
		return fmt.Errorf("%sconfig file does not exist%s", utils.ColorRed, utils.ColorReset)
	} else {
		if err := os.Remove(g_fname); err != nil {
			return err
		}
		if err := utils.CreateFile(g_fname); err != nil {
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

//
// func (p *ParseTemplate) ReadOrUpdateParameter(param string, default_value string) error {
// 	defer func() {
// 		fmt.Println()
// 	}()
// 	if v, exists := p.configs[param]; !exists {
// 		p.configs[param] = default_value
// 		if value, err := utils.PromptString(fmt.Sprintf("`%s` not found\n  %s\n\nenter new path or leave blank for default", param, p.config_file.Fname()), default_value); err != nil {
// 			return err
// 		} else {
// 			if value != "" {
// 				p.configs[param] = value
// 			}
// 		}
// 		return p.WriteToFile()
// 	} else {
// 		v_str := v.(string)
// 		if leave, err := utils.PromptString(fmt.Sprintf("`%s` found\nenter new value or leave blank to use existing value", param), v_str); err != nil {
// 			return err
// 		} else {
// 			if leave != "" {
// 				p.configs[param] = leave
// 				return p.WriteToFile()
// 			} else {
// 				return nil
// 			}
// 		}
// 	}
// }
//
// func StoreSecret(fname, param, value, key string) error {
// 	v := goencode.File(key, fname)
// 	return v.Set(param, value)
// }
//
// func AssertSecretStored(fname, param, key string) error {
// 	v := goencode.File(key, fname)
// 	_, err := v.Get(param)
// 	return err
// }
//
// func StoreOrCheckSecret(fname, param, description string) error {
// 	encoding_key := base64.StdEncoding.EncodeToString([]byte(os.Getenv("USER")))
// 	defer func() {
// 		AssertSecretStored(fname, param, encoding_key)
// 	}()
// 	makeNew := func() error {
// 		if newparam, err := utils.PromptString(fmt.Sprintf("enter your %s:", description), ""); err != nil {
// 			return err
// 		} else {
// 			return StoreSecret(fname, param, newparam, encoding_key)
// 		}
// 	}
// 	if _, exists := os.Stat(fname); os.IsNotExist(exists) {
// 		return makeNew()
// 	} else {
// 		v := goencode.File(encoding_key, fname)
// 		_, err := v.Get(param)
// 		v.List()
// 		if err != nil {
// 			return makeNew()
// 		} else {
// 			if overwrite, err := utils.PromptBool(fmt.Sprintf("vault contains `%s`\noverwrite?", param), false); err != nil {
// 				return err
// 			} else {
// 				if overwrite {
// 					if err := v.Delete(param); err != nil {
// 						return err
// 					}
// 					return makeNew()
// 				}
// 				return nil
// 			}
// 		}
// 	}
// }

// func ReadOrPutConfig(
// 	cfg Config,
// 	silent bool,
// 	desc, param string,
// 	def interface{},
// ) (interface{}, error) {
// 	if v, exists := cfg[param]; !exists {
// 		if silent {
// 			return nil, fmt.Errorf(fmt.Sprintf("`%s` is not defined in config", param))
// 		}
// 		cfg[param] = ""
// 		if value, err := utils.Prompt(
// 			fmt.Sprintf(
// 				"`%s` : %s.\n parameter not set in config.\n",
// 				param,
// 				desc,
// 			),
// 			false,
// 			def,
// 		); err != nil {
// 			return nil, err
// 		} else {
// 			if value != "" {
// 				cfg[param] = value
// 			}
// 		}
// 		return cfg[param], WriteToFile(cfg)
// 	} else {
// 		if silent {
// 			return cfg[param], nil
// 		}
// 		if leave, err := utils.Prompt(
// 			fmt.Sprintf("`%s` set in config.", param),
// 			true,
// 			v,
// 		); err != nil {
// 			return nil, err
// 		} else {
// 			if leave != "" {
// 				cfg[param] = leave
// 				return cfg[param], WriteToFile(cfg)
// 			} else {
// 				return cfg[param], nil
// 			}
// 		}
// 	}
// }
//
// var localConfig = Config{
// 	configPath: os.Getenv("HOME") + "/.config/nogo/",
// 	configFile: "config.toml",
// }

// type Config = map[string]interface{}

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

func (c Config) EnsureExists(silent bool) error {
	if _, exists := os.Stat(c.fname); os.IsNotExist(exists) {
		if silent {
			return fmt.Errorf("config file does not exist")
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

// func CreateOrParseConfig(silent bool) (Config, error) {
// 	config_fname := os.Getenv("HOME") + "/.config/nogo/config.toml"
// 	configs := Config{}
// 	if _, exists := os.Stat(config_fname); os.IsNotExist(exists) {
// 		if silent {
// 			return nil, fmt.Errorf("config file does not exist")
// 		}
// 		utils.Message(
// 			fmt.Sprintf("config file does not exist. creating @ %s\n", config_fname),
// 			utils.Normal,
// 			true,
// 		)
// 		if err := utils.CreateFile(config_fname); err != nil {
// 			return nil, err
// 		}
// 	} else if silent {
// 		if _, err := toml.DecodeFile(config_fname, &configs); err != nil {
// 			return nil, err
// 		}
// 		return configs, nil
// 	}
// 	utils.Message(
// 		fmt.Sprintf("reading config file @ %s\n", config_fname),
// 		utils.Normal,
// 		true,
// 	)
// 	if _, err := toml.DecodeFile(config_fname, &configs); err != nil {
// 		return nil, err
// 	}
// 	ReadOrPutConfig(
// 		configs,
// 		silent,
// 		"use legacy mode for encoding the Notion token when saving locally (not safe)",
// 		"legacy_encode",
// 		false,
// 	)
// 	ReadOrPutConfig(
// 		configs,
// 		silent,
// 		"Notion ID for the stack (todo-list) block",
// 		"stack_id",
// 		"",
// 	)
// 	return configs, nil
// 	// if err := parsed_l_config.ReadOrUpdateParameter("nogo_vault", parsed_l_config.config_file.configPath+"nogo_vault"); err != nil {
// 	// 	return nil, err
// 	// }
// 	//
// 	// if secret_fname, ok := parsed_l_config.configs["nogo_vault"].(string); !ok {
// 	// 	return nil, fmt.Errorf("undefined API token file")
// 	// } else {
// 	// 	StoreOrCheckSecret(secret_fname, "api_token", "Notion API token")
// 	// 	StoreOrCheckSecret(secret_fname, "stack_page_id", "Stack page ID")
// 	// }
// 	// return parsed_l_config, nil
// }

// type LocalParseTemplate struct {
// 	ParseTemplate
// }
//
// func CreateOrReadLocalConfig(silent bool) (LocalParseTemplate, error) {
// 	parsed_l_config := LocalParseTemplate{
// 		ParseTemplate{
// 			config_file: localConfig,
// 			configs:     map[string]interface{}{},
// 		},
// 	}
// 	l_fname := parsed_l_config.config_file.Fname()
// 	if _, exists := os.Stat(l_fname); os.IsNotExist(exists) {
// 		if silent {
// 			return LocalParseTemplate{}, fmt.Errorf("config file does not exist")
// 		}
// 		utils.Message(fmt.Sprintf("local config file does not exist. creating...\n  %s", l_fname), utils.Normal, true)
// 		if err := utils.CreateFile(l_fname); err != nil {
// 			return LocalParseTemplate{}, err
// 		}
// 	} else if silent {
// 		if _, err := toml.DecodeFile(l_fname, &parsed_l_config.configs); err != nil {
// 			return LocalParseTemplate{}, err
// 		}
// 		return parsed_l_config, nil
// 	}
// 	utils.Message(fmt.Sprintf("Reading local config file...\n  %s", l_fname), utils.Normal, true)
// 	if _, err := toml.DecodeFile(l_fname, &parsed_l_config.configs); err != nil {
// 		return LocalParseTemplate{}, err
// 	}
// 	if err := parsed_l_config.ReadOrUpdateParameter("nogo_vault", parsed_l_config.config_file.configPath+"nogo_vault"); err != nil {
// 		return LocalParseTemplate{}, err
// 	}
//
// 	if secret_fname, ok := parsed_l_config.configs["nogo_vault"].(string); !ok {
// 		return LocalParseTemplate{}, fmt.Errorf("undefined API token file")
// 	} else {
// 		StoreOrCheckSecret(secret_fname, "api_token", "Notion API token")
// 		StoreOrCheckSecret(secret_fname, "stack_page_id", "Stack page ID")
// 	}
// 	return parsed_l_config, nil
// }
