package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init() {
	zerolog.TimeFieldFormat = time.RFC3339

	log = zerolog.New(os.Stdout).With().
		Timestamp().
		// CallerWithSkipFrameCount(3).
		Logger().
		Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
			FormatCaller: func(i interface{}) string {
				if str, ok := i.(string); ok {
					parts := strings.Split(str, "/")
					if len(parts) > 0 {
						return parts[len(parts)-1]
					}
				}
				return ""
			},
		})
}

func Info(msg string, fields ...interface{}) {
	log.Info().Fields(toMap(fields...)).Msg(msg)
}

func Error(msg string, fields ...interface{}) {
	log.Error().Fields(toMap(fields...)).Msg(msg)
}

func Debug(msg string, fields ...interface{}) {
	log.Debug().Fields(toMap(fields...)).Msg(msg)
}

func Warn(msg string, fields ...interface{}) {
	log.Warn().Fields(toMap(fields...)).Msg(msg)
}

func Fatal(msg string, fields ...interface{}) {
	log.Fatal().Fields(toMap(fields...)).Msg(msg)
}

func toMap(fields ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(fields)-1; i += 2 {
		k, ok := fields[i].(string)
		if !ok {
			continue
		}
		m[k] = fields[i+1]
	}
	return m
}
