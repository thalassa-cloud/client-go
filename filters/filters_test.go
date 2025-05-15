package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilters_GetLabelFilter(t *testing.T) {
	tests := []struct {
		name     string
		filters  Filters
		expected *LabelFilter
	}{
		{
			name: "when label filter exists",
			filters: Filters{
				&LabelFilter{MatchLabels: map[string]string{"env": "prod"}},
				&FilterKeyValue{Key: FilterRegion, Value: "nl-01"},
			},
			expected: &LabelFilter{MatchLabels: map[string]string{"env": "prod"}},
		},
		{
			name: "when label filter does not exist",
			filters: Filters{
				&FilterKeyValue{Key: FilterRegion, Value: "nl-01"},
				&FilterKeyValue{Key: FilterZone, Value: "zone-1"},
			},
			expected: nil,
		},
		{
			name:     "when filters is empty",
			filters:  Filters{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filters.GetLabelFilter()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilters_GetKeyValueFilter(t *testing.T) {
	tests := []struct {
		name     string
		filters  Filters
		key      FilterKey
		expected *FilterKeyValue
	}{
		{
			name: "when key value filter exists",
			filters: Filters{
				&FilterKeyValue{Key: FilterRegion, Value: "nl-01"},
				&FilterKeyValue{Key: FilterZone, Value: "zone-1"},
			},
			key:      FilterRegion,
			expected: &FilterKeyValue{Key: FilterRegion, Value: "nl-01"},
		},
		{
			name: "when key value filter does not exist",
			filters: Filters{
				&FilterKeyValue{Key: FilterRegion, Value: "nl-01"},
				&FilterKeyValue{Key: FilterZone, Value: "zone-1"},
			},
			key:      FilterVpcIdentity,
			expected: nil,
		},
		{
			name:     "when filters is empty",
			filters:  Filters{},
			key:      FilterRegion,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filters.GetKeyValueFilter(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLabelFilter_ToParams(t *testing.T) {
	tests := []struct {
		name     string
		filter   *LabelFilter
		expected map[string]string
	}{
		{
			name: "with single label",
			filter: &LabelFilter{
				MatchLabels: map[string]string{"env": "prod"},
			},
			expected: map[string]string{
				"matchLabels[env]": "prod",
			},
		},
		{
			name: "with multiple labels",
			filter: &LabelFilter{
				MatchLabels: map[string]string{
					"env":     "prod",
					"region":  "nl-01",
					"version": "1.0",
				},
			},
			expected: map[string]string{
				"matchLabels[env]":     "prod",
				"matchLabels[region]":  "nl-01",
				"matchLabels[version]": "1.0",
			},
		},
		{
			name: "with empty labels",
			filter: &LabelFilter{
				MatchLabels: map[string]string{},
			},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ToParams()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterKeyValue_ToParams(t *testing.T) {
	tests := []struct {
		name     string
		filter   *FilterKeyValue
		expected map[string]string
	}{
		{
			name: "with valid key and value",
			filter: &FilterKeyValue{
				Key:   FilterRegion,
				Value: "nl-01",
			},
			expected: map[string]string{
				"region": "nl-01",
			},
		},
		{
			name: "with empty key",
			filter: &FilterKeyValue{
				Key:   "",
				Value: "nl-01",
			},
			expected: map[string]string{},
		},
		{
			name: "with empty value",
			filter: &FilterKeyValue{
				Key:   FilterRegion,
				Value: "",
			},
			expected: map[string]string{},
		},
		{
			name: "with whitespace key and value",
			filter: &FilterKeyValue{
				Key:   "   ",
				Value: "   ",
			},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ToParams()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterType(t *testing.T) {
	tests := []struct {
		name     string
		filter   Filter
		expected FilterType
	}{
		{
			name:     "label filter type",
			filter:   &LabelFilter{},
			expected: FilterTypeLabel,
		},
		{
			name:     "key value filter type",
			filter:   &FilterKeyValue{},
			expected: FilterTypeKeyValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.FilterType()
			assert.Equal(t, tt.expected, result)
		})
	}
}
