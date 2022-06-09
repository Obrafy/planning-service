package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

const (
	ConfigDirectory       = "/configuration/"
	DefaultConfigFileName = "default.yaml"
)

func loadOneFile(path string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var c Configuration

	err = yaml.Unmarshal(bytes, &c)

	return &c, err
}

func merge(a, b *Configuration) (*Configuration, error) {
	if err := mergo.Merge(a, b); err != nil {
		return nil, err
	}

	return a, nil
}

func getAppLocation() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func getDefaultConfigurationFileName() (string, error) {
	dir, err := getAppLocation()
	return dir + ConfigDirectory + DefaultConfigFileName, err
}

func getEnvironmentSpecificConfigurationFileName(environment string) (string, error) {
	dir, err := getAppLocation()
	return dir + ConfigDirectory + environment + ".yaml", err
}

func LoadConfiguration(environment string) (*Configuration, error) {
	fmt.Println("Loading configuration for environment ", environment)

	defaultConfigurationFile, err := getDefaultConfigurationFileName()

	if err != nil {
		return nil, err
	}

	configDefault, err := loadOneFile(defaultConfigurationFile)

	if err != nil {
		return nil, err
	}

	if len(environment) == 0 {
		fmt.Println("Loaded configuration: ", defaultConfigurationFile)
		return configDefault, nil
	}

	environmentSpecificConfigurationFile, err := getEnvironmentSpecificConfigurationFileName(environment)

	if err != nil {
		return nil, err
	}

	environmentConfiguration, err := loadOneFile(environmentSpecificConfigurationFile)

	if err != nil {
		return nil, err
	}

	fmt.Println("Loaded configuration for environment:", environmentSpecificConfigurationFile)

	merged, err := merge(environmentConfiguration, configDefault)

	fmt.Println("Loaded configuration merge:", merged)

	return merged, err
}
