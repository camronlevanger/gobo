package utils

import (
	"log"
)

// ILogger is the interface to implement for simple console logging.
type ILogger interface {
	Info(message string)
	Warn(message string)
	Error(message string)
	Fatal(message string)
	Panic(message string)
}

// Logger is the struct for this implementation of the ILogger interface.
type Logger struct {
	noisy bool
}

// GetLogger returns a pointer to an implementation of the ILogger interface.
func GetLogger(verbose bool) ILogger {

	var logger = Logger{
		verbose,
	}

	return &logger
}

// Info calls log.Println if verbose == true.
func (logger *Logger) Info(message string) {
	if logger.noisy {
		log.Println(message)
	}
}

// Warn calls log.Println if verbose == true.
func (logger *Logger) Warn(message string) {
	if logger.noisy {
		log.Println(message)
	}
}

// Error calls log.Println.
func (logger *Logger) Error(message string) {
	log.Println(message)
}

// Fatal calls log.Fatal.
func (logger *Logger) Fatal(message string) {
	log.Fatalln(message)
}

// Panic calls log.Panic.
func (logger *Logger) Panic(message string) {
	log.Panicln(message)
}
