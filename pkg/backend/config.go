package backend

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Targets *[]Target `yaml:"targets"`
	Address *string   `yaml:"address"`
	Assets  *string   `yaml:"assets"`
}

var DefaultConfig Config = Config{}

func init() {
	addr := ":8999"
	assets := "./pkg/frontend/public"
	DefaultConfig.Address = &addr
	DefaultConfig.Assets = &assets
}

func (c *Config) FromYAML(p string) error {
	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	c2 := Config{}
	err = yaml.Unmarshal(b, &c2)
	if err != nil {
		return err
	}

	if c2.Address != nil {
		c.Address = c2.Address
	}
	if c2.Assets != nil {
		c.Assets = c2.Assets
	}
	if c2.Targets != nil {
		c.Targets = c2.Targets
	}

	return nil
}
