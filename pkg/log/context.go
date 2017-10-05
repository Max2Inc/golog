package log

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// The primary purpose of the logs is software development (not end user consumption)
// We have two main use cases:
// 1. Development - being able to control (limit) logs. Typically we have control over the MAS and the users/meshes
//                  connected, hence may not require to enable/disable on a user or mesh basis
//    -> These systems require detailed logging in some areas (e.g. a UI query on the database) while having little or
//       no logs from other areas (e.g. polling).
// 2. Support - being able to analyse and debug logs report from other systems (lab, system test, field) Typically we have
//              little control over the connected user/meshes or traffic. Often we dont have the ability to re-run test cases
//              and are analysing after the crash.
//    -> These system should run with logging enabled to a sufficient level to permit diagnosis of problems (typically info)

// Should be moved into common context golog
// Package adds/extracts the logger as part of the context. This allows a logging context to be passed carring details
// about the context rather than duplicate logs of code/formatting with in individual log lines
// The logging context is safe for use by multiple goroutines

// A common pattern is to re-use fields between logging statements by re-using
// the logrus.Entry returned from WithFields()
// contextLogger := log.WithFields(log.Fields{
//     "common": "this is a common field",
//     "other": "I also should be logged always",
// })

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// loggerKey is the context key for the logger.  Its value of zero is
// arbitrary.  If this package defined other context keys, they would have
// different integer values.
const loggerKey key = 0


// Default global logger is used for logging when no context logger is found.
var logger = NewLogger()

// NewContext returns a new Context carrying logger .
func NewContext(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext extracts the logger from the ctx, if present.
func FromContext(ctx context.Context) (*logrus.Entry) {

	// ctx.Value returns nil if ctx has no value for the key;
	// the logrus.Entry type assertion returns ok=false for nil.
	ctxLogger, ok := ctx.Value(loggerKey).(*logrus.Entry)
	if !ok {
		return logger
	}

	return ctxLogger
}
