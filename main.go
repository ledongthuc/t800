package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"t800/internal/common"
	"t800/internal/processor"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create processor
	proc, err := processor.NewProcessor(ctx)
	if err != nil {
		fmt.Printf("Error creating processor: %v\n", err)
		os.Exit(1)
	}

	// Start the processor
	if err := proc.Start(); err != nil {
		fmt.Printf("Error starting processor: %v\n", err)
		os.Exit(1)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Channel to signal when threats are eliminated
	threatsEliminated := make(chan struct{})

	// Create a test threat
	threat := common.Threat{
		ID:          "THREAT-001",
		Severity:    8,
		Location:    common.Location{X: 100, Y: 100, Z: 0},
		Health:      100.0, // Initial health at 100%
		Type:        "hostile_robot",
		Description: "Hostile combat robot detected",
		Timestamp:   time.Now().Unix(),
	}

	// Report the threat
	if err := proc.ReportThreat(threat); err != nil {
		fmt.Printf("Error reporting threat: %v\n", err)
	}

	// Monitor threat health
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if activeThreat := proc.GetActiveThreat(); activeThreat != nil {
					fmt.Printf("Threat %s Health: %.2f%%\n", activeThreat.ID, activeThreat.Health)
				} else {
					// If no active threat, signal elimination and stop monitoring
					close(threatsEliminated)
					return
				}
			}
		}
	}()

	// Wait for either a signal, timeout, or all threats to be eliminated
	select {
	case <-sigChan:
		fmt.Println("\nReceived shutdown signal")
	case <-time.After(30 * time.Second):
		fmt.Println("\nTimeout reached")
	case <-threatsEliminated:
		fmt.Println("\nAll threats have been eliminated")
	}

	// Stop the processor
	if err := proc.Stop(); err != nil {
		fmt.Printf("Error stopping processor: %v\n", err)
	}
} 