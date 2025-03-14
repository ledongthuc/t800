package common

import "math"

// Location represents 3D coordinates
type Location struct {
	X float64
	Y float64
	Z float64
}

// MovementSpeed represents the robot's movement capabilities
type MovementSpeed struct {
	Linear  float64 // meters per second
	Angular float64 // radians per second
}

// DefaultSpeed returns the default movement speed configuration
func DefaultSpeed() MovementSpeed {
	return MovementSpeed{
		Linear:  5.0,  // 5 meters per second
		Angular: math.Pi / 2, // 90 degrees per second
	}
}

// MoveTowards calculates new position when moving towards a target
func (loc *Location) MoveTowards(target Location, speed float64, deltaTime float64) Location {
	direction := Location{
		X: target.X - loc.X,
		Y: target.Y - loc.Y,
		Z: target.Z - loc.Z,
	}
	
	distance := CalculateDistance(*loc, target)
	if distance == 0 {
		return *loc
	}
	
	// Normalize direction and multiply by speed and time
	moveDistance := speed * deltaTime
	if moveDistance > distance {
		moveDistance = distance
	}
	
	scale := moveDistance / distance
	return Location{
		X: loc.X + direction.X * scale,
		Y: loc.Y + direction.Y * scale,
		Z: loc.Z + direction.Z * scale,
	}
}

// Threat represents a potential threat to the robot
type Threat struct {
	ID          string
	Type        string
	Location    Location
	Severity    int
	Timestamp   int64
	Description string
	Health      float64 // Health percentage (0-100)
}

// OperationMode defines the current operation mode
type OperationMode int

const (
	Normal OperationMode = iota
	Combat
	Emergency
	Maintenance
)

// CalculateDistance computes the Euclidean distance between two locations
func CalculateDistance(loc1, loc2 Location) float64 {
	return math.Sqrt(
		math.Pow(loc1.X-loc2.X, 2) +
			math.Pow(loc1.Y-loc2.Y, 2) +
			math.Pow(loc1.Z-loc2.Z, 2),
	)
} 