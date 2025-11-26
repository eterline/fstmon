
# fstmon

Monitoing microservice for [homepage](https://gethomepage.dev) or another. Written in ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

## In Action

![screen_1](./screenshots/1.png)
![screen_2](./screenshots/2.png)

## Build

To build project in full prod variant
```sh
task build-prod
```

for testing build
```sh
task build
```

## Arguments

### Logging

| Option                  | Description                                     | Default |
| ----------------------- | ----------------------------------------------- | ------- |
| `--log-level LOG-LEVEL` | Logging level: `debug` `info`, `warn`, `error`  | `info`  |
| `--log-json`, `-j`      | Set logs to JSON format                         | `false` |

### Server

| Option                               | Description                          | Default     |
| ------------------------------------ | ------------------------------------ | ----------- |
| `--listen LISTEN`, `-l LISTEN`       | Server listen address                | `:3000`     |
| `--certfile CERTFILE`, `-c CERTFILE` | Server SSL certificate file          | —           |
| `--keyfile KEYFILE`, `-k KEYFILE`    | Server SSL key file                  | —           |
| `--sni SNI`                          | Server allowed request hosts         | `[]`        |
| `--subnets SUBNETS`, `-s SUBNETS`    | Server allowed source subnets/IPs    | `[]`        |
| `--token TOKEN`, `-t TOKEN`          | Server auth token string             | env `TOKEN` |
| `--ip-header`                        | Enable parsing reverse proxy headers | `false`     |

### Metric loops

| Option                              | Description                      | Default |
| ----------------------------------- | -------------------------------- | ------- |
| `--cpu-loop CPU-LOOP`               | CPU metric update loop (seconds) | `5`     |
| `--avgload-loop AVGLOAD-LOOP`       | Avgload update loop (seconds)    | `10`    |
| `--system-loop SYSTEM-LOOP`         | System update loop (seconds)     | `30`    |
| `--network-loop NETWORK-LOOP`       | Network update loop (seconds)    | `5`     |
| `--partitions-loop PARTITIONS-LOOP` | Partitions update loop (seconds) | `30`    |

### Help

| Option         | Description                |
| -------------- | -------------------------- |
| `--help`, `-h` | Display this help and exit |


## Running

Homepage config:
```yml
// There will be your service config
          - type: customapi
            url: http://<host>:3300/monitoring/system
            method: GET
            refreshInterval: 10000
            mappings:
             - field: data.ram
               label: ram
               format: text
             - field: data.uptime
               label: uptime
               format: text
```

Another usage will be describe later

## License

[MIT](https://choosealicense.com/licenses/mit/)