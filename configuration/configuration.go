package configuration

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

//Configuration app configuration
type Configuration struct {
	Addr      string   `mapstructure:"addr"`
	APIPrefix string   `mapstructure:"api_prefix"`
	Timeout   uint     `mapstructure:"timeout"`
	Storage   *Storage `mapstructure:"storage"`
}

//Storage configuration of storage
type Storage struct {
	URI            string `mapstructure:"uri"`
	DBName         string `mapstructure:"db_name"`
	CollectionName string `mapstructure:"collection_name"`
	Timeout        uint   `mapstructure:"timeout"`
}

//NewConfiguration reads configuration from file and environment variables
func NewConfiguration(configFilePath string) (*Configuration, error) {
	viper.SetConfigFile(configFilePath)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("KNIHOMOL")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return nil, fmt.Errorf("unable to BindPFlags, %s", err)
	}

	viper.SetDefault("addr", "0.0.0.0:80")
	viper.SetDefault("api_prefix", "/v1")
	viper.SetDefault("timeout", "3")
	viper.SetDefault("storage.uri", "mongodb://localhost:27017")
	viper.SetDefault("storage.db_name", "knihomol")
	viper.SetDefault("storage.collection_name", "books")
	viper.SetDefault("storage.timeout", "3")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			return nil, fmt.Errorf("failed to read config file: %s", err)
		}
	}

	var configuration Configuration
	err = viper.Unmarshal(&configuration)
	if err != nil {
		return nil, fmt.Errorf("unable to decode configration to struct, %s", err)
	}

	return &configuration, nil
}

//PrintConfiguration prints configuration to stdout
func (c *Configuration) PrintConfiguration() {
	data, err := yaml.Marshal(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to print configuration: %s", err)
		return
	}
	fmt.Printf("%s\n", data)
}
