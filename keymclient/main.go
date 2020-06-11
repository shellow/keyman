package main

import (
	"go.uber.org/zap"
	"github.com/urfave/cli"
)

var Logger *zap.SugaredLogger

var HOSTURL string

func main() {
	var language string

	app := cli.NewApp()
	app.Name = "Key manage"
	app.Usage = "hello world"
	app.Version = "1.2.3"

}

