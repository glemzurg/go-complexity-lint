package common

import "fmt"

// Zone represents a metric zone classification.
type Zone int

const (
	ZoneGreen  Zone = iota // within acceptable limits
	ZoneYellow             // warning zone
	ZoneRed                // failure zone
)

// Thresholds defines the warn and fail boundaries for a metric.
// Values up to and including WarnAt are green.
// Values from WarnAt+1 to FailAt are yellow (warning).
// Values above FailAt are red (failure).
type Thresholds struct {
	WarnAt int
	FailAt int
}

// Validate returns an error if the thresholds are invalid.
// Both values must be non-negative and WarnAt must not exceed FailAt.
func (t Thresholds) Validate(name string) error {
	if t.WarnAt < 0 {
		return fmt.Errorf("%s: warn threshold must be non-negative, got %d", name, t.WarnAt)
	}
	if t.FailAt < 0 {
		return fmt.Errorf("%s: fail threshold must be non-negative, got %d", name, t.FailAt)
	}
	if t.WarnAt > t.FailAt {
		return fmt.Errorf("%s: warn threshold (%d) must not exceed fail threshold (%d)", name, t.WarnAt, t.FailAt)
	}
	return nil
}

// Classify returns the zone for a given metric value.
func (t Thresholds) Classify(value int) Zone {
	switch {
	case value <= t.WarnAt:
		return ZoneGreen
	case value <= t.FailAt:
		return ZoneYellow
	default:
		return ZoneRed
	}
}

// Category returns the analysis.Diagnostic Category string for the zone.
func (z Zone) Category() string {
	switch z {
	case ZoneYellow:
		return "warning"
	case ZoneRed:
		return "error"
	default:
		return ""
	}
}
