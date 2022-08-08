package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"nogo/utils"
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
		return "", fmt.Errorf("`nogo_secret_file` not found in config file.")
	}
	v := goencode.File(encoding_key, api_fname)
	return v.Get(param)
}

// func (c *ParseTemplate) GetPageID() string {}

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
		p.configs[param] = utils.PromptString(fmt.Sprintf("`%s` not found in %s.\nEnter new value or leave blank for default.", param, p.config_file.Fname()), default_value, utils.Normal)
		if p.configs[param] == "" {
			p.configs[param] = default_value
		}
		p.WriteToFile()
	} else {
		v_str := v.(string)
		leave := utils.PromptString(fmt.Sprintf("`%s` found in %s.\nEnter new value or leave blank to use existing value.", param, p.config_file.Fname()), v_str, utils.Normal)
		if leave != "" {
			p.configs[param] = leave
			p.WriteToFile()
		}
	}
}

func StoreOrCheckSecret(fname, param, description string) {
	encoding_key := base64.StdEncoding.EncodeToString([]byte(os.Getenv("USER")))
	defer func() {
		v := goencode.File(encoding_key, fname)
		if _, err := v.Get(param); err != nil {
			log.Fatal(err)
		}
	}()
	makeNew := func() {
		// utils.Message(fmt.Sprintf("%s will be stored in the encrypted %s file.\nWhile this is slightly more secure than storing as flat text,\n... it is not entirely secure, so keep this file private.", description, fname), utils.Warning, true)
		newparam := utils.PromptString(fmt.Sprintf("Enter your %s:", description), "", utils.Normal)
		v := goencode.File(encoding_key, fname)
		if err := v.Set(param, newparam); err != nil {
			log.Fatal(err)
		}
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
			overwrite := utils.PromptBool(fmt.Sprintf("nogo secret file `%s` contains `%s`.\nOverwrite with the new parameter?", fname, param), false, utils.Normal)
			if overwrite {
				v.Delete(param)
				makeNew()
			}
		}
	}
}

var localConfig = Config{
	//!DEBUG:
	//configPath: os.Getenv("HOME") + "/.config/nogo/",
	configPath: "./",
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
		msg := fmt.Sprintf("Local configuration file `%s` does not exist.\nCreating?", l_fname)
		create := utils.PromptBool(msg, true, utils.Warning)
		if create {
			utils.CreateFile(l_fname)
		} else {
			utils.Message("Local configuration file not created; local defaults will not be saved.", utils.Warning, true)
			utils.Message(fmt.Sprintf("You can still use the tool if you:\n(1) provide a configuration file via command line arguments;\n\t<EXAMPLE>\n(2) use a file in the default directory: %s.\n\t<EXAMPLE>", localConfig.Fname()), utils.Hint, true)
			return parsed_l_config
		}
	} else if silent {
		toml.DecodeFile(l_fname, &parsed_l_config.configs)
		return parsed_l_config
	}
	utils.Message(fmt.Sprintf("Reading local configuration file: %s.", l_fname), utils.Normal, true)
	toml.DecodeFile(l_fname, &parsed_l_config.configs)
	parsed_l_config.ReadOrUpdateParameter("nogo_secret_file", parsed_l_config.config_file.configPath+"nogo_secret")

	if secret_fname, ok := parsed_l_config.configs["nogo_secret_file"].(string); !ok {
		log.Fatal("Undefined API token file.")
	} else {
		StoreOrCheckSecret(secret_fname, "api_token", "Notion API token")
		StoreOrCheckSecret(secret_fname, "main_page_id", "Main page ID")
	}
	return parsed_l_config
}

// var globalConfig = Config{
// 	//!DEBUG:
// 	//config_path: "/usr/local/etc/nogo/",
// 	configPath: "./",
// 	configFile: "global.toml",
// }

// type GlobalParseTemplate struct {
// 	ParseTemplate
// }

// func CreateOrReadGlobalConfig() GlobalParseTemplate {
// 	parsed_g_config := GlobalParseTemplate{
// 		ParseTemplate{
// 			config_file: globalConfig,
// 			configs: map[string]interface{}{
// 				"ConfigPath": "",
// 			},
// 		},
// 	}
// 	g_fname := parsed_g_config.config_file.Fname()
// 	if _, exists := os.Stat(g_fname); os.IsNotExist(exists) {
// 		utils.CreateFile(g_fname)
// 	}
// 	toml.DecodeFile(g_fname, &parsed_g_config)
// 	if parsed_g_config.configs["ConfigPath"] == "" {
// 		parsed_g_config.configs["ConfigPath"] = localConfig.Fname()
// 		parsed_g_config.WriteToFile()
// 	}
// 	return parsed_g_config
// }
