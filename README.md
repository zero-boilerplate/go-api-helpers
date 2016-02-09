# go-api-generator
Common boilerplate code used in golang API (backend) apps

## Goals

The main of this repo is to reduce the amount of boiler-plate code. Using builder classes for elegance.

So essentially this repo just abstracts as much as possible boilerplate into this repo and expose it via a simple api.

## Features

Easily create a installable service with this code:
```
package main

import (
    "github.com/zero-boilerplate/go-api-helpers/service"
    "time"

    service2 "github.com/ayufan/golang-kardianos-service"
)

type runHandler struct{}

func (r *runHandler) Run(logger service2.Logger) {
    for {
        dur := 10 * time.Second
        logger.Infof("I am a cool new app running, now sleeping for %s", dur.String())
        time.Sleep(dur)
    }
}

func main() {
    r := &runHandler{}
    service.NewServiceRunnerBuilder("TestSleepService", r).Run()
}
```

To install as a service just build your app and call `BUILT_BINARY -service install`. Thanks to **https://github.com/ayufan/golang-kardianos-service** for making this so simple.
