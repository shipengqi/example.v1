package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	DefaultLogDir    = "/var/run/cert-manager"
	DefaultLogLevel  = "debug"
	DefaultLogPrefix = "cm"
)

type Config struct {
	Level  string
	Output string
	Prefix string
}

var logger zerolog.Logger

func Init(c *Config) (string, error) {
	if c.Output == "" {
		c.Output = DefaultLogDir
	}
	if c.Level == "" {
		c.Level = DefaultLogLevel
	}
	if c.Prefix == "" {
		c.Prefix = DefaultLogPrefix
	}

	logLevel := convertLogLevel(c.Level)
	zerolog.SetGlobalLevel(logLevel)

	logfile, err := getLogFileName(c.Output, c.Prefix)
	if err != nil {
		return "", err
	}

	var w io.Writer

	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger = zerolog.New(w).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(f).With().Timestamp().Logger()
	}

	return logfile, nil
}

func Tracef(format string, args ...interface{}) {
	logger.Trace().Msgf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debug().Msgf(format, args...)
}

func Infof(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	logger.Info().Msgf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	logger.Warn().Msgf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	logger.Error().Msgf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	logger.Fatal().Msgf(format, args...)
}

func Trace(msg string) {
	logger.Trace().Msg(msg)
}

func Debug(msg string) {
	logger.Debug().Msg(msg)
}

func Info(msg string) {
	fmt.Println(msg)
	logger.Info().Msg(msg)
}

func Warn(msg string) {
	fmt.Println(msg)
	logger.Warn().Msg(msg)
}

func Error(msg string) {
	fmt.Println(msg)
	logger.Error().Msg(msg)
}

func Fatal(msg string) {
	fmt.Println(msg)
	logger.Fatal().Msg(msg)
}

func Print(msg string) {
	fmt.Print(msg)
}

func getLogFileName(out, prefix string) (string, error) {
	if !checkDir(out) {
		if err := os.MkdirAll(out, 0644); err != nil {
			return "", err
		}
	}
	logFileName := filepath.Join(out, fmt.Sprintf("%s.%s.log", prefix, time.Now().Format("20060102150405")))
	return logFileName, nil
}

func convertLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "trace":
		return zerolog.TraceLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.DebugLevel
	}
}

// ----------------------------------------------------------------------------
// Helpers...

func checkDir(name string) bool {
	if f, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		if !f.IsDir() {
			return false
		}
	}
	return true
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
