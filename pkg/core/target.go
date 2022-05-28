package core

import (
	"strconv"
	"strings"
)

type Target struct {
	Name           string  `json:"name"`
	Method         string  `json:"method"`
	Address        string  `json:"address"`
	Port           int     `json:"port"`
	Status         string  `json:"status"`
	ExtendedStatus string  `json:"extendedStatus"`
	Delay          float64 `json:"delay"`
}

func (t *Target) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// First see if it is a string definition.
	var str string
	if err := unmarshal(&str); err == nil {
		parts := strings.SplitN(str, ":", 4)
		if len(parts) == 1 {
			t.Method = "tcp"
			t.Address = parts[0]
			t.Port = 80
		} else {
			if len(parts) >= 1 {
				t.Method = parts[0]
			}
			if len(parts) >= 2 {
				t.Address = parts[1]
			}
			if len(parts) >= 3 {
				t.Port, _ = strconv.Atoi(parts[2])
			} else {
				if t.Method == "tcp" {
					t.Port = 80
				} else if t.Method == "udp" {
					t.Port = 53
				} else if t.Method == "tls" {
					t.Port = 443
				}
			}
			if len(parts) >= 4 {
				t.Name = parts[3]
			}
		}
		return nil
	}

	// Otherwise try to unmarshal to Target.
	if err := unmarshal(t); err != nil {
		return err
	}

	return nil
}
