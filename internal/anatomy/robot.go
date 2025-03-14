package anatomy

import (
	"fmt"
	"sync"
)

// RobotAnatomy defines the physical structure of the robot
type RobotAnatomy struct {
	mu    sync.RWMutex
	Head  *BodyPart
	Body  *BodyPart
	Arms  []*BodyPart
	Legs  []*BodyPart
	Parts map[string]*BodyPart
}

// NewRobotAnatomy creates a new robot anatomy with standard T800 specifications
func NewRobotAnatomy() *RobotAnatomy {
	ra := &RobotAnatomy{
		Parts: make(map[string]*BodyPart),
	}

	// Initialize head
	headDims, _ := NewDimensions(0.3, 0.4, 0.3, 15.0)
	ra.Head = NewBodyPart(Head, "head", *headDims, true)
	ra.Parts["head"] = ra.Head

	// Initialize body
	bodyDims, _ := NewDimensions(0.5, 0.8, 0.4, 45.0)
	ra.Body = NewBodyPart(Body, "body", *bodyDims, true)
	ra.Parts["body"] = ra.Body

	// Initialize arms
	ra.Arms = make([]*BodyPart, 2)
	armDims, _ := NewDimensions(0.2, 0.7, 0.2, 20.0)
	for i := range ra.Arms {
		side := "left"
		if i == 1 {
			side = "right"
		}
		ra.Arms[i] = NewBodyPart(Arm, fmt.Sprintf("arm_%s", side), *armDims, false)
		ra.Parts[fmt.Sprintf("arm_%s", side)] = ra.Arms[i]
	}

	// Initialize legs
	ra.Legs = make([]*BodyPart, 2)
	legDims, _ := NewDimensions(0.25, 0.9, 0.25, 25.0)
	for i := range ra.Legs {
		side := "left"
		if i == 1 {
			side = "right"
		}
		ra.Legs[i] = NewBodyPart(Leg, fmt.Sprintf("leg_%s", side), *legDims, true)
		ra.Parts[fmt.Sprintf("leg_%s", side)] = ra.Legs[i]
	}

	return ra
}

// GetPart returns a body part by name
func (ra *RobotAnatomy) GetPart(name string) (*BodyPart, error) {
	ra.mu.RLock()
	defer ra.mu.RUnlock()

	part, exists := ra.Parts[name]
	if !exists {
		return nil, fmt.Errorf("part not found: %s", name)
	}
	return part, nil
}

// UpdatePart updates a body part's health
func (ra *RobotAnatomy) UpdatePart(name string, damage float64) error {
	ra.mu.Lock()
	defer ra.mu.Unlock()

	part, exists := ra.Parts[name]
	if !exists {
		return fmt.Errorf("part not found: %s", name)
	}

	part.TakeDamage(damage)
	return nil
}

// GetCriticalParts returns all critical body parts
func (ra *RobotAnatomy) GetCriticalParts() []*BodyPart {
	ra.mu.RLock()
	defer ra.mu.RUnlock()

	var critical []*BodyPart
	for _, part := range ra.Parts {
		if part.IsCritical {
			critical = append(critical, part)
		}
	}
	return critical
}

// CalculateTotalWeight returns the total weight of the robot
func (ra *RobotAnatomy) CalculateTotalWeight() float64 {
	ra.mu.RLock()
	defer ra.mu.RUnlock()

	var total float64
	for _, part := range ra.Parts {
		total += part.Dimensions.Weight
	}
	return total
}

// UpdateAllParts applies regeneration to all parts
func (ra *RobotAnatomy) UpdateAllParts(currentTime int64) {
	ra.mu.Lock()
	defer ra.mu.Unlock()

	for _, part := range ra.Parts {
		part.health.Update(currentTime)
	}
}

// IsPartCritical checks if a part is critical by its name
func (ra *RobotAnatomy) IsPartCritical(name string) bool {
	ra.mu.RLock()
	defer ra.mu.RUnlock()

	part, exists := ra.Parts[name]
	if !exists {
		return false
	}
	return part.IsCritical
}

// GetHealthStatus returns the health status of all parts
func (ra *RobotAnatomy) GetHealthStatus() map[string]float64 {
	ra.mu.RLock()
	defer ra.mu.RUnlock()

	status := make(map[string]float64)
	for name, part := range ra.Parts {
		status[name] = part.GetHealth()
	}
	return status
} 