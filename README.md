# nbminer-reporter
A simple app that reads NBMiner status REST API data and sends it to InfluxDB

```bash
docker run --rm -it --name nbreporter \
-v `pwd`:/src -w /src golang:1.17.3-alpine3.14 sh
```
