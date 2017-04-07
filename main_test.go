package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHelp(t *testing.T) {
	goroutinesHelp := "# HELP go_goroutines Number of goroutines that currently exist."

	name, help := parseHelp(goroutinesHelp)
	assert.Equal(t, "go_goroutines", name)
	assert.Equal(t, "Number of goroutines that currently exist.", help)
}

func TestParseType(t *testing.T) {
	goroutinesType := "# TYPE go_goroutines gauge"

	name, types := parseType(goroutinesType)
	assert.Equal(t, "go_goroutines", name)
	assert.Equal(t, "gauge", types)
}

func TestParseMetric(t *testing.T) {
	lines := []string{
		`go_gc_duration_seconds 0.00019936000000000002`,
		`go_gc_duration_seconds{quantile="0"} 0.00019936000000000002`,
		`go_gc_duration_seconds{quantile="0.25"} 0.000281616`,
		`go_gc_duration_seconds{quantile="0.5"} 0.000290125`,
		`go_gc_duration_seconds{quantile="0.75"} 0.000317352`,
		`go_gc_duration_seconds{quantile="1"} 0.004849631`,
	}

	for _, line := range lines {
		assert.Equal(t, "go_gc_duration_seconds", parseMetric(line))
	}
}
