package logging

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.SugaredLogger

// newLumberjackLogger initializes a new lumberjack writer with the provided
// configuration options
func newLumberjackLogger(path string, maxSize int, maxBackups int, maxAge int) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	}
}

// InitLogToFile initializes a logger which outputs to a file
func InitLogToFile() {
	writer := newRollingFileWriter()

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		writer,
		zap.ErrorLevel,
	)

	log = zap.New(core).Sugar()
}

func newRollingFileWriter() zapcore.WriteSyncer {
	// TODO: Use configuration to set the dir/file
	_, err := os.Stat("log")

	if os.IsNotExist(err) {
		fmt.Printf("Log directory %s does not exist! Trying to create... ", "log")
		mkDirErr := os.Mkdir("log", 0777)

		if mkDirErr != nil {
			fmt.Fprintln(os.Stderr, "Could not create log directory! Exiting...")
			os.Exit(1)
		}

		print("Successfully created log file\n")
	}

	// TODO: These should all come from configuration
	return zapcore.AddSync(newLumberjackLogger("log/wisent.log", 20, 30, 3))
}

// InitLogToStdOut initializes a logger which outputs to standard out
func InitLogToStdOut() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()

	log = logger.Sugar()
}

func init() {
	// TODO: Get env from config
	env := "DEV"

	switch env {
	case "DEV":
		InitLogToStdOut()
	case "PROD":
		InitLogToFile()
	}

	log.Infof("Logging is enabled with %s", env)
}

// Debug logs a message with the debug log level
func Debug(msg string) {
	log.Debug(msg)
}

// Debugf logs a formatted message with the debug log level
func Debugf(templ string, args ...interface{}) {
	log.Debugf(templ, args)
}

// Info logs a message with the info log level
func Info(msg string) {
	log.Info(msg)
}

// Infof logs a formatted message with the info log level
func Infof(templ string, args ...interface{}) {
	log.Infof(templ, args)
}

// Warn logs a message with the warn log level
func Warn(msg string) {
	log.Warn(msg)
}

// Warnf logs a formatted string with the warn log level
func Warnf(templ string, args ...interface{}) {
	log.Warnf(templ, args)
}

// Error logs a message with the error log level
func Error(msg string) {
	log.Error(msg)
}

// Errorf logs a formatted message with the error log level
func Errorf(templ string, args ...interface{}) {
	log.Errorf(templ, args)
}

// Fatal logs a message with the fatal log level then calls os.Exit
func Fatal(msg string) {
	log.Fatal(msg)
}

// Fatalf logs a formatted message with the fatal log level then calls os.Exit
func Fatalf(templ string, args ...interface{}) {
	log.Fatalf(templ, args)
}

// Panic logs a message with the panic log level then panics
func Panic(msg string) {
	log.Panic(msg)
}

// Panicf logs a formatted message with the panic log level then panics
func Panicf(templ string, args interface{}) {
	log.Panicf(templ, args)
}

// With adds a variadic number of fields (key-value pairs) to the logging context
func With(args ...interface{}) {
	log.With(args)
}
