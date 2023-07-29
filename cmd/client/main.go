package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rnesytov/ddos_pow_protect/internal/pow"
	"github.com/rs/zerolog"
)

type Config struct {
	Addr string `default:":8080"`
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		panic(err)
	}
	pow := pow.New(pow.NewDefaultScryptConf())
	log := zerolog.New(&zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.InfoLevel).
		With().Timestamp().
		Logger()

	conn, err := net.Dial("tcp", cfg.Addr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect")
		return
	}
	defer conn.Close()

	difficulty := make([]byte, 1)
	conn.Read(difficulty)
	log.Info().Int("difficulty", int(difficulty[0])).Msg("got difficulty")
	challengeLen := make([]byte, 1)
	conn.Read(challengeLen)
	log.Info().Int("challengeLen", int(challengeLen[0])).Msg("got challengeLen")
	challenge := make([]byte, challengeLen[0])
	conn.Read(challenge)
	log.Info().Str("challenge", hex.EncodeToString(challenge)).Msg("got challenge")

	start := time.Now()
	nonce, err := pow.Solve(challenge, difficulty[0])
	if err != nil {
		log.Fatal().Err(err).Msg("failed to solve")
		return
	}
	log.Info().Uint64("nonce", nonce).Str("elapsed", time.Since(start).String()).Msg("calculated nonce")

	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, nonce)
	conn.Write(data)
	log.Info().Msg("sent nonce")

	for {
		b := make([]byte, 1)
		_, err := conn.Read(b)
		if err == nil {
			fmt.Printf("%s", b)
		} else {
			break
		}
	}
	fmt.Println()
}

func getConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
