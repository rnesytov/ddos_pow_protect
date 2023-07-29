package server

import (
	"context"
	"encoding/binary"
	"errors"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type (
	PoW interface {
		GetChallenge(uint) []byte
		Verify(challenge []byte, difficulty uint8, nonce uint64) (bool, error)
	}

	Config struct {
		ChallengeLen uint
		Addr         string
		Difficulty   uint8
		ReadTimeout  time.Duration
		Quotes       []string
	}

	Server struct {
		cfg *Config
		log zerolog.Logger

		listener net.Listener
		pow      PoW
	}
)

func NewDefaultConfig() *Config {
	return &Config{
		ChallengeLen: 32,
		Addr:         ":8080",
		Difficulty:   2,
		ReadTimeout:  10,
	}
}

func New(cfg *Config, log zerolog.Logger, pow PoW) (*Server, error) {
	listener, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	log.Info().Str("addr", cfg.Addr).Msg("listening")
	return &Server{
		cfg:      cfg,
		log:      log,
		listener: listener,
		pow:      pow,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.log.Info().Msg("shutting down")
				return
			default:
				conn, err := s.listener.Accept()
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						return
					}
					s.log.Error().Err(err).Msg("failed to accept connection")
					continue
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					s.handle(conn)
				}()
			}
		}
	}()

	<-ctx.Done()
	wg.Wait()
	return s.listener.Close()
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	// Setup connection
	log := s.log.With().Str("remote_addr", conn.RemoteAddr().String()).Logger()
	log.Info().Msg("new connection")
	err := conn.SetReadDeadline(time.Now().Add(s.cfg.ReadTimeout))
	if err != nil {
		log.Error().Err(err).Msg("failed to set read deadline")
		return
	}
	// Write challenge
	// packet structure:
	// 0 – difficulty
	// 1 – challenge length
	// rest – challenge
	packet := make([]byte, 2+s.cfg.ChallengeLen)
	packet[0] = s.cfg.Difficulty
	packet[1] = byte(s.cfg.ChallengeLen)
	challenge := s.pow.GetChallenge(s.cfg.ChallengeLen)
	copy(packet[2:], challenge)

	if _, err = conn.Write(packet); err != nil {
		log.Error().Err(err).Msg("failed to write challenge")
		return
	}
	// Read nonce with timeout
	nonceC := make(chan uint64)
	errC := make(chan error)
	go func() {
		defer close(nonceC)
		defer close(errC)

		resp := make([]byte, 8)
		if _, err = conn.Read(resp); err != nil {
			errC <- err
		}
		nonceC <- binary.BigEndian.Uint64(resp)
	}()
	select {
	case nonce := <-nonceC:
		// verify nonce
		ok, err := s.pow.Verify(challenge, s.cfg.Difficulty, nonce)
		if err != nil {
			log.Error().Err(err).Msg("failed to verify nonce")
			return
		}
		if ok {
			conn.Write([]byte(s.cfg.Quotes[rand.Intn(len(s.cfg.Quotes))]))
		} else {
			log.Error().Msg("invalid nonce")
		}
	case err := <-errC:
		log.Error().Err(err).Msg("failed to read nonce")
		return
	case <-time.After(s.cfg.ReadTimeout):
		log.Error().Msg("timeout")
		return
	}
}
