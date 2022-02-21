package logging

import (
	stdLog "log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func InitZerolog(levelInput string) {
	level, err := zerolog.ParseLevel(levelInput)
	if err != nil {
		log.Panic()
	}

	zerolog.SetGlobalLevel(level)
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: false,
	})

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339

	log.Logger = log.Logger.With().Logger()

	stdLog.SetFlags(0)
	stdLog.SetOutput(log.Logger)
}
