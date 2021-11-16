# NBMiner Reporter
A simple Go app that reads [NBMiner](https://github.com/NebuTech/NBMiner#nbminer) status REST API data and sends it to InfluxDB.

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

    You can check all the info at NBMiner official docs: [here](https://github.com/NebuTech/NBMiner#api-reference).

    **IMPORTANT**: This API does exposes your wallet id, but is still only accessible from the machine. NBMiner Reporter WON'T send the wallet id to InfluxDB, that's the only field ignored while sending the data.

1. Save and close the file, and now restart the miner.
1. Let's check the API is working. Open a web browser, and go to http://localhost:8000.

    You should see a website with a bunch of information about your miner, the same that is shown in the terminal while the miner is running, plus some extra data.

That's it, you're ready to continue with NBMiner Reporter installation.

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

That's it, you're now sending your miner status to InfluxDB.

### On Linux

1. Get the **Linux** files from the latests release: [Download Now](https://github.com/Cruzoft/nbminer-reporter/releases)
1. Untar the file on a separated folder than NBMiner.
1. Edit the `start_lnx_nbreporter.sh` file:
    1. Change `rigX` with the friendly name you want your miner to be called
    2. Change `influxdb.my.organization.net` URL with your InfluxDB server URL.
    3. If you need to set a different InfluxDB port than default, use the option `-p 8080`.
    4. If you need to set a different InfluxDB schema than default, use the option `-l https`.
1. Execute `start_lnx_nbreporter.sh`. A terminal window will open, and you should see an output like the following:

    ```shell
    INFO[0000] NBMiner Status Reporter Initiated            
    INFO[0000] Using Friendly Name: rigX             
    INFO[0005] Checking Status.                   
    ```

That's it, you're now sending your miner status to InfluxDB.

## CMD Options

Customize the way NBMiner Reporter works by using the following options:

```bash
nbreporter [-v] [-b string] [-f number] [--help] [-h string] [-l string] [-n string] [-o string] [-p number] [-r strinumberng] [-s string] [-t string]
```

Check the options details.

| Short Flag | Long Flag | Description                                     |
|----|-------------------|-------------------------------------------------|
| -n | --name=string     | A friendly name for miner. Default: hostname    |
| -f | --freq=number     | Status check frequency in seconds. Default: 60  |
| -h | --ihost=string    | InfluxDB Host.  Default: localhost              |
| -p | --iport=number    | InfluxDB Port. Default: 8086                    |
| -t | --itoken=string   | InfluxDB Access Token.                          |
| -l | --iproto=string   | InfluxDB Protocol.  Default: http               |
| -b | --ibucket=string  | InfluxDB Bucket. Default: miner                 |
| -o | --iorg=string     | InfluxDB Organization.  Default: miner-org      |
| -r | --nbport=number   | NBMiner API Port. Default: 8000                 |
| -s | --nbhost=string   | NBMiner API Host. Default: localhost            |
| -v |                   | Run in Verbose mode. Default: false             |
|    | --help            | Show usage options.                             |

## Compatibility

NBMiner Reporter has been tested using the following setups:

| NBM Reporter Ver. | NBMiner Vers. | OS                 | InfluxDB             |
|-------------------|---------------|--------------------|----------------------|
| v1.0.X            | 39.7          | Windows 10, HiveOS | 1.8.10, 2.0.x, 2.1.x |

## Contribute

### Local Build

To work on your local machine, and build the go app, we recommend using a docker container, so you don't have to install GO SDK if you don't have it, and ensure you have a clean environment to work with.

1. Fork this repo, and clone it in your computer.
1. Open a terminal, and move to the repository folder.
1. Use this command to start the container:

    ```sh
    docker run --rm -it --name nbreporter -v `pwd`:/src -w /src golang:1.17.3-alpine3.14 sh
    ```

1. Once on the container, build the app by running  `go build ./cmd/nbreporter/...`

### Developing locally

Now you have your app built, you need to be able to run it. This app takes data from a REST API, and sends it to an InfluxDB service, so in order to be able to test it, you'll need both services accessibe from the development machine.

Worry not, you won't need to access an actual miner rig, we got you covered. We've created a NBMiner Status Simulator, basically a very simple web app that exposes the same REST API as the NBMiner, but with random numbers. Enough for develpment and testing right? Also, we've set it up on a docker-compose file with an InfluxDB service to use.

1. Build the simulator image locally:

    ```sh
    docker-compose build
    ```

1. Start the services:

    ```sh
    docker-compose up
    ```

1. Back on your development container, run the app:

    ```sh
    go run ./cmd/nbreporter/... -s host.docker.internal -t shhh-secret-token -f 5 -h host.docker.internal
    ```

1. Go to your browser and point it to http://localhost:8888. You should see a Cronograf UI wich you can use to query InfluxDB and check the data points.