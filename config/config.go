package config

import (
  "fmt"
  "os"
  "log"
  "nogo/utils"

  "github.com/BurntSushi/toml"
)

func Babadook() {
  fmt.Println("babadook")
}

type Config struct {
  configPath string
  configFile string
}

func (c *Config) Fname() string {
  return c.configPath + c.configFile
}

var globalConfig = Config{
  //!DEBUG:
  //config_path: "/usr/local/etc/nogo/",
  configPath: "./",
  configFile: "global.toml",
}

var localConfig = Config{
  //!DEBUG:
  //configPath: os.Getenv("HOME") + "/.config/nogo/",
  configPath: "./",
  configFile: "config.toml",
}

type ParseGlobalConfigTemplate struct {
  ConfigPath string `toml:"config_path"`
}

func (p *ParseGlobalConfigTemplate) EncodeTemplate() map[string]interface{} {
  return map[string]interface{}{
    "config_path": p.ConfigPath,
  }
}

func (p *ParseGlobalConfigTemplate) WriteToFile() {
  g_fname := globalConfig.Fname()
  f, err := os.OpenFile(g_fname, os.O_WRONLY, 0777)
  if err != nil {
    log.Fatal(err)
  }
  if err := toml.NewEncoder(f).Encode(p.EncodeTemplate()); err != nil {
    log.Fatal(err)
  }
  if err := f.Close(); err != nil {
    log.Fatal(err)
  }
}


/*
 * Checks if global configuration file exists:
 *   [exists] -> parses file, returns ParseGlobalConfigTemplate
 *   [not exists] -> prompts to create
 *     [create] -> creates file
 *     [not create] -> uses default local config
 *
 * returns parsed global config file and bool if default local_config is used
 * @returns ParseGlobalConfigTemplate, bool
*/
func CreateOrReadGlobalConfig() ParseGlobalConfigTemplate {
  var parsed_g_config ParseGlobalConfigTemplate
  g_fname := globalConfig.Fname()
  if _, exists := os.Stat(g_fname); os.IsNotExist(exists) {
    create := utils.PromptBool(fmt.Sprintf("[ nogo WARNING ]: Global configuration file %s does not exist\n[ nogo WARNING ]: Creating?", g_fname), true)
    if create {
      utils.CreateFile(g_fname)
    } else {
      utils.Warn("Global configuration file not created; global defaults will not be saved")
      utils.Hint(fmt.Sprintf("You can still use the tool if you (1) provide a configuration file via command line arguments\n... or (2) use a file in the default directory: %s\n\n    (1): SOME EXAMPLE HERE ...\n    (2): SOME OTHER EXAMPLE HERE", localConfig.Fname()))
      parsed_g_config.ConfigPath = localConfig.Fname()
      return parsed_g_config
    }
  }
  toml.DecodeFile(g_fname, &parsed_g_config)
  if parsed_g_config.ConfigPath == "" {
    parsed_g_config.ConfigPath = localConfig.Fname()
    parsed_g_config.WriteToFile()
  }
  return parsed_g_config
}

