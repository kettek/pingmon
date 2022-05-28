package core

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/go-ping/ping"
)

func RefreshTargets(c *Config) {
	if c.Targets == nil {
		return
	}
	for _, t := range *c.Targets {
		if t.Method == "tcp" || t.Method == "udp" {
			start := time.Now()
			_, err := net.DialTimeout(t.Method, fmt.Sprintf("%s:%d", t.Address, t.Port), time.Duration(*c.Timeout)*time.Second)
			if err != nil {
				t.Status = "offline"
				t.ExtendedStatus = err.Error()
				switch err := err.(type) {
				case *net.OpError:
					switch err := err.Err.(type) {
					case *net.DNSError:
						t.Status = "error"
						t.ExtendedStatus = err.Error()
					default:
						t.Status = "offline"
						t.ExtendedStatus = err.Error()
					}
				}
			} else {
				t.Status = "online"
				t.ExtendedStatus = ""
			}
			t.Delay = float64(time.Now().Sub(start).Microseconds())
		} else if t.Method == "tls" {
			start := time.Now()
			dialer := net.Dialer{
				Timeout: time.Duration(*c.Timeout) * time.Second,
			}
			_, err := tls.DialWithDialer(&dialer, "tcp", fmt.Sprintf("%s:%d", t.Address, t.Port), nil)

			if err != nil {
				switch err := err.(type) {
				case *net.OpError:
					switch err := err.Err.(type) {
					case *net.DNSError:
						t.Status = "error"
						t.ExtendedStatus = err.Error()
					default:
						t.Status = "offline"
						t.ExtendedStatus = err.Error()
					}
				default:
					if strings.Contains(err.Error(), "certificate has expired") {
						t.Status = "expired"
					} else if strings.Contains(err.Error(), "signed by unknown") {
						t.Status = "untrusted"
					} else {
						t.Status = "offline"
					}
					t.ExtendedStatus = err.Error()
				}
			} else {
				t.Status = "online"
				t.ExtendedStatus = ""
			}

			t.Delay = float64(time.Now().Sub(start).Microseconds())
		} else if t.Method == "ping" {
			pinger, err := ping.NewPinger(t.Address)
			pinger.SetPrivileged(*c.PrivilegedPing)
			if err != nil {
				t.Status = "error"
				continue
			}
			pinger.Timeout = time.Duration(*c.Timeout) * time.Second
			pinger.Count = 1
			pinger.Interval = 1
			start := time.Now()
			if err := pinger.Run(); err != nil {
				fmt.Println(err)
				t.Status = "offline"
			} else if pinger.Statistics().PacketsRecv == 0 {
				t.Status = "offline"
			} else {
				t.Status = "online"
			}
			t.Delay = float64(time.Now().Sub(start).Microseconds())
		}
	}
}
