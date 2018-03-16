package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Region            string `viper:"region"`
	Profile           string `viper:"profile"`
	AwsConfigPath     string `viper:aws_config_path`
	AwsCredentialPath string `viper:aws_credential_path`
	SecretKey         string `viper:secret_key`
	AccessKey         string `viper:secret_key`

	// congfig key value for gralde.properties
	GradleAccessKey string `viper:gradle_secret_key`
	GradleSecretKey string `viper:gradle_secret_key`
}

const (
	pathConfig = ".ark"
)

func getConfig() *Config {
	viper.AddConfigPath("$HOME/" + pathConfig)
	viper.SetConfigName("config")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file must exist in ~/"+pathConfig+"./config.yaml: %s \n", err))
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		panic(fmt.Errorf("Fatal error unmarhsal config struct : %s \n", err))
	}

	return config
}
