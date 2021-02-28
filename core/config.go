package core

import (
	"encoding/json"
	io "io/ioutil"
)

// Config json configuration
var Config *Configuration

// Configuration json configuration struct
type Configuration struct {
	LogPath string `json:"log_path"`
}

// Loads load json configuration
func (v *Configuration) Loads(conf string) error {
	if f, err := io.ReadFile(conf); err != nil {
		return err
	} else {
		d := []byte(f)
		if err := json.Unmarshal(d, v); err != nil {
			return err
		}
	}
	LogPath = Config.LogPath
	// fmt.Println(Config)
	return nil
}

// NewConfig create new configuration
func NewConfig() *Configuration {
	c := &Configuration{}
	return c
}
