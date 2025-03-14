package anatomy

import "fmt"

// Dimensions represents physical dimensions of a body part
type Dimensions struct {
	Width  float64 // in meters
	Height float64 // in meters
	Depth  float64 // in meters
	Weight float64 // in kilograms
}

// NewDimensions creates a new Dimensions instance with validation
func NewDimensions(width, height, depth, weight float64) (*Dimensions, error) {
	if width <= 0 || height <= 0 || depth <= 0 || weight <= 0 {
		return nil, fmt.Errorf("invalid dimensions: all values must be positive")
	}

	return &Dimensions{
		Width:  width,
		Height: height,
		Depth:  depth,
		Weight: weight,
	}, nil
}

// Volume returns the volume in cubic meters
func (d *Dimensions) Volume() float64 {
	return d.Width * d.Height * d.Depth
}

// SurfaceArea returns the surface area in square meters
func (d *Dimensions) SurfaceArea() float64 {
	return 2 * (d.Width*d.Height + d.Height*d.Depth + d.Depth*d.Width)
}

// Density returns the density in kg/m³
func (d *Dimensions) Density() float64 {
	return d.Weight / d.Volume()
}

// Scale returns a new Dimensions instance scaled by the given factor
func (d *Dimensions) Scale(factor float64) *Dimensions {
	return &Dimensions{
		Width:  d.Width * factor,
		Height: d.Height * factor,
		Depth:  d.Depth * factor,
		Weight: d.Weight * factor,
	}
} 