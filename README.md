# nbminer-reporter
A simple app that reads NBMiner status REST API data and sends it to InfluxDB

Get a local influx and nbminer simulator

```bash
docker-compose up
```

Run the dev container

```bash
docker run --rm -it --name nbreporter \
-v `pwd`:/src -w /src golang:1.17.3-alpine3.14 sh
```

Run the go app inside the container

```bash
go run ./cmd/nbminer-reporter/... -n rig03-sim -s host.docker.internal -t shhh-secret-token -f 5 -h host.docker.internal
```

GOOS=windows GOARCH=amd64 go build -o bin/app-amd64.exe app.go
nbreporter -n rig03 -l https -h influxdb-miner.orion.net.ar -p 443
