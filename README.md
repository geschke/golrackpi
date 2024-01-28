# golrackpi is a Go (Golang) Library Rest Api Client (for) Kostal Plenticore Inverters (with CLI)

[![Go Reference](https://pkg.go.dev/badge/github.com/geschke/golrackpi.svg)](https://pkg.go.dev/github.com/geschke/golrackpi)

This repository provides a Go (Golang) library for the undocumented REST-API of Kostal Plenticore Inverters. It uses the PIKO IQ / PLENTICORE plus API with its swagger documentation found at "inverter ip address"/api/v1/.

This library is not affiliated with Kostal and is no offical product of KOSTAL Solar Electric GmbH or any subsidiary company of Kostal Gruppe.

## Features

* Authenticate (Login, Logout, Check authentication)
* Read/Write settings
* Read processdata
* Read events

Additional:

* Commandline interface (CLI) to get any kind of returned inverter data

## Getting Started

Please be patient. The library is in development and currently in the last stage before a release will be published. So I'm cleaning up a bit and try to add some senseful comments and documentation in the next step.

...todo...

### Installing the library

...todo...

### Using the command line interface




```shell

 golrackpi is a small CLI application to read values from Kostal Plenticore Inverters.

Usage:
  golrackpi [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  events      Get the latest events
  help        Help about any command
  info        Returns miscellaneous information
  modules     List modules content
  processdata List processdata values
  settings    List settings content

Flags:
  -h, --help              help for golrackpi
  -p, --password string   Password (required)
  -m, --scheme string     Scheme (http or https, default http)
  -s, --server string     Server (e.g. inverter IP address) (required)

Use "golrackpi [command] --help" for more information about a command.


```

 
### Using the library from Go

...todo...


## Documentation

...todo...
## Write settings

Available settings can be found in the swagger documentation of the inverter or by calling `client.Settings()``. The following example shows how to activate smart battery control:

```go
  module := golrackpi.ModuleSettings{ModuleId: "devices:local"}
  module.Settings = []golrackpi.SettingsValues{golrackpi.SettingsValues{Id: "Battery:SmartBatteryControl:Enable", Value: "0"}}
  client.UpdateSettings([]golrackpi.ModuleSettings{module})
```

## License

MIT

## Thanks to

* [kilianknoll](https://github.com/kilianknoll) for the kostal-RESTAPI project 
* [stegm](https://github.com/stegm) for the pykoplenti Python REST client API project
* Marco Tröster ([Bonifatius94](https://github.com/Bonifatius94)) in the issue of openhab-addons for some Java code provided in https://github.com/openhab/openhab-addons/issues/7492
