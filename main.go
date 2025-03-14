package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"t800/internal/common"
	"t800/internal/processor"
)

func main() {
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create processor
	p := processor.NewProcessor(ctx)

	// Start the system
	if err := p.Start(); err != nil {
		log.Fatalf("Failed to start system: %v", err)
	}

	// Get initial system status
	status := p.GetStatus()
	log.Printf("System started in %s mode", status.Mode)

	// Get robot anatomy information
	anatomy := p.GetAnatomy()
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
		Location:  common.Location{X: 100, Y: 100, Z: 0}, // Place threat further away to see movement
		Severity:  8, // Increase severity to ensure proactive engagement
		Timestamp: time.Now().Unix(),
		Description: "High-priority hostile target detected",
	}

	// Report the threat
	if err := p.ReportThreat(threat); err != nil {
		log.Fatalf("Failed to report threat: %v", err)
	}

	// Monitor system and wait for threat elimination
	threatEliminated := make(chan bool)
	go func() {
		for {
			status := p.GetStatus()
			if status.Mode == common.Normal && p.GetActiveThreat() == nil {
				threatEliminated <- true
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for either threat elimination or shutdown signal
	select {
	case <-sigChan:
		log.Println("Received shutdown signal")
	case <-threatEliminated:
		log.Println("All threats eliminated")
	}

	// Get final health status
	healthStatus := anatomy.GetHealthStatus()
	log.Println("Final health status:")
	for part, health := range healthStatus {
		log.Printf("- %s: %.2f%%", part, health)
	}

	// Stop the system
	if err := p.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
} 