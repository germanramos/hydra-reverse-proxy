hydra-reverse-proxy
===================
[![Build Status](https://travis-ci.org/innotech/hydra-reverse-proxy.svg?branch=master)](https://travis-ci.org/innotech/hydra-reverse-proxy) [![Coverage Status](https://coveralls.io/repos/innotech/hydra-reverse-proxy/badge.png)](https://coveralls.io/r/innotech/hydra-reverse-proxy) [![GoDoc](https://godoc.org/github.com/innotech/hydra-reverse-proxy/reverse_proxy?status.png)](https://godoc.org/github.com/innotech/hydra-reverse-proxy/reverse_proxy)

Reverse proxy that takes an incoming request and sends it to hydra server, proxying the response back to the client.

## Installation

### Ubuntu/Debian 
```
sudo dpkg -i hydra-reverse-proxy-0-1.x86_64.deb
sudo apt-get install -f
```
### CentOS/RedHat/Fedora
```
sudo yum install hydra-reverse-proxy-0-1.x86_64.rpm
```

### Run
Edit the configuration file /etc/hydra-reverse-proxy.conf. After that just run:
```
sudo /etc/init.d/hydra-reverse-proxy start
```

## Configuration

Configuration options can be set in two places:

 1. Command line flags
 2. Configuration file

Options set on the command line take precedence over all other sources.

### Command Line Flags

* `-app-id` - (Required) The registered id in hydra for the application which wishes to reach through the proxy (i.e "app001").
* `-hydra-servers` - (Required) URLs list of known hydra servers (i.e "http://215.10.111.201:7772,http://215.10.111.202:7772")
* `-proxy-addr` - (Required) The address where the proxy remains listening (i.e ":3000").
* `-apps-cache-duration` - The number of milliseconds to wait before attempting to update the application cache. Default to 20000 milliseconds.
* `-duration-between-servers-retries` - The number of milliseconds in between retries for the set of servers. Default to 0 milliseconds.
* `-hydra-servers-cache-duration` - The number of milliseconds to wait before attempting to update the hydra cache. Default to 60000 milliseconds.
* `-max-number-of-retries` - Maximum number of retries per server. Default to 10 milliseconds.

## Configuration File

The hydra reverse proxy configuration file is written in [TOML](https://github.com/mojombo/toml)
and reads from `/etc/hydra-reverse-proxy.conf` by default.

```TOML
app_id = "app001"
hydra_servers = ["http://215.10.111.201:7772","http://215.10.111.202:7772"]
proxy_addr = ":3000"
[HydraClient]
apps_cache_duration = 20000
duration_between_all_servers_retry = 0
hydra_servers_cache_duration = 60000
max_number_of_retries = 10
```
