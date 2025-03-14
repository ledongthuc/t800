package offense

import (
	"fmt"
	"t800/internal/anatomy"
	"t800/internal/common"
)

// AttackAction represents an offensive action function
type AttackAction func(*anatomy.BodyPart, *common.Threat) error

// AttackStrategy defines an offensive strategy
type AttackStrategy struct {
	Priority    int
	Action      AttackAction
	Description string
	PowerUsage  float64
	Range       float64
	Preemptive  bool
}

// PlasmaCannonAttack fires a concentrated plasma beam
func PlasmaCannonAttack(part *anatomy.BodyPart, threat *common.Threat) error {
	if part == nil || threat == nil {
		return fmt.Errorf("invalid parameters")
	}

	// Only arms can perform plasma cannon attacks
	if part.Type != anatomy.Arm {
		return fmt.Errorf("plasma cannon can only be fired from arms")
	}

	// Log attack execution
	fmt.Printf("Firing plasma cannon from %s at threat %s\n", part.Name, threat.ID)
	return nil
}

// MissileLaunch launches guided missiles at the threat
func MissileLaunch(part *anatomy.BodyPart, threat *common.Threat) error {
	if part == nil || threat == nil {
		return fmt.Errorf("invalid parameters")
	}

	if part.Type != anatomy.Body {
		return fmt.Errorf("missiles can only be launched from body")
	}

	// Log missile launch
	fmt.Printf("Launching missiles from %s at threat %s\n", part.Name, threat.ID)
	return nil
}

// EMPPulse generates an electromagnetic pulse to disable electronic threats
func EMPPulse(part *anatomy.BodyPart, threat *common.Threat) error {
	if part == nil || threat == nil {
		return fmt.Errorf("invalid parameters")
	}

	if part.Type != anatomy.Body {
		return fmt.Errorf("EMP can only be generated from body")
	}

	// Log EMP activation
	fmt.Printf("Activating EMP pulse from %s at threat %s\n", part.Name, threat.ID)
	return nil
}

// LaserBeam fires a high-energy laser beam
func LaserBeam(part *anatomy.BodyPart, threat *common.Threat) error {
	if part == nil || threat == nil {
		return fmt.Errorf("invalid parameters")
	}

	if part.Type != anatomy.Head {
		return fmt.Errorf("laser can only be fired from head")
	}

	// Log laser activation
	fmt.Printf("Firing laser beam from %s at threat %s\n", part.Name, threat.ID)
	return nil
}

// OffenseManager handles offensive strategies
type OffenseManager struct {
	strategies map[anatomy.PartType][]AttackStrategy
}

// NewOffenseManager creates a new offense manager
func NewOffenseManager() *OffenseManager {
	om := &OffenseManager{
		strategies: make(map[anatomy.PartType][]AttackStrategy),
	}
	om.initializeStrategies()
	return om
}

// initializeStrategies sets up default attack strategies for each part type
func (om *OffenseManager) initializeStrategies() {
	// Arm strategies
	om.strategies[anatomy.Arm] = []AttackStrategy{
		{
			Priority:    1,
			Action:      PlasmaCannonAttack,
			Description: "Plasma cannon attack",
			PowerUsage:  75.0,
			Range:       50.0,
			Preemptive:  true,
		},
	}

	// Body strategies
	om.strategies[anatomy.Body] = []AttackStrategy{
		{
			Priority:    2,
			Action:      MissileLaunch,
			Description: "Guided missile launch",
			PowerUsage:  90.0,
			Range:       100.0,
			Preemptive:  true,
		},
		{
			Priority:    3,
			Action:      EMPPulse,
			Description: "EMP pulse",
			PowerUsage:  85.0,
			Range:       30.0,
			Preemptive:  true,
		},
	}

	// Head strategies
	om.strategies[anatomy.Head] = []AttackStrategy{
		{
			Priority:    4,
			Action:      LaserBeam,
			Description: "Laser beam attack",
			PowerUsage:  60.0,
			Range:       40.0,
			Preemptive:  true,
		},
	}
}

// GetOffensiveStrategies returns available attack strategies for a body part
func (om *OffenseManager) GetOffensiveStrategies(part *anatomy.BodyPart) []AttackStrategy {
	if strategies, exists := om.strategies[part.Type]; exists {
		return strategies
	}
	return nil
}

// GetPreemptiveStrategies returns strategies that can be used for preemptive strikes
func (om *OffenseManager) GetPreemptiveStrategies(part *anatomy.BodyPart) []AttackStrategy {
	allStrategies := om.GetOffensiveStrategies(part)
	preemptive := make([]AttackStrategy, 0)
	for _, strategy := range allStrategies {
		if strategy.Preemptive {
			preemptive = append(preemptive, strategy)
		}
	}
	return preemptive
} 