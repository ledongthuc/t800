package anatomy

import (
	"fmt"
	"sync"
)

// SafeHealth provides thread-safe health management
type SafeHealth struct {
	mu         sync.RWMutex
	current    float64
	maximum    float64
	regenRate  float64
	lastUpdate int64
}

// NewSafeHealth creates a new SafeHealth instance
func NewSafeHealth(maximum float64) *SafeHealth {
	return &SafeHealth{
		current:    maximum,
		maximum:    maximum,
		regenRate:  0.1, // 10% regeneration per second
		lastUpdate: 0,
	}
}

// Get returns the current health value
func (h *SafeHealth) Get() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.current
}

// Reduce decreases health by the specified amount
func (h *SafeHealth) Reduce(amount float64) float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	if amount <= 0 {
		return 0
	}

	actualDamage := min(amount, h.current)
	h.current = max(0, h.current-actualDamage)
	return actualDamage
}

// Heal increases health by the specified amount
func (h *SafeHealth) Heal(amount float64) float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	if amount <= 0 {
		return 0
	}

	actualHealing := min(amount, h.maximum-h.current)
	h.current = min(h.maximum, h.current+actualHealing)
	return actualHealing
}

// SetRegenRate sets the health regeneration rate
func (h *SafeHealth) SetRegenRate(rate float64) error {
	if rate < 0 {
		return fmt.Errorf("regeneration rate cannot be negative")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.regenRate = rate
	return nil
}

// Update applies regeneration over time
func (h *SafeHealth) Update(currentTime int64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.lastUpdate == 0 {
		h.lastUpdate = currentTime
		return
	}

	deltaTime := float64(currentTime-h.lastUpdate) / 1000.0 // Convert to seconds
	if deltaTime <= 0 {
		return
	}

	regenAmount := h.regenRate * deltaTime * h.maximum
	h.current = min(h.maximum, h.current+regenAmount)
	h.lastUpdate = currentTime
}

// IsCritical returns true if health is below 20%
func (h *SafeHealth) IsCritical() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.current < (h.maximum * 0.2)
}

// Percentage returns the current health as a percentage
func (h *SafeHealth) Percentage() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return (h.current / h.maximum) * 100
} 