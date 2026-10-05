package policy

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestStrategyValid(t *testing.T) {
	tests := []struct {
		name     string
		strategy Strategy
		want     bool
	}{
		{name: "keep-source is valid", strategy: KeepSource, want: true},
		{name: "keep-target is valid", strategy: KeepTarget, want: true},
		{name: "manual is valid", strategy: Manual, want: true},
		{name: "fast-forward is no longer valid", strategy: Strategy("fast-forward"), want: false},
		{name: "empty string is invalid", strategy: Strategy(""), want: false},
		{name: "typo is invalid", strategy: Strategy("keep-sourcee"), want: false},
		{name: "unrelated word is invalid", strategy: Strategy("banana"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.strategy.Valid(); got != tt.want {
				t.Errorf("Strategy(%q).Valid() = %v; want %v", tt.strategy, got, tt.want)
			}
		})
	}
}

func TestStrategyUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name      string
		yamlValue string
		want      Strategy
		wantError bool
	}{
		{name: "valid keep-target", yamlValue: "keep-target", want: KeepTarget},
		{name: "valid keep-source", yamlValue: "keep-source", want: KeepSource},
		{name: "valid manual", yamlValue: "manual", want: Manual},
		{name: "fast-forward is rejected", yamlValue: "fast-forward", wantError: true},
		{name: "invalid value fails", yamlValue: "keep-sourcee", wantError: true},
		{name: "empty string value fails", yamlValue: `""`, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Strategy
			err := yaml.Unmarshal([]byte(tt.yamlValue), &got)
			if tt.wantError {
				if err == nil {
					t.Errorf("yaml.Unmarshal(%q) expected error, got nil", tt.yamlValue)
				}
				return
			}
			if err != nil {
				t.Errorf("yaml.Unmarshal(%q) unexpected error: %v", tt.yamlValue, err)
				return
			}
			if got != tt.want {
				t.Errorf("yaml.Unmarshal(%q) = %v; want %v", tt.yamlValue, got, tt.want)
			}
		})
	}
}
