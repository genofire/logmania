# logmania

[![pipeline status](https://dev.sum7.eu/sum7/logmania/badges/master/pipeline.svg)](https://dev.sum7.eu/genofire/logmania/pipelines)
[![coverage report](https://dev.sum7.eu/sum7/logmania/badges/master/coverage.svg)](https://dev.sum7.eu/genofire/logmania/pipelines)
[![Go Report Card](https://goreportcard.com/badge/dev.sum7.eu/sum7/logmania)](https://goreportcard.com/report/dev.sum7.eu/genofire/logmania)
[![GoDoc](https://godoc.org/dev.sum7.eu/sum7/logmania?status.svg)](https://godoc.org/dev.sum7.eu/genofire/logmania)


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

## Get logmania

#### Download

Latest Build binary from ci here:

[Download All](https://dev.sum7.eu/sum7/logmania/-/jobs/artifacts/master/download/?job=build-my-project) (with config example)

[Download Binary](https://dev.sum7.eu/sum7/logmania/-/jobs/artifacts/master/raw/bin/logmania?inline=false&job=build-my-project)

#### Build

```bash
go get -u dev.sum7.eu/sum7/logmania
```

## Configure

see `config_example.toml`


## Related Projects

- [hook2xmpp](https://dev.sum7.eu/sum7/hook2xmpp) for e.g. grafana, alertmanager(prometheus), gitlab, git, circleci
