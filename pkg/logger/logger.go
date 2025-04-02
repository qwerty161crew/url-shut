package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
        defaultFramesSkip int = 3
)

var (
        globalOutput io.Writer
        globalLogger zerolog.Logger
)

type LoggerConfig interface {
        Level() zerolog.Level
        OutputFile() string
}

func Setup(cfg LoggerConfig) (err error) {
        globalOutput = os.Stderr

        if outputFile := cfg.OutputFile(); outputFile != "" {
                runLogFile, err := os.OpenFile(
                        outputFile,
                        os.O_APPEND|os.O_CREATE|os.O_WRONLY,
                        0664,
                )

                if err != nil {
                        msg := fmt.Sprintf("failed to open log output file %s", cfg.OutputFile())
                        log.Error().Err(err).Msg(msg)
                        return errors.Wrap(err, msg)
                }

                multi := zerolog.MultiLevelWriter(globalOutput, runLogFile)
                globalLogger = zerolog.New(multi)
        } else {
                globalLogger = zerolog.New(globalOutput)
        }

        globalLogger = globalLogger.With().Stack().Logger()
        globalLogger = globalLogger.With().Timestamp().Logger()
        globalLogger = globalLogger.With().CallerWithSkipFrameCount(defaultFramesSkip).Logger()

        zerolog.SetGlobalLevel(zerolog.Level(cfg.Level()))
        return
}

func Debug(msg string, kv ...any) {
        globalLogger.Debug().Fields(kv).Msg(msg)
}

func Info(msg string, kv ...any) {
        globalLogger.Info().Fields(kv).Msg(msg)
}

func Warn(msg string, kv ...any) {
        globalLogger.Warn().Fields(kv).Msg(msg)
}

func Error(msg string, err error, kv ...any) {
        globalLogger.Error().Fields(kv).Err(err).Msg(msg)
}
