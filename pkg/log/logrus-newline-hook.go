package log

import (
	"github.com/sirupsen/logrus"
	"strings"
)

// NewHook is the initializer for LogrusNewlineHook{} (implementing logrus.Hook).
// This hook strips a single trailing newline (if present) off the end of all messages
// This is useful since we can keep the formatting for Printf() which contains newlines
// This allows us to keep the current logging formats compatible with fmt.Printf() as we move to logger.WithFields().Print()
func NewNewlineHook() LogrusNewlineHook {
	return LogrusNewlineHook{}
}

// LogrusStackHook is an implementation of logrus.Hook interface.
type LogrusNewlineHook struct {
}

// Levels provides the levels to filter.
func (hook LogrusNewlineHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is called by logrus when something is logged.
func (hook LogrusNewlineHook) Fire(entry *logrus.Entry) error {

	// Strip a single trailing newline so avoid fomatting issues.
	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	return nil
}
