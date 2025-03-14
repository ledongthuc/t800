package monitoring

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger provides structured logging for the T800 system
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	logger := &Logger{
		Logger: logrus.New(),
	}

	// Configure logger
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	return logger
}

// LogThreat logs threat detection events
func (l *Logger) LogThreat(threatID string, severity int, location interface{}) {
	l.WithFields(logrus.Fields{
		"threat_id": threatID,
		"severity":  severity,
		"location":  location,
		"timestamp": time.Now(),
	}).Warn("Threat detected")
}

// LogDefensiveAction logs defensive actions taken
func (l *Logger) LogDefensiveAction(action string, target string, success bool) {
	l.WithFields(logrus.Fields{
		"action":    action,
		"target":    target,
		"success":   success,
		"timestamp": time.Now(),
	}).Info("Defensive action executed")
}

// LogHealthStatus logs the health status of robot parts
func (l *Logger) LogHealthStatus(partName string, health float64, isCritical bool) {
	l.WithFields(logrus.Fields{
		"part_name":  partName,
		"health":     health,
		"is_critical": isCritical,
		"timestamp":  time.Now(),
	}).Info("Health status update")
}

// LogSystemStatus logs overall system status
func (l *Logger) LogSystemStatus(mode string, active bool, partsCount int) {
	l.WithFields(logrus.Fields{
		"mode":        mode,
		"active":      active,
		"parts_count": partsCount,
		"timestamp":   time.Now(),
	}).Info("System status update")
}

// LogError logs error events with context
func (l *Logger) LogError(err error, context string) {
	l.WithFields(logrus.Fields{
		"error":     err.Error(),
		"context":   context,
		"timestamp": time.Now(),
	}).Error("System error")
} 