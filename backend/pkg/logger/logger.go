package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	log.Logger = logger
}
