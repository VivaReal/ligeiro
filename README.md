# Ligeiro

![ligeiro](https://user-images.githubusercontent.com/379894/33691641-4a943e4e-dad0-11e7-846a-540dc7923d86.png)

Minimalist tools for your golang applications.

## Motivation

After copy/pasting very often small pieces of code that rely on some conventions such as how to handle environment configurations, format logs according with company requirements a few times, it got painful to change everywhere when it need to change or have a new requirement.

This library attempts to provide minimalist tools that you help your application to comply with some platform requirements such as configs, logs, etc.

## Usage

### Configs

Basically pass a map of env vars with their respective defaults, it will be merged with bundled
configs (ENVIRONMENT, LOG_LEVEL and VERSION):

```go
import (
	"fmt"

	"github.com/vivareal/ligeiro/envcfg"
)

func main() {
	Config := envcfg.Load(envcfg.EnvCfgMap{
		"APPLICATION": "myappname",
		"LISTEN_PORT": ":8080",
		"MAX_BYTES": "3000000"
	})

	fmt.Println(Config.Get("listenPort"))
	fmt.Println(Config.Get("application"))
	fmt.Println(Config.Get("logLevel"))
	fmt.Println(Config.GetInt("maxBytes"), fmt.Sprintf("%T", Config.GetInt("maxBytes")))
}
```

Running previous file with LOG_LEVEL env var set:

```
$ LOG_LEVEL=error go run main.go

:8080
myappname
error
3000000 int
```

### Logs

Offers a thin wrapper to github.com/sirupsen/logrus to log GELF like format to stdout, ready to docker GELF drivers. The
best way to use it to define your own `applog` package with custom fields registered:

```go
package myapi

import "github.com/vivareal/ligeiro/logger"

var Applog = logger.WithFields(logger.Fields{
	"application": "myapi",
	"squad":       "growth",
})

func example() {
	Applog.Info("ligeiro")
	// Output: {"application":"myapi","environment":"dev","fields.level":6,"full_message":"ligeiro","level":6,"level_name":"info","time":"2017-12-07T17:25:13-02:00","timestamp":1512674713370,"version":"detached"}
}
```

Note: _`logger` already checks `envcfg` to add some custom fields to formatted log output_
