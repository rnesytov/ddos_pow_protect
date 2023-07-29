package main

import (
	"context"
	_ "embed"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"

	"github.com/rnesytov/ddos_pow_protect/internal/pow"
	"github.com/rnesytov/ddos_pow_protect/internal/server"
)

//go:embed quotes.txt
var rawQuotes []byte

type Config struct {
	Addr        string        `default:":8080"`
	Difficulty  uint8         `default:"2"`
	ReadTimeout time.Duration `default:"10s"`
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		panic(err)
	}
	log := zerolog.New(&zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.InfoLevel).
		With().Timestamp().
		Logger()

	srvConfig := server.NewDefaultConfig()
	srvConfig.Addr = cfg.Addr
	srvConfig.Difficulty = cfg.Difficulty
	srvConfig.ReadTimeout = cfg.ReadTimeout
	srvConfig.Quotes = getQuotes()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pow := pow.New(pow.NewDefaultScryptConf())
	server, err := server.New(srvConfig, log, pow)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server")
		return
	}

	if err := server.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to run server")
	}
}

func getConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func getQuotes() []string {
	return strings.Split(string(rawQuotes), "\n")
}
