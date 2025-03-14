package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"t800/internal/anatomy"
	"t800/internal/ai"
	"t800/internal/common"
	"t800/internal/defense"
	"t800/internal/monitoring"
	"t800/internal/offense"
	"t800/internal/scanner"
)

// Processor represents the main T800 defensive system
type Processor struct {
	logger          *monitoring.Logger
	anatomy         *anatomy.RobotAnatomy
	defense         *defense.StrategyManager
	offense         *offense.OffenseManager
	scanner         *scanner.Scanner
	status          *Status
	location        common.Location
	speed           common.MovementSpeed
	ctx             context.Context
	cancel          context.CancelFunc
	activeThreat    *common.Threat
	engagementDistance float64
	decisionMaker   *ai.DecisionMaker
	mode            common.OperationMode
	availableWeapons []string
}

// Status maintains the processor's current state
type Status struct {
	mu       sync.RWMutex
	active   bool
	Mode     common.OperationMode
	lastScan time.Time
}

// NewProcessor creates a new T800 processor
func NewProcessor(ctx context.Context) (*Processor, error) {
	ctx, cancel := context.WithCancel(ctx)
	logger := monitoring.NewLogger()
	
	decisionMaker, err := ai.NewDecisionMaker(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create decision maker: %v", err)
	}

	return &Processor{
		logger:          logger,
		anatomy:         anatomy.NewRobotAnatomy(),
		defense:         defense.NewStrategyManager(),
		offense:         offense.NewOffenseManager(),
		scanner:         scanner.NewScanner(),
		status:          &Status{Mode: common.Normal},
		location:        common.Location{X: 0, Y: 0, Z: 0},
		speed:           common.DefaultSpeed(),
		ctx:             ctx,
		cancel:          cancel,
		engagementDistance: 50.0,
		decisionMaker:   decisionMaker,
		mode:            common.Normal,
		availableWeapons: []string{"plasma_cannon", "missile", "emp_pulse", "laser_beam"},
	}, nil
}

// Start initializes the defensive system
func (p *Processor) Start() error {
	p.logger.Info("Initializing T800 defensive system")
	
	p.status.mu.Lock()
	p.status.active = true
	p.status.mu.Unlock()

	// Start monitoring routines
	go p.monitorThreats()
	go p.monitorHealth()
	go p.scanEnvironment()

	return nil
}

// Stop safely shuts down the system
func (p *Processor) Stop() error {
	p.logger.Info("Initiating shutdown sequence")
	
	p.status.mu.Lock()
	p.status.active = false
	p.status.mu.Unlock()
	
	p.cancel()
	return nil
}

// GetStatus returns the current system status
func (p *Processor) GetStatus() Status {
	p.status.mu.RLock()
	defer p.status.mu.RUnlock()
	return *p.status
}

// GetAnatomy returns the robot's anatomy
func (p *Processor) GetAnatomy() *anatomy.RobotAnatomy {
	return p.anatomy
}

// GetActiveThreat returns the current active threat
func (p *Processor) GetActiveThreat() *common.Threat {
	return p.activeThreat
}

// ReportThreat reports a new threat to the system
func (p *Processor) ReportThreat(threat common.Threat) error {
	if !p.status.active {
		return fmt.Errorf("system is not active")
	}

	// Log the threat
	p.logger.LogThreat(threat.ID, threat.Severity, threat.Location)

	// Set as active threat and enter combat mode
	p.activeThreat = &threat
	p.status.mu.Lock()
	p.status.Mode = common.Combat
	p.status.mu.Unlock()
	p.logger.Info(fmt.Sprintf("New primary target acquired: %s (Severity: %d)", threat.ID, threat.Severity))

	// Get critical parts that need protection
	criticalParts := p.anatomy.GetCriticalParts()

	// Apply defensive strategies
	for _, part := range criticalParts {
		strategies := p.defense.GetDefensiveStrategies(part)
		for _, strategy := range strategies {
			if err := strategy.Action(part, &threat); err != nil {
				p.logger.LogError(err, "defensive action failed")
				continue
			}
			p.logger.LogDefensiveAction(strategy.Description, part.Name, true)
		}
	}

	// Add offensive response
	// Use arms for attack if available
	for _, arm := range p.anatomy.Arms {
		attackStrategies := p.offense.GetOffensiveStrategies(arm)
		for _, strategy := range attackStrategies {
			if err := strategy.Action(arm, &threat); err != nil {
				p.logger.LogError(err, "offensive action failed")
				continue
			}
			p.logger.LogDefensiveAction(strategy.Description, arm.Name, true)
		}
	}

	// Use body-mounted weapons as backup
	bodyAttacks := p.offense.GetOffensiveStrategies(p.anatomy.Body)
	for _, strategy := range bodyAttacks {
		if err := strategy.Action(p.anatomy.Body, &threat); err != nil {
			p.logger.LogError(err, "offensive action failed")
			continue
		}
		p.logger.LogDefensiveAction(strategy.Description, p.anatomy.Body.Name, true)
	}

	return nil
}

