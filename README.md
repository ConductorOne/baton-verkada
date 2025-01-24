# `baton-verkada` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-verkada.svg)](https://pkg.go.dev/github.com/conductorone/baton-verkada) ![main ci](https://github.com/conductorone/baton-verkada/actions/workflows/main.yaml/badge.svg)

`baton-verkada` is a connector for Verkada built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the Verkada API to sync data about access users and access groups in your Verkada organization.
Check out [Baton](https://github.com/conductorone/baton) to learn more about the project in general.

# Getting Started

## Prerequisites

- Access to the Verkada Command platform.
- API key. Only Organization Admins can create API keys. To create a new API key, go to Organization settings -> Verkada API -> New API Key
- Set the API key permissions to 'read only' if your want to just fetch the data about users and access. If you want provissioning support (adding and removing users from access groups) set the permissions to 'read/write'. 

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-verkada

BATON_API_KEY=verkadaApiKey baton-verkada
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_API_KEY=verkadaApiKey ghcr.io/conductorone/baton-verkada:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-verkada/cmd/baton-verkada@main

BATON_API_KEY=verkadaApiKey baton-verkada
baton resources
```

# Data Model

`baton-verkada` will pull down information about the following Verkada resources:

- Access Users
- Access Groups

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-verkada` Command Line Usage

```
baton-verkada

Usage:
  baton-verkada [flags]
  baton-verkada [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --api-key string         API key used to authenticate to Verkada API. ($BATON_API_KEY)
      --client-id string       The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string   The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string            The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                   help for baton-verkada
      --log-format string      The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string       The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -p, --provisioning           This must be set in order for provisioning actions to be enabled. ($BATON_PROVISIONING)
      --region string          API region. Default is US. In case of EU based organization, pass region as EU. ($BATON_REGION) (default "US")
  -v, --version                version for baton-verkada

Use "baton-verkada [command] --help" for more information about a command.
```
