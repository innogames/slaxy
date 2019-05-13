package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/innogames/slaxy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfg slaxy.Config

	splash = `  _________.__                        
 /   _____/|  | _____  ___  ______.__.
 \_____  \ |  | \__  \ \  \/  <   |  |
 /        \|  |__/ __ \_>    < \___  |
/_______  /|____(____  /__/\_ \/ ____|
        \/           \/      \/\/     `
	v        = viper.New()
	logger   = logrus.New()
	slaxyCmd = &cobra.Command{
		Use:   "slaxy",
		Long:  splash,
		Short: "Sentry webhooks to slack message converter proxy",
		Run:   run,
	}
)

// logrusLogger wraps a logrus logger for compatibility with the slaxy library
type logrusLogger struct {
	slaxy.Logger
	l *logrus.Logger
}

// Debug logs debug messages
func (l *logrusLogger) Debug(msg string) {
	l.l.Debug(msg)
}

// Debugf logs debug messages
func (l *logrusLogger) Debugf(msg string, args ...interface{}) {
	l.l.Debugf(msg, args...)
}

// Info logs debug messages
func (l *logrusLogger) Info(msg string) {
	l.l.Info(msg)
}

// Infof logs debug messages
func (l *logrusLogger) Infof(msg string, args ...interface{}) {
	l.l.Infof(msg, args...)
}

// Warn logs debug messages
func (l *logrusLogger) Warn(msg string) {
	l.l.Warn(msg)
}

// Warnf logs debug messages
func (l *logrusLogger) Warnf(msg string, args ...interface{}) {
	l.l.Warnf(msg, args...)
}

// Error logs debug messages
func (l *logrusLogger) Error(msg string) {
	l.l.Error(msg)
}

// Errorf logs debug messages
func (l *logrusLogger) Errorf(msg string, args ...interface{}) {
	l.l.Errorf(msg, args...)
}

// init initializes the CLI
func init() {
	cobra.OnInitialize(loadConfig)
	slaxyCmd.PersistentFlags().StringP("config", "c", "", "path to config file if any")
	slaxyCmd.PersistentFlags().DurationP("grace-period", "g", 60*time.Second, "grace period for stopping the server")
	slaxyCmd.PersistentFlags().StringP("addr", "a", "localhost:3000", "listen address")
	slaxyCmd.PersistentFlags().StringP("token", "t", "", "slack token")
	slaxyCmd.PersistentFlags().StringP("channel", "n", "", "slack channel")
	slaxyCmd.PersistentFlags().StringSliceP("excluded-fields", "e", nil, "excluded sentry fields")

	_ = v.BindPFlag("grace-period", slaxyCmd.PersistentFlags().Lookup("grace-period"))
	_ = v.BindPFlag("addr", slaxyCmd.PersistentFlags().Lookup("addr"))
	_ = v.BindPFlag("token", slaxyCmd.PersistentFlags().Lookup("token"))
	_ = v.BindPFlag("channel", slaxyCmd.PersistentFlags().Lookup("channel"))
	_ = v.BindPFlag("excluded-fields", slaxyCmd.PersistentFlags().Lookup("excluded-fields"))
}

func main() {
	if err := slaxyCmd.Execute(); err != nil {
		logger.WithError(err).Fatal("Failed to run the command")
	}
}

// run runs the command
func run(cmd *cobra.Command, args []string) {
	l := &logrusLogger{l: logger}
	srv := slaxy.New(cfg, l)

	// run it
	logger.Info("Starting server")
	err := srv.Start()
	if err != nil {
		logger.WithError(err).Fatal("Failed starting the server")
	}

	go func() {
		for err := range srv.Errors() {
			logger.WithError(err).Error("Encountered an unexpected error")
		}
	}()

	// wait for graceful shutdown
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go handleInterrupt(srv, wg)
	wg.Wait()
	os.Exit(0)
}

// loadConfig loads and parses the config file
func loadConfig() {
	path, err := slaxyCmd.PersistentFlags().GetString("config")
	if err != nil {
		logger.WithError(err).Fatal("Could not get config path")
	}

	// configure viper
	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.SetConfigName("config")
		v.AddConfigPath("/etc/slaxy")
		v.AddConfigPath(".")
	}

	// read config
	err = v.ReadInConfig()
	if err != nil {
		logger.WithError(err).Fatal("Could not read config")
	}

	// parse config
	err = v.Unmarshal(&cfg)
	if err != nil {
		logger.WithError(err).Fatal("Could not parse config")
	}
}

// handleInterrupt takes care of signals and graceful shutdowns
func handleInterrupt(srv slaxy.Server, wg *sync.WaitGroup) {
	c := make(chan os.Signal, 1)
	defer close(c)

	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-c
	logger.Info("Shutting down. Kill again to force")
	go stop(srv, wg)

	<-c
	logger.Warn("Forced shutdown")
	os.Exit(1)
}

func stop(srv slaxy.Server, wg *sync.WaitGroup) {
	if err := srv.Stop(); err != nil {
		logger.WithError(err).Fatalf("Failed to properly stop the server")
	}

	wg.Done()
}
