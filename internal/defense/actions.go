package defense

import (
	"fmt"
	"t800/internal/anatomy"
	"t800/internal/common"
)

// ActivateEmergencyShields activates emergency shielding for critical parts
func ActivateEmergencyShields(part *anatomy.BodyPart, threat *common.Threat) error {
	if part == nil {
		return fmt.Errorf("invalid body part")
	}

	// Increase shield strength temporarily
	part.Protection.ShieldStrength = min(100, part.Protection.ShieldStrength*1.5)
	return nil
}

// InitiateEvasiveManeuver calculates and executes evasive movement
func InitiateEvasiveManeuver(part *anatomy.BodyPart, threat *common.Threat) error {
	if part == nil || threat == nil {
		return fmt.Errorf("invalid parameters")
	}

	// Calculate optimal evasive position
	// This would typically involve path planning and movement control
	return nil
}

// ReinforceCriticalSystems strengthens protection of critical systems
func ReinforceCriticalSystems(part *anatomy.BodyPart, threat *common.Threat) error {
	if part == nil {
		return fmt.Errorf("invalid body part")
	}

	if !part.IsCritical {
		return fmt.Errorf("part is not critical")
	}

	// Increase armor rating temporarily
	part.Protection.ArmorRating = min(100, part.Protection.ArmorRating*1.3)
	return nil
}

// DistributeShieldPower optimizes shield power distribution
func DistributeShieldPower(part *anatomy.BodyPart, threat *common.Threat) error {
	if part == nil {
		return fmt.Errorf("invalid body part")
	}

	// Optimize shield strength based on threat type and severity
	threatMultiplier := float64(threat.Severity) / 10.0
	part.Protection.ShieldStrength = min(100, part.Protection.ShieldStrength*threatMultiplier)
	return nil
}

// getDefaultStrategies returns default defensive strategies
func (sm *StrategyManager) getDefaultStrategies() []Strategy {
	return []Strategy{
		{
			Priority:    1,
			Action:      ActivateEmergencyShields,
			Description: "Standard shield activation",
		},
		{
			Priority:    2,
			Action:      InitiateEvasiveManeuver,
			Description: "Basic evasive movement",
		},
	}
}

// Helper function to calculate minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
} 