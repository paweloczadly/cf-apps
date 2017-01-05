# cf-apps

Displays IP addresses and ports of Cloud Foundry applications in each space.

## Running with Docker

```bash
$ docker run -ti paweloczadly/cf-apps https://cloudfoundry.url "$(cf oauth-token)"
```

## Usage

```bash
$ go install
$ cf-apps https://cloudfoundry.url "$(cf oauth-token)"
```

## Building Docker image

```bash
$ go build
$ curl -o ca-certificates.crt https://raw.githubusercontent.com/bagder/ca-bundle/master/ca-bundle.crt
$ docker build .
```
