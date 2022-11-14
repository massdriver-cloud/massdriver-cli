package config

import (
	"context"

	log "github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	OrgID  string `env:"MASSDRIVER_ORG_ID,required"`
	APIKey string `env:"MASSDRIVER_API_KEY,required"`
}

var c Config

func Get() *Config {
	ctx := context.Background()
	err := envconfig.Process(ctx, &c)
	if err != nil {
		log.Fatal().Err(err).Msg("Required environment variable not set.")
	}

	return &c
}
