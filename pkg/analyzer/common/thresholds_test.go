package common

import "testing"

func TestClassify(t *testing.T) {
	th := Thresholds{WarnAt: 4, FailAt: 6}

	tests := []struct {
		value int
		want  Zone
	}{
		{0, ZoneGreen},
		{4, ZoneGreen},
		{5, ZoneYellow},
		{6, ZoneYellow},
		{7, ZoneRed},
		{100, ZoneRed},
	}

	for _, tt := range tests {
		got := th.Classify(tt.value)
		if got != tt.want {
			t.Errorf("Classify(%d) = %d, want %d", tt.value, got, tt.want)
		}
	}
}

func TestZoneCategory(t *testing.T) {
	tests := []struct {
		zone Zone
		want string
	}{
		{ZoneGreen, ""},
		{ZoneYellow, "warning"},
		{ZoneRed, "error"},
	}

	for _, tt := range tests {
		got := tt.zone.Category()
		if got != tt.want {
			t.Errorf("Zone(%d).Category() = %q, want %q", tt.zone, got, tt.want)
		}
	}
}
