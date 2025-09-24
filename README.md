
# fstmon

Monitoing microservice for [homepage](https://gethomepage.dev) or another. Written in ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

## In Action

![screen_1](./screenshots/1.png)
![screen_2](./screenshots/2.png)

## Build

To build project in full prod variant
```
make build-prod
```

for testing build
```
make build
```

All binariaes will be in './build' folder at repository


## Usage:

```bash
Usage: fstmon [--debug] [--log-json] [--listen LISTEN] [--certfile CERTFILE] [--keyfile KEYFILE] [--sni SNI] [--subnets SUBNETS] [--token TOKEN] [--ip-header]

Options:
  --debug                Allow debug logging level
  --log-json             Set logs to JSON format
  --listen LISTEN        Server listen address [default: :8100]
  --certfile CERTFILE    Server SSL certificate file [env: CERT]
  --keyfile KEYFILE      Server SSL key file [env: KEY]
  --sni SNI              Server allowed request hosts [default: []]
  --subnets SUBNETS      Server allowed source subnets/IPs [default: []]
  --token TOKEN          Server auth token string
  --ip-header            Enable parsing reverse proxy headers
  --help, -h             display this help and exit
```

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