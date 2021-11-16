# NBMiner Reporter
A simple Go app that reads NBMiner status REST API data and sends it to InfluxDB.

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/Cruzoft/nbminer-reporter/cicd)
![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/Cruzoft/nbminer-reporter/main?label=Go%20Version)
![GitHub](https://img.shields.io/github/license/Cruzoft/nbminer-reporter?color=orange)

## Usage

Using the reporter is quite easy, specially if you've already setup NBMiner before.

### Before you start

This app assumes your NBMiner is exposing status data through its REST API. If you haven't enable it, follow these simple steps:

1. Go NBMiner folder.
1. Open and edit the file you execute to start the miner (the one that has your wallet id in it).
1. At the end of the line, add the following flag `--api 127.0.0.1`.

    What this flag does is enabling a REST API that exposes the miner status data, such as temp or hasrates.

    You can check all the info at NBMiner official docs: [here](https://github.com/NebuTech/NBMiner#cmd-options).



1. Save and close the file, and now restart the miner.
1. Let's check the API is working. Open a web browser, and go to http://localhost:8000.

### On Windows

1. Get the **Windows** files from the latests release: [Download Now](https://github.com/Cruzoft/nbminer-reporter/releases)
1. Unzip the file on a separated folder than NBMiner.
1. Edit the `start_win_nbreporter.bat` file:
    1. Change `rigX` with the friendly name you want your miner to be called
    2. Change `influxdb.my.organization.net` URL with your InfluxDB server URL.
    3. If you need to set a different InfluxDB port than default, use the option `-p 8080`.
    4. If you need to set a different InfluxDB schema than default, use the option `-l https`.
1. Execute `start_win_nbreporter.bat`. A shell window will open, and you should see an output like the following:

    ```shell
    INFO[0000] NBMiner Status Reporter Initiated            
    INFO[0000] Using Friendly Name: rigX             
    INFO[0005] Checking Status.                   
    ```

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
