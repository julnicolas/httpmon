package config

import (
	"flag"
	"fmt"
	"time"
)

type Config struct {
	Debug          bool // Debug flag, waits a few seconds before starting
	Period         time.Duration
	File           string // file name or "stdin"
	ReadBufferSize uint
	Alert          Alert
}

// Default creates a new default configuration structure
func Default() Config {
	return Config{
		Period:         10 * time.Second,
		ReadBufferSize: 100,
		Alert:          Alert{}.Default(),
	}
}

// CLI parses the cli, overriding matching existing config fields
func CLI(conf Config) (Config, error) {
	cli := cliInput{}
	flag.BoolVar(&cli.debug, "debug", false, "wait a few seconds before starting")
	flag.StringVar(&cli.file, "file", conf.File, "csv file to read http traces from")
	flag.BoolVar(&cli.stdin, "stdin", false, "read http logs from stdin, takes precendence over --file")
	flag.DurationVar(&cli.period, "period", conf.Period, "log aggregation period used to generate metrics values (go duration format)")
	flag.UintVar(&cli.bufferLen, "lines", conf.ReadBufferSize, "size of the line buffer when reading logs")
	flag.DurationVar(&cli.alertDuration, "alert-duration", conf.Alert.RequestsPerSecond.Period, "if requests/s > --threshold for --alert-duration then the alert is active (go duration format)")
	flag.UintVar(&cli.alertThreshold, "alert-threshold", uint(conf.Alert.RequestsPerSecond.Threshold), "requests/s threshold over wich the alert becomes active")
	flag.Parse()

	if err := fromCLI(&conf, cli); err != nil {
		return conf, err
	}

	return conf, nil
}

type cliInput struct {
	debug          bool
	file           string
	stdin          bool
	period         time.Duration
	bufferLen      uint
	alertDuration  time.Duration
	alertThreshold uint
}

func cliValidation(cli cliInput) error {
	if !cli.stdin {
		if cli.file == "" {
			return fmt.Errorf("--file - name missing")
		}
	}

	if cli.period < time.Second {
		return fmt.Errorf("--period - minimum period is 1s, received %s", cli.period)
	}

	if cli.bufferLen == 0 {
		return fmt.Errorf("--lines - buffer length must be greater than 0")
	}

	if cli.alertDuration < time.Second {
		return fmt.Errorf("--alert-duration - minimum period is 1s, received %s", cli.period)
	}

	return nil
}

func fromCLI(conf *Config, cli cliInput) error {
	if conf == nil {
		return fmt.Errorf("fromCli - nil configuration")
	}

	if err := cliValidation(cli); err != nil {
		return err
	}

	conf.Debug = cli.debug
	if cli.stdin {
		conf.File = "stdin"
	} else {
		conf.File = cli.file
	}

	conf.Period = cli.period
	conf.ReadBufferSize = cli.bufferLen
	conf.Alert.RequestsPerSecond.Period = cli.alertDuration
	conf.Alert.RequestsPerSecond.Threshold = float64(cli.alertThreshold)

	return nil
}
