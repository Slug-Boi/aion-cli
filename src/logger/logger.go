package logger

import "go.uber.org/zap"

// The setup logger function will live in the root command as most logging should be propagated up to the CMD commands
// This allows us to create a local logger for each command that can be used to log errors and info messages
func SetupLogger() *zap.SugaredLogger {
	// setup suggered zap logger
	// https://github.com/uber-go/zap
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	return sugar
}
