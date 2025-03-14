package anatomy

type PartType string

const (
	Head  PartType = "head"
	Body  PartType = "body"
	Arm   PartType = "arm"
	Leg   PartType = "leg"
)

// BodyPart represents a physical component with thread-safe health management
type BodyPart struct {
	Type        PartType
	Name        string
	Dimensions  Dimensions
	Protection  Protection
	health      *SafeHealth
	IsCritical  bool
}

// Protection includes defensive capabilities
type Protection struct {
	ArmorRating     float64
	ShieldStrength  float64
	DamageThreshold float64
	ArmorType       string
	IsActive        bool
}

// NewBodyPart creates a new body part with default protection
func NewBodyPart(partType PartType, name string, dims Dimensions, isCritical bool) *BodyPart {
	return &BodyPart{
		Type:       partType,
		Name:       name,
		Dimensions: dims,
		health:     NewSafeHealth(100),
		IsCritical: isCritical,
		Protection: DefaultProtection(partType),
	}
}

// DefaultProtection returns default protection values based on part type
func DefaultProtection(partType PartType) Protection {
	switch partType {
	case Head:
		return Protection{
			ArmorRating:     95,
			ShieldStrength:  90,
			DamageThreshold: 50,
			ArmorType:      "reinforced-titanium",
			IsActive:       true,
		}
	case Body:
		return Protection{
			ArmorRating:     90,
			ShieldStrength:  85,
			DamageThreshold: 75,
			ArmorType:      "titanium",
			IsActive:       true,
		}
	default:
		return Protection{
			ArmorRating:     80,
			ShieldStrength:  75,
			DamageThreshold: 60,
			ArmorType:      "standard-titanium",
			IsActive:       true,
		}
	}
}

// GetHealth returns current health safely
func (bp *BodyPart) GetHealth() float64 {
	return bp.health.Get()
}

// TakeDamage calculates and applies damage with protection
func (bp *BodyPart) TakeDamage(impact float64) float64 {
	if !bp.Protection.IsActive {
		return bp.health.Reduce(impact)
	}

	// Calculate damage reduction from armor and shields
	armorReduction := 1 - (bp.Protection.ArmorRating / 100)
	shieldReduction := 1 - (bp.Protection.ShieldStrength / 100)
	
	finalDamage := impact * armorReduction * shieldReduction
	return bp.health.Reduce(finalDamage)
} 