package main

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type config struct {
	Version         string        `mapstructure:"VERSION"`
	APIHost         string        `mapstructure:"API_HOST"`
	DebugHost       string        `mapstructure:"DEBUG_HOST"`
	ReadTimeout     time.Duration `mapstructure:"READ_TIMEOUT"`
	WriteTimeout    time.Duration `mapstructure:"WRITE_TIMEOUT"`
	IdleTimeout     time.Duration `mapstructure:"IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
}

var cfg *config

var build = "develop"

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("server")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%v", err)
	}

	cfg = &config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}
}

func main() {
	fmt.Println("")
	fmt.Println(build)
	fmt.Println(cfg)
}
