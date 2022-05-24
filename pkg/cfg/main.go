package cfg

import (
	"os"
	"pingmon/pkg/pinger"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Targets []pinger.Target `yaml:"targets"`
	Address string          `yaml:"address"`
}

func (c *Config) FromYAML(p string) error {
	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, &c)
}
