package main

import (
	"context"
	"expvar"
	"fmt"
	"os"
	"runtime"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/mroobert/go-tickets/auth/internal/foundation/logger"
	"github.com/spf13/viper"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

var build = "develop"

func main() {
	// Construct the application logger.
	log, err := logger.New("AUTH-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	// =========================================================================
	// GOMAXPROCS

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}

	// =========================================================================
	// Configuration

	cfg := struct {
		Version string `mapstructure:"VERSION"`
		Web     struct {
			APIHost         string        `mapstructure:"API_HOST"`
			DebugHost       string        `mapstructure:"DEBUG_HOST"`
			ReadTimeout     time.Duration `mapstructure:"READ_TIMEOUT"`
			WriteTimeout    time.Duration `mapstructure:"WRITE_TIMEOUT"`
			IdleTimeout     time.Duration `mapstructure:"IDLE_TIMEOUT"`
			ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
		}
	}{
		Version: build,
	}

	viper.AddConfigPath("./cmd/server/")
	viper.SetConfigName("server")
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%v", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	// =========================================================================
	// App Starting

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "config", cfg)
	defer log.Infow("shutdown complete")

	expvar.NewString("build").Set(cfg.Version)

	// =========================================================================
	// Initialize Firebase Support

	fbClient, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error initializing firebase client: %w", err)
	}

	fbAuthClient, err := fbClient.Auth(context.Background())
	if err != nil {
		return fmt.Errorf("error initializing firebase auth client: %w", err)
	}

	fmt.Println(fbAuthClient)
	return nil
}
