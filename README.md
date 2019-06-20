# logmania
[![DroneCI](https://ci.sum7.eu/api/badges/genofire/logmania/status.svg?branch=master)](https://ci.sum7.eu/genofire/logmania)
[![CircleCI](https://circleci.com/gh/genofire/logmania/tree/master.svg?style=shield)](https://circleci.com/gh/genofire/logmania/tree/master)
[![Coverage Status](https://coveralls.io/repos/github/genofire/logmania/badge.svg?branch=master)](https://coveralls.io/github/genofire/logmania?branch=master)
[![Go Report Card](https://goreportcard.com/badge/dev.sum7.eu/genofire/logmania)](https://goreportcard.com/report/dev.sum7.eu/genofire/logmania)
[![GoDoc](https://godoc.org/dev.sum7.eu/genofire/logmania?status.svg)](https://godoc.org/dev.sum7.eu/genofire/logmania)


This is a little logging server.

## input
It receive logs (events) by:
- syslog
- journald (with service nc)

## output
And forward this logs (events) to multiple different output:
- xmpp (client and muc)
- file

there a multi filter possible
- regex
- priority

it could replace text by regex expression

configuration live possible by bot (on input e.g. xmpp)

## Related Projects

- [hook2xmpp](https://dev.sum7.eu/genofire/hook2xmpp) for e.g. grafana, alertmanager(prometheus), gitlab, git, circleci
