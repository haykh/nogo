package config

import (
	"encoding/base64"
	"fmt"
	"log"
	utils "nogo/utils"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/haykh/goencode"
)

type Config struct {
	configPath string
	configFile string
}

func (c *Config) Fname() string {
	return fmt.Sprintf("%s%s", c.configPath, c.configFile)
}

type ParseTemplate struct {
	config_file Config
	configs     map[string]interface{}
}

func (c *ParseTemplate) GetSecret(param string) (string, error) {
	encoding_key := base64.StdEncoding.EncodeToString([]byte(os.Getenv("USER")))
	api_fname, ok := c.configs["nogo_secret_file"].(string)
	if !ok {
		return "", fmt.Errorf("`nogo_secret_file` not found in config file")
	}
	v := goencode.File(encoding_key, api_fname)
	return v.Get(param)
}

func (c *ParseTemplate) SetSecret(param, newvalue string) error {
	encoding_key := base64.StdEncoding.EncodeToString([]byte(os.Getenv("USER")))
	api_fname, ok := c.configs["nogo_secret_file"].(string)
	if !ok {
		return fmt.Errorf("`nogo_secret_file` not found in config file")
	}
	v := goencode.File(encoding_key, api_fname)
	if err := v.Set(param, newvalue); err != nil {
		return err
	}
	return nil
}

func (p *ParseTemplate) WriteToFile() {
	g_fname := p.config_file.Fname()
	if _, exists := os.Stat(g_fname); os.IsNotExist(exists) {
		log.Fatal("Config file does not exist.")
	} else {
		if err := os.Remove(g_fname); err != nil {
			log.Fatal(err)
		}
		utils.CreateFile(g_fname)
	}
	f, err := os.OpenFile(g_fname, os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}
	if err := toml.NewEncoder(f).Encode(p.configs); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func (p *ParseTemplate) ReadOrUpdateParameter(param string, default_value string) {
	if v, exists := p.configs[param]; !exists {
		p.configs[param] = utils.PromptString(fmt.Sprintf("`%s` not found\n  %s\n\nEnter new path or leave blank for default", param, p.config_file.Fname()), default_value, utils.Normal, utils.ColorRed)
		if p.configs[param] == "" {
			p.configs[param] = default_value
		}
		p.WriteToFile()
		fmt.Println()
	} else {
		v_str := v.(string)
		leave := utils.PromptString(fmt.Sprintf("`%s` found\nEnter new value or leave blank to use existing value", param), v_str, utils.Normal, utils.ColorBlue)
		if leave != "" {
			p.configs[param] = leave
			p.WriteToFile()
		}
		fmt.Println()
	}
}

func StoreSecret(fname, param, value, key string) {
	v := goencode.File(key, fname)
	if err := v.Set(param, value); err != nil {
		log.Fatal(err)
	}
}

func AssertSecretStored(fname, param, key string) {
	v := goencode.File(key, fname)
	if _, err := v.Get(param); err != nil {
		log.Fatal(err)
	}
}

func StoreOrCheckSecret(fname, param, description string) {
	encoding_key := base64.StdEncoding.EncodeToString([]byte(os.Getenv("USER")))
	defer func() {
		AssertSecretStored(fname, param, encoding_key)
	}()
	makeNew := func() {
		newparam := utils.PromptString(fmt.Sprintf("Enter your %s:", description), "", utils.Normal)
		StoreSecret(fname, param, newparam, encoding_key)
	}
	if _, exists := os.Stat(fname); os.IsNotExist(exists) {
		makeNew()
	} else {
		v := goencode.File(encoding_key, fname)
		_, err := v.Get(param)
		v.List()
		if err != nil {
			makeNew()
		} else {
			overwrite := utils.PromptBool(fmt.Sprintf("nogo secret file `%s` contains `%s`\nOverwrite with the new parameter?", fname, param), false, utils.Normal, utils.ColorBlue)
			if overwrite {
				v.Delete(param)
				makeNew()
			}
			fmt.Println()
		}
	}
}

var localConfig = Config{
	configPath: os.Getenv("HOME") + "/.config/nogo/",
	configFile: "config.toml",
}

type LocalParseTemplate struct {
	ParseTemplate
}

func CreateOrReadLocalConfig(silent bool) LocalParseTemplate {
	parsed_l_config := LocalParseTemplate{
		ParseTemplate{
			config_file: localConfig,
			configs:     map[string]interface{}{},
		},
	}
	l_fname := parsed_l_config.config_file.Fname()
	if _, exists := os.Stat(l_fname); os.IsNotExist(exists) {
		if silent {
			log.Fatal("Config file does not exist.")
		}
		utils.Message(fmt.Sprintf("Local configuration file does not exist\n  Creating:\n  %s\n", l_fname), utils.Normal, true, utils.ColorRed)
		utils.CreateFile(l_fname)
	} else if silent {
		toml.DecodeFile(l_fname, &parsed_l_config.configs)
		return parsed_l_config
	}
	utils.Message(fmt.Sprintf("Reading local configuration file:\n  %s", l_fname), utils.Normal, true)
	toml.DecodeFile(l_fname, &parsed_l_config.configs)
	parsed_l_config.ReadOrUpdateParameter("nogo_secret_file", parsed_l_config.config_file.configPath+"nogo_secret")

	if secret_fname, ok := parsed_l_config.configs["nogo_secret_file"].(string); !ok {
		log.Fatal("Undefined API token file.")
	} else {
		StoreOrCheckSecret(secret_fname, "api_token", "Notion API token")
		StoreOrCheckSecret(secret_fname, "stack_page_id", "Stack page ID")
	}
	return parsed_l_config
}
