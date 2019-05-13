# Slaxy

Sentry webhooks to **sla**ck message converter pro**xy**.

[![Build Status](https://travis-ci.org/innogames/slaxy.svg)](https://travis-ci.org/innogames/slaxy)
[![GoDoc](https://godoc.org/github.com/innogames/slaxy?status.svg)](https://godoc.org/github.com/innogames/slaxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/innogames/slaxy)](https://goreportcard.com/report/github.com/innogames/slaxy)
[![Release](https://img.shields.io/github/release/innogames/slaxy.svg)](https://github.com/innogames/slaxy/releases)

## Contents

- [General](#general)
- [Installation](#installation)
- [Usage](#usage)
  - [Example Config](#example-config)
  - [CLI](#cli)

## General

Slaxy provides an HTTP web interface that is capable of receiving Sentry webhooks, converting them to Slack messages and posting them to Slack on behalf of a bot user for example.
A token is required to authenticate to Slack.

Once the server is up and running it will continuously receive webhooks from Sentry and post them to the configured channel in your Slack workspace.

## Installation

You can simply go install the binary that serves as a HTTP server:
```
$ go install github.com/innogames/slaxy/slaxy
```

## Usage

### Example Config

```
grace-period: 60s
addr: 127.0.0.1:3000
bot-token: xoxb-###-###-###
channel: xxx
excluded-fields:
  - ^sentry:.*$
```

### CLI

```
  _________.__
 /   _____/|  | _____  ___  ______.__.
 \_____  \ |  | \__  \ \  \/  <   |  |
 /        \|  |__/ __ \_>    < \___  |
/_______  /|____(____  /__/\_ \/ ____|
        \/           \/      \/\/

Usage:
  slaxy [flags]

Flags:
  -a, --addr string               listen address (default "localhost:3000")
  -n, --channel string            slack channel
  -c, --config string             path to config file if any
  -e, --excluded-fields strings   excluded sentry fields
  -g, --grace-period duration     grace period for stopping the server (default 1m0s)
  -h, --help                      help for slaxy
  -t, --token string              slack token
```
