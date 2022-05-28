# Introduction
pingmon is a *backend and frontend service*(**pingmon**), a *cli tool*(**pingmon-cli**), or a *system tray item*(**pingmon-systray**) that can ping a list of addresses to see if they are alive. Pinging can be done through TCP, UDP, or ICMP pings. A list of servers are defined through a configuration file that the backend pings over time and serves their status via a frontend HTTP service.

![pingmon browser view](screenshot.png)

# Building

## Requisites

  * node/npm of at least 2020
  * go of at least 1.17

## Build

pingmon uses the gobl build system. If you have go installed, both the backend and frontend can be built by issuing:

```
go run . buildBackend
go run . buildFrontend
```

Once this is done, you can run the corresponding `pingmon` program.

# Configuration
Configuration is done through a simple yaml file in the working directory named `cfg.yml`. The following example should give a good idea of how it works.

## Complete Example
```yaml
rate: 10                    // Ping every 10 seconds
timeout: 1                  // Timeout pings after 1 second
privilegedPing: true        // Use a privileged ping. See permissions.
targets:
  - myaddress               // defaults to tcp:myaddress:80
  - tcp:myaddress           // defaults to 80
  - udp:myaddress           // defaults to 53
  - tcp:myaddress:port
  - udp:myaddress:port
  - ping:myaddress          // send ICMP ping
  - tls:myaddress           // Attempt SSL/TLS handshaking. defaults to 443
  - tls:myaddress:port
address: "*:8888"           // Listen to on all addresses on port 8888
assets: pkg/frontend/public // Look for frontend assets in this directory.
title:                      // will print "> pingmon <" in the browser
  prefix: "> "
  name: "pingmon"
  suffix: " <"
```

# Service Installation

## Linux (gobl)
pingmon can be installed to the system and a systemd unit installed by simply issuing:

```
sudo go run . installSystemdUnit
```

After this, you can start and enable pingmon:

```
sudo systemctl enable --now pingmon
```

## Linux (manual)
If you wish to run as a service, you can copy the provided `extra/pingmon.service` into `/etc/systemd/system/` or otherwise. If you install pingmon into a different location instead of the default, then modify the service file accordingly.

File installation is most readily done from within the `pingmon` source directory:

```
sudo mkdir -p /opt/pingmon
sudo cp pingmon /opt/pingmon/
sudo cp -r pkg/frontend/public /opt/pingmon/public
```

Then create a `cfg.yml` file in `/opt/pingmon/` with the following:

```
assets: public
```

### Permissions
To run as a service, you may need to set the pingmon binary to use `CAP_NET_RAW`:

```
sudo setcap cap_net_raw=+ep pingmon
```

If this is set, ensure the `privilegedPing` option is set to `true`.

----

Or allow unprivileged UDP pings:

```
sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
```

If this is set, ensure the `privilegedPing` option is set to `false`.

