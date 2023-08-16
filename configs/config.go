package configs

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Server struct {
		Debug              bool   `yaml:"debug", envconfig:"SERVER_DEBUG"`
		Host               string `yaml:"host", envconfig:"SERVER_HOST"`
		Port               int    `yaml:"port", envconfig:"SERVER_PORT"`
		WriteTimeout       int    `yaml:"write_timeout", envconfig:"SERVER_WRITE_TIMEOUT"`
		ReadTimeout        int    `yaml:"read_timeout", envconfig:"SERVER_READ_TIMEOUT"`
		WsWriteBufferSize  int    `yaml:"ws_write_buffer_size", envconfig:"SERVER_WS_WRITE_BUFFER_SIZE"`
		WsReadBufferSize   int    `yaml:"ws_read_buffer_sie", envconfig:"SERVER_WS_READ_BUFFER_SIZE"`
		WsHandshakeTimeout int    `yaml:"ws_handshake_timeout", envconfig:"SERVER_WS_HANDSHAKE_TIMEOUT"`
	} `yaml:"server"`

	Game struct {
		Fps             int  `yaml:"fps", envconfig:"GAME_FPS"`
		UniversePrepend bool `yaml:"universe_prepend", envconfig:"GAME_UNIVERSE_PREPEND"`
	} `yaml:"game"`

	Log struct {
		Level           string `yaml:"level", envconfig:"LOG_LEVEL"`
		SetReportCaller bool   `yaml:"set_report_caller", envconfig:"LOG_SET_REPORT_CALLER"`
	} `yaml:"log"`
}

func NewConfig() *Config {
	var cfg Config

	// Read default.yml file
	f, err := os.Open("configs/default.yml")
	if err != nil {
		handleError(err)
	}
	defer f.Close()

	// Decode into Config struct instance.
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		handleError(err)
	}

	// Read env vars and update default values.
	err = envconfig.Process("", &cfg)
	if err != nil {
		handleError(err)
	}

	return &cfg
}

func handleError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
