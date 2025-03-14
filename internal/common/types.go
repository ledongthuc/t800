package common

// Location represents 3D coordinates
type Location struct {
	X float64
	Y float64
	Z float64
}

// Threat represents a potential threat to the robot
type Threat struct {
	ID          string
	Type        string
	Location    Location
	Severity    int
	Timestamp   int64
	Description string
}

// OperationMode defines the current operation mode
type OperationMode int

const (
	Normal OperationMode = iota
	Combat
	Emergency
	Maintenance
) 