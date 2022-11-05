package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var defaultConf []byte = []byte(`
stage: "dev"

# https://github.com/uber-go/zap
logger:
  level: 0
  developement: true

# see https://github.com/golang-migrate/migrate
migration:
  source: "migration_source"
  url: "migration_url"

linebot:
  channel:
    secret: "your_secret"
    token: "your_token"

mongo:
  url: "mongo_url"
  database: "mongo_database"
`)

type ConfYaml struct {
	Stage  string `yaml:"stage"`
	Logger struct {
		Level        int  `yaml:"level"`
		Developement bool `yaml:"developement"`
	} `yaml:"logger"`

	Migration struct {
		Source string `yaml:"source"`
		URL    string `yaml:"url"`
	} `yaml:"migration"`

	Linebot struct {
		Channel struct {
			Secret string `yaml:"secret"`
			Token  string `yaml:"token"`
		} `yaml:"channel"`
	} `yaml:"linebot"`

	Mongo struct {
		URL      string `yaml:"url"`
		Database string `yaml:"database"`
	} `yaml:"mongo"`
}

func LoadConf() (*ConfYaml, error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("myapp") // will be uppercased automatically
	viper.AutomaticEnv()        // read in environment variables that match

	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // use config in current path
	viper.SetConfigName("config")

	conf := &ConfYaml{}
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
		return conf, err
	}

	conf.Stage = viper.GetString("stage")

	conf.Logger.Level = viper.GetInt("logger.level")
	conf.Logger.Developement = viper.GetBool("logger.developement")

	conf.Migration.Source = viper.GetString("migration.source")
	conf.Migration.URL = viper.GetString("migration.url")

	conf.Linebot.Channel.Secret = viper.GetString("linebot.channel.secret")
	conf.Linebot.Channel.Token = viper.GetString("linebot.channel.token")

	conf.Mongo.Database = viper.GetString("mongo.database")
	conf.Mongo.URL = viper.GetString("mongo.url")

	return conf, nil
}
