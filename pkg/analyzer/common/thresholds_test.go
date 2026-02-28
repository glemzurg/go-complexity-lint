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

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		th      Thresholds
		wantErr bool
	}{
		{
			name: "valid thresholds",
			th:   Thresholds{WarnAt: 4, FailAt: 6},
		},
		{
			name: "equal thresholds",
			th:   Thresholds{WarnAt: 5, FailAt: 5},
		},
		{
			name: "zero thresholds",
			th:   Thresholds{WarnAt: 0, FailAt: 0},
		},
		{
			name:    "negative warn",
			th:      Thresholds{WarnAt: -1, FailAt: 6},
			wantErr: true,
		},
		{
			name:    "negative fail",
			th:      Thresholds{WarnAt: 4, FailAt: -1},
			wantErr: true,
		},
		{
			name:    "warn exceeds fail",
			th:      Thresholds{WarnAt: 10, FailAt: 5},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.th.Validate("test")
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
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
