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

// For context logging we want different contexts to support different log levels.
// We do not want to create a new logger for each context (goroutine)
//  - https://stackoverflow.com/questions/18361750/correct-approach-to-global-logging-in-golang
// So we have separate loggers for each level, but share the same logger between goroutines for a specific log level.
// The logger supports concurrent write access
func NewPkgLogger(envPrefix string) *logrus.Logger {

	theLogger := logrus.New()
	Setup(envPrefix, theLogger)

	return theLogger
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

// Setup the logger from the environment
func Setup(envPrefix string, theLogger *logrus.Logger) {
	var logConfig LogConfig
	err := envconfig.Process(envPrefix, &logConfig)
	if err != nil {
		fmt.Errorf("Failed to parse log environment: %s", err)
	}

	// Validate LogLevel
	logLevel, err := logrus.ParseLevel(logConfig.Level)
	if err != nil {
		fmt.Errorf("Failed to parse %s_LEVEL(%s) environment: %s", envPrefix, logConfig.Level, err)
		logLevel = logrus.InfoLevel
	}

	// Validate caller levels
	var callerLevelsStr []string
	var callerLevels []logrus.Level
	for _, level := range logConfig.CallerLevels {
		logLevel, err := logrus.ParseLevel(level)
		if err != nil {
			fmt.Errorf("Failed to parse %s_CALLER_LEVELS(%s) environment: %s", envPrefix, level, err)
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
			fmt.Errorf("Failed to parse %s_STACK_LEVELS(%s) environment: %s", envPrefix, level, err)
		} else {
			stackLevelsStr = append(stackLevelsStr, logLevel.String())
			stackLevels = append(stackLevels, logLevel)
		}
	}
	// Add a trace to standard out - useful for errors configuring logging
	//fmt.Printf("%s_LEVEL: %s %s_DISABLE_TIMESTAMP: %t, %s_JSON_FORMAT: %t %s_CALLER_LEVELS %v %s_STACK_LEVELS %v\n",
	//	envPrefix, logConfig.Level,
	//	envPrefix, logConfig.DisableTimestamp,
	//	envPrefix, logConfig.JsonFormat,
	//	envPrefix, callerLevelsStr,
	//	envPrefix, stackLevelsStr)

	// Now apply the configuration
	if logConfig.JsonFormat {
		// JSON with/without timestamp
		theLogger.Formatter = &logrus.JSONFormatter{
			DisableTimestamp: logConfig.DisableTimestamp,
		}
	} else {
		// Text with/without timestamp
		theLogger.Formatter = &logrus.TextFormatter{
			FullTimestamp:    true,
			DisableTimestamp: logConfig.DisableTimestamp,
		}
	}
	theLogger.SetLevel(logLevel)

	// Output to stdout instead of the default stderr.
	theLogger.Out = os.Stdout

	// Add the stack hook.
	theLogger.AddHook(logrus_stack.NewHook(callerLevels, stackLevels))

	// Add the stack hook.
	theLogger.AddHook(NewNewlineHook())
}

// Configure global logger.
// We base this on the environment (and deliberately avoid the configuration file).
//     1. The environment is always available - even before the configuration file has been parsed.
//     2. The environment is easier to override with docker (as opposed to a file inside the containers which is more difficult)
func init() {
	// This bypasses the mutex protection setting up the logrus global logger - but is only called from this package init
	// Alternative would be to create a new vsys specific "global" logger - since all vsys packages use NewLogger()
	Setup("LOG", logrus.StandardLogger())
}
