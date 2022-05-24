# Building
pingmon uses the gobl build system. If you have go installed, both the backend and frontend can be built by issuing:

```
go run . buildBackend
go run . buildFrontend
```

Once this is done, you can run the corresponding `pingmon` program.

# Configuring

TODO

# Service Installation

## Linux
If you wish to run as a service, you can copy the provided `extra/pingmon.service` into `/etc/systemd/system/` or otherwise. If you installing pingmon into a different location than the default, then modify the service file accordingly.


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

To install as a service, you need to set the pingmon binary to use `CAP_NET_RAW`:

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

