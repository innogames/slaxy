package slaxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

type handler func(l net.Listener)

// Config holds all config values
type Config struct {
	GracePeriod    time.Duration `mapstructure:"grace-period"`
	Addr           string
	Token          string `mapstructure:"bot-token"`
	Channel        string
	ExcludedFields []string `mapstructure:"excluded-fields"`
}

// server types
type server struct {
	cfg            Config
	logger         Logger
	done           chan struct{}
	srv            *http.Server
	wg             *sync.WaitGroup
	errChan        chan error
	slack          *slack.Client
	excludedFields []*regexp.Regexp
}

// Server represents a server instance
type Server interface {
	Start() error
	Stop() error
	Errors() <-chan error
}

// New creates a new server instance
func New(cfg Config, logger Logger) Server {
	return &server{
		cfg:     cfg,
		logger:  logger,
		done:    make(chan struct{}, 1),
		wg:      new(sync.WaitGroup),
		errChan: make(chan error, 100),
	}
}

// Start starts up the server
func (s *server) Start() error {
	return s.setup(s.cfg.Addr, s.handleWeb)
}

// Stop gracefully shuts down the server
func (s *server) Stop() error {
	s.done <- struct{}{}

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.GracePeriod)
	err := s.srv.Shutdown(ctx)
	cancel()

	s.wg.Wait()

	return err
}

// Errors returns the error channel
func (s *server) Errors() <-chan error {
	return s.errChan
}

// setup starts up a server with its own listener and handler function
func (s *server) setup(addr string, handler handler) error {
	// pre-compile regexes
	excludedFields := make([]*regexp.Regexp, 0, len(s.cfg.ExcludedFields))
	for _, regex := range s.cfg.ExcludedFields {
		excludedFields = append(excludedFields, regexp.MustCompile(regex))
	}
	s.excludedFields = excludedFields

	// connect to slack
	client := slack.New(s.cfg.Token)
	_, err := client.AuthTest()
	if err != nil {
		return err
	}
	s.slack = client

	// start tcp listener
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on %s", addr)
	}

	s.logger.Info(fmt.Sprintf("Listening on %s", addr))
	go s.handleListener(l, addr, handler)

	return nil
}

// handleListener handles a listener using the specified handler function
func (s *server) handleListener(l net.Listener, addr string, handler handler) {
	defer s.logger.Info(fmt.Sprintf("Listener %s shutdown", addr))

	handler(l)
}

// handleWeb handles all incoming connections to the webhook server
func (s *server) handleWeb(l net.Listener) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleWebhook)

	s.srv = &http.Server{
		Handler: mux,
	}
	err := s.srv.Serve(l)

	// server closed abnormally
	if err != nil && err != http.ErrServerClosed {
		err = errors.Wrap(err, "server failed")
		s.errChan <- err
	}
}
