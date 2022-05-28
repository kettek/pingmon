package core

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Title struct {
	Prefix *string `yaml:"prefix"`
	Name   *string `yaml:"name"`
	Suffix *string `yaml:"suffix"`
}

type Config struct {
	Rate           int        `yaml:"rate"`
	Timeout        *int       `yaml:"timeout"`
	Targets        *[]*Target `yaml:"targets"`
	Address        *string    `yaml:"address"`
	Assets         *string    `yaml:"assets"`
	Title          *Title     `yaml:"title"`
	PrivilegedPing *bool      `yaml:"privilegedPing"`
}

var DefaultConfig Config = Config{}

func init() {
	addr := ":8999"
	assets := "./pkg/frontend/public"
	timeout := 5
	privilegedPing := true
	DefaultConfig.Address = &addr
	DefaultConfig.Assets = &assets
	DefaultConfig.Rate = 30
	DefaultConfig.Timeout = &timeout
	DefaultConfig.PrivilegedPing = &privilegedPing

	prefix := ""
	name := "pingmon"
	suffix := ""
	DefaultConfig.Title = &Title{
		Prefix: &prefix,
		Name:   &name,
		Suffix: &suffix,
	}
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
	if c2.Title != nil {
		if c2.Title.Name != nil {
			c.Title.Name = c2.Title.Name
		}
		if c2.Title.Prefix != nil {
			c.Title.Prefix = c2.Title.Prefix
		}
		if c2.Title.Suffix != nil {
			c.Title.Suffix = c2.Title.Suffix
		}
	}
	if c2.PrivilegedPing != nil {
		c.PrivilegedPing = c2.PrivilegedPing
	}

	return nil
}
