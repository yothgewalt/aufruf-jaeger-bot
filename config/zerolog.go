package config

import (
	"os"

	"github.com/rs/zerolog"
)

func NewZeroLog() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
