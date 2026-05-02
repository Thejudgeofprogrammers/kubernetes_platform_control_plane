package logger

import (
	clog "github.com/charmbracelet/log"
	"os"
	"strings"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(key string, value any) Logger
}

type CharmLogger struct {
	log     *clog.Logger
	allowed map[string]bool
	layer   string
}

func NewLogger(level, layers string) Logger {
	l := clog.NewWithOptions(os.Stdout, clog.Options{
		Level: parseLevel(level),
	})

	return &CharmLogger{
		log:     l,
		allowed: parseLayers(layers),
	}
}

func (l *CharmLogger) With(key string, value any) Logger {
	newLogger := l.log.With(key, value)

	newLayer := l.layer
	if key == "layer" {
		if v, ok := value.(string); ok {
			newLayer = v
		}
	}

	return &CharmLogger{
		log:     newLogger,
		allowed: l.allowed,
		layer:   newLayer,
	}
}

func (l *CharmLogger) shouldLog() bool {
	if len(l.allowed) == 0 {
		return true
	}
	return l.allowed[l.layer]
}

func (l *CharmLogger) format(msg string) string {
	if l.layer == "" {
		return msg
	}
	return "[" + strings.ToUpper(l.layer) + "] " + msg
}

func (l *CharmLogger) Debug(msg string, args ...any) {
	if l.shouldLog() {
		l.log.Debug(l.format(msg), args...)
	}
}

func (l *CharmLogger) Info(msg string, args ...any) {
	if l.shouldLog() {
		l.log.Info(l.format(msg), args...)
	}
}

func (l *CharmLogger) Warn(msg string, args ...any) {
	if l.shouldLog() {
		l.log.Warn(l.format(msg), args...)
	}
}

func (l *CharmLogger) Error(msg string, args ...any) {
	if l.shouldLog() {
		l.log.Error(l.format(msg), args...)
	}
}

func parseLayers(env string) map[string]bool {
	result := make(map[string]bool)

	if env == "" {
		return result
	}

	for _, l := range strings.Split(env, ",") {
		result[strings.TrimSpace(l)] = true
	}

	return result
}

func parseLevel(level string) clog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return clog.DebugLevel
	case "info":
		return clog.InfoLevel
	case "warn":
		return clog.WarnLevel
	case "error":
		return clog.ErrorLevel
	default:
		return clog.InfoLevel
	}
}