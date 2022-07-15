package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func NewLogger(logLevel string) *logrus.Logger {
	logger := logrus.New()
	ll, err := logrus.ParseLevel(logLevel)
	if err != nil {
		ll = logrus.DebugLevel
	}
	logger.SetLevel(ll)
	logger.SetReportCaller(false)
	logger.SetFormatter(&logrus.TextFormatter{})
	return logger
}

// structuredLogger is a adaptor of logrus for chi logging middleware
type structuredLogger struct {
	*logrus.Logger
}

// RequestLogger returns a logger handler using a custom LogFormatter.
func RequestLogger(f *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&structuredLogger{Logger: f})
}

// NewLogEntry is called by chi to create log entry on start of each request
func (l *structuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	fields := logrus.Fields{
		"@timestamp": time.Now().UTC().Format(time.RFC3339),
		"method":     r.Method,
		"uri":        r.RequestURI,
	}
	return &structuredLoggerEntry{entry: logrus.NewEntry(l.Logger).WithFields(fields)}
}

// structuredLoggerEntry is an adaptor of logrus's entry for chi middleware
type structuredLoggerEntry struct {
	entry *logrus.Entry
}

// Write is called by chi at the end of each request
func (l *structuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra any) {
	l.entry = l.entry.WithFields(logrus.Fields{
		"status":  status,
		"elapsed": elapsed.Milliseconds(),
	})
	l.entry.Tracef("written %d bytes", bytes)
}

// Panic is called by chi's recoverer middleware on panic
func (l *structuredLoggerEntry) Panic(v any, stack []byte) {
	l.entry.WithField(logrus.ErrorKey, fmt.Sprintf("%+v", v)).Error(string(stack))
}

// getLogEntry returns request logger
func getLogEntry(r *http.Request) *logrus.Entry {
	if log, ok := r.Context().Value(middleware.LogEntryCtxKey).(*structuredLoggerEntry); ok {
		return log.entry
	}
	log := logrus.New()
	log.SetOutput(ioutil.Discard)
	return logrus.NewEntry(log)
}
