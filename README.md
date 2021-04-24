# dislog

[![Go Reference](https://pkg.go.dev/badge/github.com/DisgoOrg/dislog.svg)](https://pkg.go.dev/github.com/DisgoOrg/dislog)
[![Go Report](http://goreportcard.com/badge/DisgoOrg/dislog)](https://goreportcard.com/report/DisgoOrg/dislog)
[![Go Version](https://img.shields.io/github/go-mod/go-version/DisgoOrg/dislog)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/DisgoOrg/dislog/blob/master/LICENSE)
[![Disgo Version](https://img.shields.io/github/v/release/DisgoOrg/dislog)](https://github.com/DisgoOrg/dislog/releases/latest)
[![Disgo Discord](https://img.shields.io/badge/Disgo%20Discord-blue.svg)](https://discord.gg/mgjJeufk)

dislog is a [logrus](https://github.com/sirupsen/logrus) [logging hook](https://github.com/sirupsen/logrus#hooks) sending logs over [Discord Webhooks](https://discord.com/developers/docs/resources/webhook) using the [disgohook](https://github.com/DisgoOrg/dislog) library


## Getting Started

### Installing

`go get`

```sh
go get github.com/DisgoOrg/dislog
```

### Usage

Import the package into your project.

```go
import "github.com/DisgoOrg/dislog"
```

Create a new Webhook by `webhook_id/webhook_token` and pass a [logger](https://github.com/DisgoOrg/log) like [logrus](https://github.com/sirupsen/logrus). This webhook then can be used to send/edit/delete messages

```go
logger := logrus.New()
dlog, err := dislog.NewDisLogByToken(nil, logrus.InfoLevel, os.Getenv("webhook_token"), dislog.TraceLevelAndAbove...)
if err != nil {
logger.Errorf("error initializing dislog %s", err)
return
}
defer dlog.Close()
logger.AddHook(dlog)
```

## Documentation

Documentation is unfinished and can be found under

* [![Go Reference](https://pkg.go.dev/badge/github.com/DisgoOrg/dislog.svg)](https://pkg.go.dev/github.com/DisgoOrg/dislog)
* [![logrus Hooks Documentation](https://img.shields.io/badge/logrus%20Documentation-blue.svg)](https://github.com/sirupsen/logrus#hooks)

## Examples

You can find and example under [example](https://github.com/DisgoOrg/dislog/tree/master/example)

## Troubleshooting

For help feel free to open an issues or reach out on [Discord](https://discord.gg/mgjJeufk)

## Contributing

Contributions are welcomed but for bigger changes please first reach out via [Discord](https://discord.gg/mgjJeufk) or create an issue to discuss your intentions and ideas.

## License

Distributed under the [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/DisgoOrg/dislog/blob/master/LICENSE). See LICENSE for more information.