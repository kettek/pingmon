package cfg

import (
	"pingmon/pkg/pinger"
)

type Config struct {
	Targets []pinger.Target `yaml:"targets"`
	Address string          `yaml:"address"`
}
