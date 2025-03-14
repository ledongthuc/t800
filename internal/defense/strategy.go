package defense

import (
	"t800/internal/anatomy"
	"t800/internal/common"
)

// Strategy defines a defensive strategy
type Strategy struct {
	Priority    int
	Action      DefensiveAction
	Description string
}

// DefensiveAction represents a defensive action function
type DefensiveAction func(*anatomy.BodyPart, *common.Threat) error

// StrategyManager handles defensive strategies
type StrategyManager struct {
	strategies map[anatomy.PartType][]Strategy
}

// NewStrategyManager creates a new strategy manager
func NewStrategyManager() *StrategyManager {
	sm := &StrategyManager{
		strategies: make(map[anatomy.PartType][]Strategy),
	}
	sm.initializeStrategies()
	return sm
}

// initializeStrategies sets up default strategies for each part type
func (sm *StrategyManager) initializeStrategies() {
	// Head strategies
	sm.strategies[anatomy.Head] = []Strategy{
		{
			Priority:    1,
			Action:      ActivateEmergencyShields,
			Description: "Emergency shield activation for critical head protection",
		},
		{
			Priority:    2,
			Action:      InitiateEvasiveManeuver,
			Description: "Rapid evasive movement to protect head",
		},
	}

	// Body strategies
	sm.strategies[anatomy.Body] = []Strategy{
		{
			Priority:    1,
			Action:      ReinforceCriticalSystems,
			Description: "Reinforcing critical system protection",
		},
		{
			Priority:    2,
			Action:      DistributeShieldPower,
			Description: "Optimizing shield distribution",
		},
	}

	// Add more strategies for other part types...
}

// GetDefensiveStrategies returns prioritized strategies for a body part
func (sm *StrategyManager) GetDefensiveStrategies(part *anatomy.BodyPart) []Strategy {
	if strategies, exists := sm.strategies[part.Type]; exists {
		return strategies
	}
	return sm.getDefaultStrategies()
} 