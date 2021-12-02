package logger

import (
	"fmt"
	"github.com/shipengqi/example.v1/apps/blog/pkg/utils"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	DefaultLogDir    = "/var/run/logs"
	DefaultLogLevel  = "debug"
	DefaultLogPrefix = "demo"
)

var logger zerolog.Logger

type Config struct {
	Level  string `ini:"LEVEL"`
	Output string `ini:"OUTPUT"`
	Prefix string `ini:"PREFIX"`
}

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

	// log output to files as well
	var w io.Writer

	filename, err := initLogFileName(c.Output, c.Prefix)
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		w = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: isWindows()}
	} else {
		w = zerolog.MultiLevelWriter(
			zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: isWindows()},
			f,
		)
	}

	logger = zerolog.New(w).With().Timestamp().Logger()
	return filename, nil
}

func initLogFileName(out, prefix string) (string, error) {
	if err := utils.IsNotExistMkDir(out); err != nil {
		return "", err
	}

	logFileName := filepath.Join(out, fmt.Sprintf("%s.%s.log", prefix, time.Now().Format("20060102150405")))
	return logFileName, nil
}

func convertLogLevel(level string) zerolog.Level {
	level = strings.TrimSpace(level)
	level = strings.ToLower(level)
	switch level {
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

func Trace() *zerolog.Event {
	return logger.Trace()
}

func Debug() *zerolog.Event {
	return logger.Debug()
}

func Info() *zerolog.Event {
	return logger.Info()
}

func Warn() *zerolog.Event {
	return logger.Warn()
}

func Error() *zerolog.Event {
	return logger.Error()
}

func Fatal() *zerolog.Event {
	return logger.Fatal()
}

func Panic() *zerolog.Event {
	return logger.Panic()
}

// ----------------------------------------------------------------------------
// Helpers...

func isExist(name string) bool {
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
