package processor

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"t800/internal/anatomy"
	"t800/internal/common"
	"t800/internal/defense"
	"t800/internal/monitoring"
	"t800/internal/offense"
	"t800/internal/scanner"
)

// Processor represents the main T800 defensive system
type Processor struct {
	logger     *monitoring.Logger
	anatomy    *anatomy.RobotAnatomy
	defense    *defense.StrategyManager
	offense    *offense.OffenseManager
	scanner    *scanner.Scanner
	status     *Status
	location   common.Location
	speed      common.MovementSpeed
	ctx        context.Context
	cancel     context.CancelFunc
	activeThreat *common.Threat
	engagementDistance float64
}

// Status maintains the processor's current state
type Status struct {
	mu       sync.RWMutex
	active   bool
	Mode     common.OperationMode
	lastScan time.Time
}

// NewProcessor creates a new T800 processor
func NewProcessor(ctx context.Context) *Processor {
	ctx, cancel := context.WithCancel(ctx)
	return &Processor{
		logger:     monitoring.NewLogger(),
		anatomy:    anatomy.NewRobotAnatomy(),
		defense:    defense.NewStrategyManager(),
		offense:    offense.NewOffenseManager(),
		scanner:    scanner.NewScanner(),
		status:     &Status{Mode: common.Normal},
		location:   common.Location{X: 0, Y: 0, Z: 0},
		speed:      common.DefaultSpeed(),
		ctx:        ctx,
		cancel:     cancel,
		engagementDistance: 20.0, // Optimal engagement distance in meters
	}
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

// scanEnvironment continuously scans the environment and responds to threats
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
			p.processThreats(threats)
		case <-movementTicker.C:
			if p.activeThreat != nil {
				p.moveAndEngage()
			}
		}
	}
}

// processThreats evaluates and prioritizes threats
func (p *Processor) processThreats(threats []*common.Threat) {
	if len(threats) == 0 {
		return
	}

	// Find the highest severity threat
	var highestThreat *common.Threat
	highestSeverity := -1

	for _, threat := range threats {
		// Log detected threat
		p.logger.LogThreat(threat.ID, threat.Severity, threat.Location)

		if threat.Severity > highestSeverity {
			highestSeverity = threat.Severity
			highestThreat = threat
		}
	}

	// Update active threat if we found a more severe one
	if highestThreat != nil && (p.activeThreat == nil || highestThreat.Severity > p.activeThreat.Severity) {
		p.activeThreat = highestThreat
		p.status.mu.Lock()
		p.status.Mode = common.Combat
		p.status.mu.Unlock()
		p.logger.Info(fmt.Sprintf("New primary target acquired: %s (Severity: %d)", highestThreat.ID, highestThreat.Severity))
	}
}

// moveAndEngage handles movement and combat with the active threat
func (p *Processor) moveAndEngage() {
	if p.activeThreat == nil {
		return
	}

	distance := common.CalculateDistance(p.location, p.activeThreat.Location)
	
	// If we're at optimal engagement distance, focus on attack
	if math.Abs(distance - p.engagementDistance) < 1.0 {
		p.engageTarget(p.activeThreat)
	} else if distance > p.engagementDistance {
		// Move to optimal engagement distance
		targetPos := p.activeThreat.Location
		p.moveTowardsTarget(targetPos)
		
		// Execute preemptive strikes while moving
		if p.shouldEngageProactively(p.activeThreat) {
			p.executePreemptiveStrike(p.activeThreat)
		}
	} else {
		// Too close, back up while attacking
		// Calculate backup position
		direction := common.Location{
			X: p.location.X - p.activeThreat.Location.X,
			Y: p.location.Y - p.activeThreat.Location.Y,
			Z: p.location.Z - p.activeThreat.Location.Z,
		}
		distance := common.CalculateDistance(p.location, p.activeThreat.Location)
		scale := p.engagementDistance / distance
		backupPos := common.Location{
			X: p.activeThreat.Location.X + direction.X * scale,
			Y: p.activeThreat.Location.Y + direction.Y * scale,
			Z: p.activeThreat.Location.Z + direction.Z * scale,
		}
		p.moveTowardsTarget(backupPos)
		p.engageTarget(p.activeThreat)
	}

	// Clear active threat if it's been eliminated
	// In a real system, this would be based on threat detection updates
	if rand.Float64() < 0.1 { // 10% chance of eliminating threat each engagement
		p.logger.Info(fmt.Sprintf("Target %s eliminated", p.activeThreat.ID))
		p.activeThreat = nil
		p.status.mu.Lock()
		p.status.Mode = common.Normal
		p.status.mu.Unlock()
	}
}

