# Miner Status REST API Simulator

## Develop

```bash
docker run --rm -it --name simulator -p 8000:8000 -v `pwd`:/src -w /src golang:1.17.3-alpine3.14 sh
```

## Build

```bash
docker build -t nbminer-simulator .
```

## Run

```bash
docker run --rm --name nbsimulator -p 8000:8000 nbminer-simulator
```