// monitorThreats continuously monitors for threats
func (p *Processor) monitorThreats() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			// Implement threat monitoring logic
		}
	}
}

// monitorHealth continuously monitors robot health
func (p *Processor) monitorHealth() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			p.anatomy.UpdateAllParts(time.Now().Unix())
			status := p.anatomy.GetHealthStatus()
			for part, health := range status {
				p.logger.LogHealthStatus(part, health, p.anatomy.IsPartCritical(part))
			}
		}
	}
}

// moveTowardsTarget moves the robot towards the current target
func (p *Processor) moveTowardsTarget(target common.Location) {
	deltaTime := 0.1 // 100ms movement update
	newLocation := p.location.MoveTowards(target, p.speed.Linear, deltaTime)
	
	// Update robot's position
	p.location = newLocation
	
	// Log movement
	distance := common.CalculateDistance(p.location, target)
	p.logger.Info(fmt.Sprintf("Moving towards target. Distance: %.2f meters", distance))
}

// scanEnvironment continuously scans for threats and processes them
func (p *Processor) scanEnvironment() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	movementTicker := time.NewTicker(100 * time.Millisecond)
	defer movementTicker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			threats := p.scanner.ScanArea(p.location)
			if err := p.processThreatsWithAI(p.ctx, threats); err != nil {
				p.logger.LogError(err, "failed to process threats with AI")
			}
		case <-movementTicker.C:
			if p.activeThreat != nil {
				if err := p.moveAndEngageWithAI(p.ctx); err != nil {
					p.logger.LogError(err, "failed to move and engage with AI")
				}
			}
		}
	}
}

// processThreatsWithAI evaluates threats using AI decision maker
func (p *Processor) processThreatsWithAI(ctx context.Context, threats []*common.Threat) error {
	if len(threats) == 0 {
		if p.activeThreat != nil {
			p.logger.Info("No threats detected, returning to normal mode")
			p.mode = common.Normal
			p.activeThreat = nil
		}
		return nil
	}

	for _, threat := range threats {
		// Skip eliminated threats
		if threat.Health <= 0 {
			continue
		}

		shouldEngage, err := p.decisionMaker.ShouldEngageProactively(ctx, *threat, p.location, p.getHealthStatus())
		if err != nil {
			return fmt.Errorf("AI decision error: %v", err)
		}

		if shouldEngage {
			p.logger.Info("New primary target acquired: " + threat.ID)
			p.activeThreat = threat
			p.mode = common.Combat
			return nil
		}
	}
	return nil
}

// moveAndEngageWithAI handles movement and combat using AI decisions
func (p *Processor) moveAndEngageWithAI(ctx context.Context) error {
	if p.activeThreat == nil {
		return nil
	}

	decision, err := p.decisionMaker.MakeCombatDecision(
		ctx,
		p.location,
		p.activeThreat,
		p.getHealthStatus(),
		p.availableWeapons,
	)
	if err != nil {
		return fmt.Errorf("AI decision error: %v", err)
	}

	switch decision.Action {
	case "move":
		p.moveTowardsTarget(p.activeThreat.Location)
	case "attack":
		p.executeAttack(decision.Weapon)
	case "defend":
		p.activateDefensiveMeasures()
	case "retreat":
		p.retreatFromThreat()
	}

	return nil
}

// executeAttack executes an attack against the current threat
func (p *Processor) executeAttack(weapon string) {
	if p.activeThreat == nil || p.activeThreat.Health <= 0 {
		return
	}

	// Calculate damage based on weapon type
	var damage float64
	switch weapon {
	case "plasma_cannon":
		damage = 25.0
	case "missile":
		damage = 40.0
	case "emp_pulse":
		damage = 15.0
	case "laser_beam":
		damage = 20.0
	default:
		damage = 10.0
	}

	// Apply damage to threat
	p.activeThreat.Health -= damage
	if p.activeThreat.Health < 0 {
		p.activeThreat.Health = 0
	}

	// Log the attack
	p.logger.Info(fmt.Sprintf("Attacked %s with %s (Damage: %.1f%%, Remaining Health: %.1f%%)",
		p.activeThreat.ID, weapon, damage, p.activeThreat.Health))

	// If threat is eliminated, clear it
	if p.activeThreat.Health <= 0 {
		p.logger.Info(fmt.Sprintf("Threat %s has been eliminated", p.activeThreat.ID))
		p.activeThreat = nil
		p.mode = common.Normal
	}
}

// activateDefensiveMeasures activates defensive systems
func (p *Processor) activateDefensiveMeasures() {
	p.logger.Info("Activating defensive measures")
	// Implement defensive measures
}

// retreatFromThreat moves away from the current threat
func (p *Processor) retreatFromThreat() {
	p.logger.Info("Retreating from threat")
	// Implement retreat logic
}

// getHealthStatus returns the current health status of all parts
func (p *Processor) getHealthStatus() map[string]float64 {
	// Implement health status retrieval
	return map[string]float64{
		"head":      100.0,
		"body":      100.0,
		"arm_left":  100.0,
		"arm_right": 100.0,
		"leg_left":  100.0,
		"leg_right": 100.0,
	}
} 