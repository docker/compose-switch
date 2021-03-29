package redirect

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	Enabled   = "enabled"
	Disabled  = "disabled"
	ComposeV2 = "composeV2"
)

type config struct {
	path string
	data map[string]interface{}
}

func (config *config) Load() error {
	config.data = map[string]interface{}{}
	if _, err := os.Stat(config.path); os.IsNotExist(err) {
		return nil
	}
	raw, err := ioutil.ReadFile(config.path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, &config.data)
	return err
}

func (config *config) Write() error {
	d, err := json.MarshalIndent(config.data, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(config.path, d, 0644)
}

func GetConfigFile() config {
	configDir := os.Getenv("DOCKER_CONFIG")
	if configDir == "" {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".docker")
	}
	configFile := filepath.Join(configDir, "features.json")
	return config{
		path: configFile,
	}
}

func (config *config) GetFeature(key string) (string, bool) {
	v, ok := config.data[key]
	if !ok {
		return "", false
	}
	return v.(string), true
}

func (config *config) SetFeature(key, value string) {
	config.data[key] = value
}
