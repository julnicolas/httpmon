package main

import (
	"os"
	"time"

	"github.com/julnicolas/httpmon/pkg/app"
	"github.com/julnicolas/httpmon/pkg/config"
	"github.com/sirupsen/logrus"
)

func exitErr(err error) {
	logrus.Errorln(err)
	os.Exit(1)
}

func main() {
	conf, err := config.CLI(config.Default())
	if err != nil {
		exitErr(err)
	}

	if conf.Debug {
		// Wait a few seconds before starting
		// to connect a debugger
		time.Sleep(5 * time.Second)
	}

	app := app.NewApp(conf)

	if err := app.Init(); err != nil {
		exitErr(err)
	}

	defer app.Close()
	if err := app.Run(); err != nil {
		exitErr(err)
	}
}
