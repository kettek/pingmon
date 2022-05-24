package backend

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Rate    int        `yaml:"rate"`
	Timeout *int       `yaml:"timeout"`
	Targets *[]*Target `yaml:"targets"`
	Address *string    `yaml:"address"`
	Assets  *string    `yaml:"assets"`
}

var DefaultConfig Config = Config{}

func init() {
	addr := ":8999"
	assets := "./pkg/frontend/public"
	timeout := 5
	DefaultConfig.Address = &addr
	DefaultConfig.Assets = &assets
	DefaultConfig.Rate = 30
	DefaultConfig.Timeout = &timeout
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
	if c2.Timeout != nil {
		c.Timeout = c2.Timeout
	}
	if c2.Rate != 0 {
		c.Rate = c2.Rate
	}

	return nil
}
