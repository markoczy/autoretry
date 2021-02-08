package logger

import (
	"flag"
	"fmt"
	"log"
)

func NewAutoFlagFactory() *AutoFlagFactory {
	return &AutoFlagFactory{}
}

// AutoFlagFactory out of the box flag definitions for -v, -vv and -vvv
type AutoFlagFactory struct {
	err, warn, info, debug *bool
}

func (f *AutoFlagFactory) InitFlags() {
	f.err = flag.Bool("v", false, "Enable error logs")
	f.warn = flag.Bool("vv", false, "Enable warning logs (implies -v)")
	f.info = flag.Bool("vvv", false, "Enable info logs (implies -vv)")
	f.debug = flag.Bool("vvvv", false, "Enable debug logs (implies -vvv)")
}

// Create creates a logger from the flags, must be called after flag.Parse
func (f *AutoFlagFactory) Create() Logger {
	return New(*f.err, *f.warn, *f.info, *f.debug)
}

// Logger a simple logger interface, original from github.com/markoczy/crawler
type Logger interface {
	Error(f string, v ...interface{})
	Warn(f string, v ...interface{})
	Info(f string, v ...interface{})
	Debug(f string, v ...interface{})
}

func New(err, warn, info, debug bool) Logger {
	return &logger{
		err:   err || warn || info || debug,
		warn:  warn || info || debug,
		info:  info || debug,
		debug: debug,
	}
}

type logger struct {
	err, warn, debug, info bool
}

func (l *logger) Error(f string, v ...interface{}) {
	if l.err {
		l.log("ERROR", f, v...)
	}
}

func (l *logger) Warn(f string, v ...interface{}) {
	if l.warn {
		l.log("WARN", f, v...)
	}
}

func (l *logger) Info(f string, v ...interface{}) {
	if l.info {
		l.log("INFO", f, v...)
	}
}

func (l *logger) Debug(f string, v ...interface{}) {
	if l.debug {
		l.log("DEBUG", f, v...)
	}
}

func (l *logger) log(prefix, f string, v ...interface{}) {
	log.Println(prefix, fmt.Sprintf(f, v...))
}
