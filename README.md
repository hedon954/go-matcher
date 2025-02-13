# go-matcher

[![Go Report Card](https://goreportcard.com/badge/github.com/hedon954/go-matcher)](https://goreportcard.com/report/github.com/hedon954/go-matcher)
[![codecov](https://codecov.io/github/hedon954/go-matcher/graph/badge.svg?token=FEW1EL1FKG)](https://codecov.io/github/hedon954/go-matcher)
[![CI](https://github.com/hedon954/go-matcher/workflows/build/badge.svg)](https://github.com/hedon954/go-matcher/actions)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/hedon954/go-matcher?sort=semver)](https://github.com/hedon954/go-matcher/releases)

Go-matcher is a game matcher implement in Go, which supports add game mode and match strategy easily.

- `GameMode`: The identifier of each different game, used to define as enum.
- `MatchStrategy`: The strategy to match players, used to define as interface.

## Features

- [ ] API
  - [x] HTTP
  - [x] TCP
  - [ ] UDP
  - [ ] KCP
  - [ ] WebSocket
  - [ ] gRPC
- [ ] Service
  - [x] match service
  - [ ] push service
  - [ ] settle service
- [x] Swagger Doc
- [ ] timer
  - [x] native timer
  - [x] asynq timer
  - [ ] redis timer
- [x] GameMode
  - [x] GoatGame
- [ ] MatchStrategy
  - [x] Glicko2
  - [ ] Gather
  - [ ] ELO
- [x] Config
  - [x] File Loader
  - [x] Nacos Dynamic Loader
- [ ] AI Generator
- [x] Open Telemetry
  - [x] Logger
  - [x] Tracer
  - [x] Metric
- [x] Graceful restart
  - [x] tableflip
  - [x] save and reload entries


## Start With Open Telemetry

```bash
make docker-up
```

## Architecture

![go-matcher architecture](assets/img/architecture.png)

## Smoke flow

![smoke flow](./assets/uml/smoke_flow.svg)

## State Machine

![state machine](./assets/img/delay-timer.png)

## Entry Design

![entry design](./assets/img/entry.png)

## Todo

- [ ] network
  - [ ] UDP
  - [ ] KCP
  - [ ] WebSocket
  - [ ] gRPC
- [ ] dynamic config
  - [ ] etcd
  - [ ] apollo
  - [ ] consul
- [ ] tracer
- [ ] repository stats
- [ ] match queue stats
- [ ] graceful restart
- [ ] repository by redis
- [ ] hot upgrade
- [ ] horizontal expansion
- [ ] AI Generator

## Help

### Before writing code

```bash
bash ./setup_pre_commit.sh
```

### Start asynq web UI `Asynqmon`

```bash
# clone repo
git clone git@github.com:hibiken/asynqmon.git
# enter repo
cd asynqmon
# check node version
node -v
# if higher 16, should degrade it
# install nvm if not
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
# install node 16
nvm install 16
# use node 16
nvm use 16
# install yarn if not
brew install yarn
# build asynqmon
make build
# start asynqmon
./asynqmon
```

### How to use Open Telemetry

[Open Telemetry](https://opentelemetry.io/docs/languages/go/getting-started/)

### Start jaeger

```bash
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest
```

### Start Nacos

run by docker compose:

```bash
git clone https://github.com/nacos-group/nacos-docker.git \
cd nacos-docker \
docker-compose -f example/standalone-mysql-8.yaml up
```

open nacos website:

```bash
open http://localhost:8848/nacos/index.html
```
