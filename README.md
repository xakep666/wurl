# wurl - console client for websocket protocol

[![Documentation](https://godoc.org/github.com/github.com/xakep666/wurl?status.svg)](http://godoc.org/github.com/xakep666/wurl)
[![Go Report Card](https://goreportcard.com/badge/github.com/xakep666/wurl)](https://goreportcard.com/report/github.com/xakep666/wurl)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/github.com/xakep666/wurl/LICENSE)

## Abstract

At the moment we already have quite few websocket clients. Starting from browser addons, ending with console clients.

But I`m not satisfied with either of them. Browser addons requires installed and running browser.
NodeJS-based clients requires node and tons of dependencies.
But most importantly, none of them allows you to specify additional headers for request.

So I decided to write own console websocket client...

## Installation
`go get -u github.com/xakep666/wurl`

## Current features
- Read text/binary messages from connection and display it
- Ability to set additional headers for connection upgrade request
- Correctly processes ping message (by default responses with pongs message)
- Can periodically send ping message to server (period can be set through flags)

### TODOs for v1
- [ ] Document all packages
- [ ] Flag to show handshake response
- [ ] Store and load options from file
- [ ] Warning about binary messages before displaying (cURL-like)
- [ ] Ability to specify output for binary and text messages
- [ ] Option to send message to server before reading
- [ ] Good description for all flags/commands