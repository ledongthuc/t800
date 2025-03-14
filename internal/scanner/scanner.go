package scanner

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	"t800/internal/common"
)

// Scanner represents the threat detection system
type Scanner struct {
	range_      float64
	resolution  float64
	lastScan    time.Time
	activeRange map[string]*common.Threat
	predictions map[string]*ThreatPrediction
}

// ThreatPrediction represents a predicted threat
type ThreatPrediction struct {
	Location     common.Location
	Probability  float64
	TimeToImpact float64
	Severity     int
}

// NewScanner creates a new scanner system
func NewScanner() *Scanner {
	return &Scanner{
		range_:      100.0, // 100 meter range
		resolution:  0.1,   // 10cm resolution
		activeRange: make(map[string]*common.Threat),
		predictions: make(map[string]*ThreatPrediction),
	}
}

// ScanArea performs a 360-degree scan of the surrounding area
func (s *Scanner) ScanArea(currentLocation common.Location) []*common.Threat {
	s.lastScan = time.Now()
	threats := make([]*common.Threat, 0)

	// Simulate finding threats in the area
	// This is where you would integrate with actual sensors
	for angle := 0.0; angle < 360.0; angle += 10.0 {
		// Convert angle to radians
		rad := angle * math.Pi / 180.0

		// Calculate potential threat position
		threatLoc := common.Location{
			X: currentLocation.X + s.range_*math.Cos(rad),
			Y: currentLocation.Y + s.range_*math.Sin(rad),
			Z: currentLocation.Z,
		}

		// Check for actual threats
		if s.detectThreat(threatLoc) {
			threat := &common.Threat{
				ID:        generateThreatID(),
				Type:      "unknown",
				Location:  threatLoc,
				Severity:  calculateThreatLevel(threatLoc, currentLocation),
				Timestamp: time.Now().Unix(),
			}
			threats = append(threats, threat)
			s.activeRange[threat.ID] = threat
		}

		// Check for potential threats
		if prediction := s.predictThreat(threatLoc, currentLocation); prediction != nil {
			// Convert prediction to threat if probability is high enough
			if prediction.Probability > 0.7 {
				threat := &common.Threat{
					ID:        generateThreatID(),
					Type:      "predicted",
					Location:  prediction.Location,
					Severity:  prediction.Severity,
					Timestamp: time.Now().Unix(),
				}
				threats = append(threats, threat)
				s.activeRange[threat.ID] = threat
			}
		}
	}

	return threats
}

// detectThreat simulates threat detection (replace with actual sensor logic)
func (s *Scanner) detectThreat(loc common.Location) bool {
	// Simulate random threat detection (5% chance)
	return rand.Float64() < 0.05
}

// calculateThreatLevel determines threat severity based on distance
func calculateThreatLevel(threatLoc, currentLoc common.Location) int {
	distance := common.CalculateDistance(threatLoc, currentLoc)

	// Closer threats are more severe
	if distance < 10.0 {
		return 9 // Critical
	} else if distance < 30.0 {
		return 6 // High
	} else if distance < 60.0 {
		return 3 // Medium
	}
	return 1 // Low
}

// predictThreat analyzes a location for potential threats
func (s *Scanner) predictThreat(loc, currentLoc common.Location) *ThreatPrediction {
	// Simulate threat prediction logic
	if rand.Float64() < 0.1 { // 10% chance of predicting a threat
		distance := common.CalculateDistance(loc, currentLoc)
		probability := 1.0 - (distance / s.range_)
		timeToImpact := distance / 10.0 // Assuming 10m/s movement speed
		severity := calculateThreatLevel(loc, currentLoc)

		return &ThreatPrediction{
			Location:     loc,
			Probability:  probability,
			TimeToImpact: timeToImpact,
			Severity:     severity,
		}
	}
	return nil
}

// Helper functions
func generateThreatID() string {
	return fmt.Sprintf("THREAT-%d", time.Now().UnixNano())
}

func init() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
} 