package common

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