// shouldEngageProactively determines if the robot should attack proactively
func (p *Processor) shouldEngageProactively(threat *common.Threat) bool {
	// Define engagement criteria
	if threat.Severity >= 6 {
		// Always engage high-severity threats
		return true
	}

	distance := common.CalculateDistance(threat.Location, p.location)
	if distance < 30.0 && threat.Severity >= 3 {
		// Engage medium-severity threats within 30 meters
		return true
	}

	return false
}

// engageTarget initiates offensive actions against a target
func (p *Processor) engageTarget(threat *common.Threat) {
	p.logger.LogDefensiveAction("Initiating proactive engagement", "system", true)

	// First, prepare defensive measures
	criticalParts := p.anatomy.GetCriticalParts()
	for _, part := range criticalParts {
		strategies := p.defense.GetDefensiveStrategies(part)
		for _, strategy := range strategies {
			if err := strategy.Action(part, threat); err != nil {
				p.logger.LogError(err, "defensive preparation failed")
			}
		}
	}

	// Launch coordinated attack
	p.executeCoordinatedAttack(threat)
}

// executeCoordinatedAttack performs a coordinated attack using all available weapons
func (p *Processor) executeCoordinatedAttack(threat *common.Threat) {
	var wg sync.WaitGroup

	// Use both arms simultaneously for maximum effect
	for _, arm := range p.anatomy.Arms {
		attackStrategies := p.offense.GetOffensiveStrategies(arm)
		for _, strategy := range attackStrategies {
			wg.Add(1)
			go func(arm *anatomy.BodyPart, strategy offense.AttackStrategy) {
				defer wg.Done()
				if err := strategy.Action(arm, threat); err != nil {
					p.logger.LogError(err, "arm attack failed")
				}
			}(arm, strategy)
		}
	}

	// Launch missiles as secondary attack
	bodyAttacks := p.offense.GetOffensiveStrategies(p.anatomy.Body)
	for _, strategy := range bodyAttacks {
		if err := strategy.Action(p.anatomy.Body, threat); err != nil {
			p.logger.LogError(err, "body attack failed")
		}
	}

	// Wait for all arm attacks to complete
	wg.Wait()
}

// executePreemptiveStrike launches preemptive attacks against predicted threats
func (p *Processor) executePreemptiveStrike(threat *common.Threat) {
	p.logger.LogDefensiveAction("Initiating preemptive strike", "system", true)

	// Use long-range weapons first
	bodyAttacks := p.offense.GetPreemptiveStrategies(p.anatomy.Body)
	for _, strategy := range bodyAttacks {
		if err := strategy.Action(p.anatomy.Body, threat); err != nil {
			p.logger.LogError(err, "preemptive strike failed")
		}
	}

	// Prepare defensive measures while engaging
	criticalParts := p.anatomy.GetCriticalParts()
	for _, part := range criticalParts {
		strategies := p.defense.GetDefensiveStrategies(part)
		for _, strategy := range strategies {
			if err := strategy.Action(part, threat); err != nil {
				p.logger.LogError(err, "defensive preparation failed")
			}
		}
	}

	// Use remaining weapons as backup
	for _, arm := range p.anatomy.Arms {
		attackStrategies := p.offense.GetPreemptiveStrategies(arm)
		for _, strategy := range attackStrategies {
			if err := strategy.Action(arm, threat); err != nil {
				p.logger.LogError(err, "arm preemptive strike failed")
			}
		}
	}
} 