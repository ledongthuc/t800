package main

import (
	"context"
	"log"
	"time"

	"t800/internal/common"
	"t800/internal/processor"
)

func main() {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a new T800 processor
	t800 := processor.NewProcessor(ctx)

	// Start the system
	if err := t800.Start(); err != nil {
		log.Fatalf("Failed to start T800: %v", err)
	}
	defer t800.Stop()

	// Get initial system status
	status := t800.GetStatus()
	log.Printf("System started in %s mode", status.Mode)

	// Get robot anatomy information
	anatomy := t800.GetAnatomy()
	log.Printf("Robot total weight: %.2f kg", anatomy.CalculateTotalWeight())
	
	// Get critical parts
	criticalParts := anatomy.GetCriticalParts()
	log.Printf("Critical parts: %d", len(criticalParts))
	for _, part := range criticalParts {
		log.Printf("- %s (Health: %.2f%%)", part.Name, part.GetHealth())
	}

	// Simulate a threat
	threat := common.Threat{
		ID:        "THREAT-001",
		Type:      "physical",
		Location:  common.Location{X: 2, Y: 2, Z: 0},
		Severity:  3,
		Timestamp: time.Now().Unix(),
		Description: "Potential physical interference detected",
	}

	// Report the threat
	if err := t800.ReportThreat(threat); err != nil {
		log.Fatalf("Failed to report threat: %v", err)
	}

	// Monitor system for a while
	monitorDuration := 5 * time.Second
	log.Printf("Monitoring system for %s...", monitorDuration)
	time.Sleep(monitorDuration)

	// Get final health status
	healthStatus := anatomy.GetHealthStatus()
	log.Println("Final health status:")
	for part, health := range healthStatus {
		log.Printf("- %s: %.2f%%", part, health)
	}
} 