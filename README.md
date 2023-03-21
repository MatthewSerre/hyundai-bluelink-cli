# Hyundai Bluelink CLI

A repository for a CLI application that uses client-side code generated from the `hyundai-bluelink-protobufs` repository [link](https://github.com/MatthewSerre/hyundai-bluelink-protobufs) to communicate with various services for authentication, vehicle information, and remote vehicle actions via Hyundai Bluelink.

## Getting Started

* [Install Go](https://go.dev/doc/install)

## Usage

Use the existing binary to run the bare CLI command

```
bin/cmd/hb
``` 

or make changes to the code, build a new binary, and run the bare CLI command

```
task build
bin/cmd/hb
```

Doing so will generate a `.env` with fields for credential values (enter the username, password, and PIN for the Hyundai account associated with the Bluelink service). Then run `bin/cmd/hb -h` to see the available commands. In a local environment, the CLI needs to connect to the services hosted in these repositories:

- `hyundai-bluelink-authentication-server` [link](https://github.com/MatthewSerre/hyundai-bluelink-authentication-server)
- `hyundai-bluelink-vehicle-information-service` [link](https://github.com/MatthewSerre/hyundai-bluelink-vehicle-information-service)
- `hyundai-bluelink-remote-action-service` [link](https://github.com/MatthewSerre/hyundai-bluelink-remote-action-service)

Follow the instructions in each repository's README to set up the services to work with the CLI.

## Contributing

Create an issue and/or a pull request and I will take a look.

***

This project is not affiliated with Hyundai in any way. Credit to [TaiPhamD](https://github.com/TaiPhamD) and his `bluelink_go` project [link](https://github.com/TaiPhamD/bluelink_go) for inspiration and some code snippets.