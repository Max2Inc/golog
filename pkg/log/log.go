package log

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/neiledgar/logrus-stack"
	"github.com/sirupsen/logrus"
)

// logging guidelines.
// 1. Dont include newline statements on print()/Println() type functions
// 2. Dont pass parameters - these should be withFields
// 3. Logrus has global StandardLogger() or create your own (per package)
// 4. A common pattern is to re-use fields between logging statements by re-using
// the logrus.Entry returned from WithFields()
// contextLogger := log.WithFields(log.Fields{
//     "common": "this is a common field",
//     "other": "I also should be logged always",
// })
// contextLogger.Info("I'll be logged with common and other field")

// TODO: Migrate code to use context loggers and pass around
// TODO: Add specific package loggers - disable polling logs etc
// TODO: Add interface to change logging of running process

// Options considered:
// Dont use Prometheus logger - has no WithFields
// Dont use standard logger - has no levels and no WithFields
// Uses hook to get file/line numbers - allows versatility on a per logger /level basis
// If this proves too inefficient may have to overload logging methods [e.g. Info()] to add caller context
// e.g. https://github.com/prometheus/common/blob/master/log/log.go#L234
// This would be similar to how the standard looger does it: https://golang.org/src/log/log.go

// We use the global logger (which supports concurrency using mutex)
// We could move to separate loggers e.g. to configure levels on a per-package basis
func NewLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{}) // no custom fields for now but will add context/packages
}

// Configure global logger.
// We base this on the environment (and deliberately avoid the configuration file).
//     1. The environment is always available - even before the configuration file has been parsed.
//     2. The environment is easier to override with docker (as opposed to a file inside the containers which is more difficult)
type LogConfig struct {
	Level            string   `default:"info"`                                                    // Set the debug log level severity or above. DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel
	DisableTimestamp bool     `split_words:"true" default:"true"`                                 // Disable timestamps (useful when output is redirected to logging system that already adds timestamps)
	JsonFormat       bool     `split_words:"true" default:"false"`                                // Log in JSON format (useful when output is redirected to logging system with JSON support)
	CallerLevels     []string `split_words:"true" default:"debug,info,warning,error,fatal,panic"` // Log levels to add the caller (file/line)
	StackLevels      []string `split_words:"true" default:"panic"`                                // Log levels to add the stack
}

// Configure global logger.
// We base this on the environment (and deliberately avoid the configuration file).
//     1. The environment is always available - even before the configuration file has been parsed.
//     2. The environment is easier to override with docker (as opposed to a file inside the containers which is more difficult)
func init() {

	var logConfig LogConfig
	err := envconfig.Process("LOG", &logConfig)
	if err != nil {
		fmt.Errorf("Failed to parse log environment: %s", err)
	}

	// Validate LogLevel
	logLevel, err := logrus.ParseLevel(logConfig.Level)
	if err != nil {
		fmt.Errorf("Failed to parse LOG_LEVEL(%s) environment: %s", logConfig.Level, err)
		logLevel = logrus.InfoLevel
	}

	// Validate caller levels
	var callerLevelsStr []string
	var callerLevels []logrus.Level
	for _, level := range logConfig.CallerLevels {
		logLevel, err := logrus.ParseLevel(level)
		if err != nil {
			fmt.Errorf("Failed to parse LOG_CALLER_LEVELS(%s) environment: %s", level, err)
		} else {
			callerLevelsStr = append(callerLevelsStr, logLevel.String())
			callerLevels = append(callerLevels, logLevel)
		}
	}

	// Validate stack levels
	var stackLevelsStr []string
	var stackLevels []logrus.Level
	for _, level := range logConfig.StackLevels {
		logLevel, err := logrus.ParseLevel(level)
		if err != nil {
			fmt.Errorf("Failed to parse LOG_STACK_LEVELS(%s) environment: %s", level, err)
		} else {
			stackLevelsStr = append(stackLevelsStr, logLevel.String())
			stackLevels = append(stackLevels, logLevel)
		}
	}
	// Add a trace to standard out - useful for errors configuring logging
	//fmt.Printf("LOG_LEVEL: %s LOG_DISABLE_TIMESTAMP: %t, LOG_JSON_FORMAT: %t LOG_CALLER_LEVELS %v LOG_STACK_LEVELS %v\n",
	//	logConfig.Level, logConfig.DisableTimestamp, logConfig.JsonFormat, callerLevelsStr, stackLevelsStr)

	// Now apply the configuration
	if logConfig.JsonFormat {
		// JSON with/without timestamp
		logrus.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: logConfig.DisableTimestamp,
		})
	} else {
		// Text with/without timestamp
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,
			DisableTimestamp: logConfig.DisableTimestamp,
		})
	}
	logrus.SetLevel(logLevel)

	// Output to stdout instead of the default stderr.
	logrus.SetOutput(os.Stdout)

	// Add the stack hook.
	logrus.AddHook(logrus_stack.NewHook(callerLevels, stackLevels))

	// Add the stack hook.
	logrus.AddHook(NewNewlineHook())
}
