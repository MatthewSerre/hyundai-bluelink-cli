# Hyundai Bluelink CLI

A repository for a CLI application that uses client-side code generated from the `hyundai-bluelink-protobufs` repository [link](https://github.com/MatthewSerre/hyundai-bluelink-protobufs) to communicate with various services for authentication, vehicle information, and remote vehicle actions via Hyundai Bluelink.

## Getting Started

* [Install Go](https://go.dev/doc/install)

## Usage

Run `go run main.go` to run the bare CLI command, which will generate a `.env` with fields for credential values. Then run `go run main.go -h` to see the available commands. In a local environment, the CLI needs to connect to the services hosted in these repositories:

- `hyundai-bluelink-authentication-server` [link](https://github.com/MatthewSerre/hyundai-bluelink-authentication-server)
- `hyundai-bluelink-vehicle-information-service` [link](https://github.com/MatthewSerre/hyundai-bluelink-vehicle-information-service)
- `hyundai-bluelink-remote-action-service` [link](https://github.com/MatthewSerre/hyundai-bluelink-remote-action-service)

Follow the instructions in each repository's README to set up the services to work with the CLI.

## Contributing

Create an issue and/or a pull request and I will take a look.

***

This project is not affiliated with Hyundai in any way.