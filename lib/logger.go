package lib

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fatih/color"
)

type Logger struct {
	infoLog    *log.Logger
	debugLog   *log.Logger
	warnLog    *log.Logger
	errorLog   *log.Logger
	fatalLog   *log.Logger
	successLog *log.Logger
}

var (
	once   sync.Once
	logger *Logger
)

func GetLogger() *Logger {
	once.Do(func() {
		// Initialize the logger
		infoColor := color.New(color.FgBlue).SprintFunc()
		debugColor := color.New(color.FgCyan).SprintFunc()
		warnColor := color.New(color.FgYellow).SprintFunc()
		errorColor := color.New(color.FgRed).SprintFunc()
		successColor := color.New(color.FgGreen).SprintFunc()
		fatalColor := color.New(color.FgHiRed).SprintFunc()

		logger = &Logger{
			infoLog:    log.New(os.Stdout, infoColor("[INFO] "), log.LstdFlags),
			debugLog:   log.New(os.Stdout, debugColor("[DEBUG] "), log.LstdFlags),
			warnLog:    log.New(os.Stdout, warnColor("[WARN] "), log.LstdFlags),
			errorLog:   log.New(os.Stderr, errorColor("[ERROR] "), log.LstdFlags),
			successLog: log.New(os.Stdout, successColor("[SUCCESS] "), log.LstdFlags),
			fatalLog:   log.New(os.Stderr, fatalColor("[FATAL] "), log.LstdFlags),
		}
	})
	return logger
}

// Log methods
func (l *Logger) Info(v ...interface{}) {
	l.infoLog.Printf("%s\n", fmt.Sprint(v...))
}

func (l *Logger) Debug(v ...interface{}) {
	l.debugLog.Printf("%s\n", fmt.Sprint(v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.warnLog.Printf("%s\n", fmt.Sprint(v...))
}

func (l *Logger) Success(v ...interface{}) {
	l.successLog.Printf("%s\n", fmt.Sprint(v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.errorLog.Printf("%s\n", fmt.Sprint(v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.fatalLog.Fatalf("%s\n", fmt.Sprint(v...))
}

// Formatted Log methods
func (l *Logger) Infof(format string, v ...interface{}) {
	l.infoLog.Printf("%s\n", color.New(color.FgBlue).SprintfFunc()(format, v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.debugLog.Printf("%s\n", color.New(color.FgCyan).SprintfFunc()(format, v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.warnLog.Printf("%s\n", color.New(color.FgYellow).SprintfFunc()(format, v...))
}

func (l *Logger) Successf(format string, v ...interface{}) {
	l.successLog.Printf("%s\n", color.New(color.FgGreen).SprintfFunc()(format, v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.errorLog.Printf("%s\n", color.New(color.FgRed).SprintfFunc()(format, v...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.fatalLog.Fatalf("%s\n", color.New(color.FgHiRed).SprintfFunc()(format, v...))
}
