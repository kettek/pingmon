package core

import (
	"strconv"
	"strings"
)

type Target struct {
	Method         string  `json:"method"`
	Address        string  `json:"address"`
	Port           int     `json:"port"`
	Status         string  `json:"status"`
	ExtendedStatus string  `json:"extendedStatus"`
	Delay          float64 `json:"delay"`
}

func (t *Target) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string

	if err := unmarshal(&str); err != nil {
		return err
	}

	parts := strings.SplitN(str, ":", 3)
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
	}

	return nil
}
