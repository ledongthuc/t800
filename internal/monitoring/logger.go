package monitoring

import (
	"fmt"
	"time"

	"t800/internal/common"
)

// Logger handles system logging
type Logger struct {
	// Add any logger-specific fields here
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{}
}

// Info logs an informational message
func (l *Logger) Info(msg string) {
	fmt.Printf("[%s] INFO: %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}

// LogThreat logs a detected threat
func (l *Logger) LogThreat(threatID string, severity int, location common.Location) {
	fmt.Printf("[%s] WARNING: Threat detected - ID: %s, Severity: %d, Location: (%.2f, %.2f, %.2f)\n",
		time.Now().Format("2006-01-02 15:04:05"),
		threatID,
		severity,
		location.X,
		location.Y,
		location.Z)
}

// LogDefensiveAction logs a defensive action
func (l *Logger) LogDefensiveAction(action, target string, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	fmt.Printf("[%s] INFO: Defensive Action - %s on %s: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		action,
		target,
		status)
}

// LogHealthStatus logs the health status of a part
func (l *Logger) LogHealthStatus(partName string, health float64, isCritical bool) {
	critical := ""
	if isCritical {
		critical = " (CRITICAL)"
	}
	fmt.Printf("[%s] INFO: Health Status - %s: %.2f%%%s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		partName,
		health,
		critical)
}

// LogSystemStatus logs the overall system status
func (l *Logger) LogSystemStatus(status string) {
	fmt.Printf("[%s] INFO: System Status - %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		status)
}

// LogError logs an error message
func (l *Logger) LogError(err error, context string) {
	fmt.Printf("[%s] ERROR: %s - %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		context,
		err)
} 