package main

import (
	"github.com/loadimpact/speedboat/stats"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseBackendStdout(t *testing.T) {
	output, err := parseBackend("-")
	assert.NoError(t, err)
	assert.IsType(t, &stats.JSONBackend{}, output)
}

func TestParseStagesSimple(t *testing.T) {
	stages, err := parseStages([]string{"10"}, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 10, stages[0].EndVUs)
	assert.Equal(t, 10*time.Second, stages[0].Duration)
}

func TestParseStagesSimpleTrailingDash(t *testing.T) {
	stages, err := parseStages([]string{"10-"}, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 10, stages[0].EndVUs)
	assert.Equal(t, 10*time.Second, stages[0].Duration)
}

func TestParseStagesSimpleRamp(t *testing.T) {
	stages, err := parseStages([]string{"10-15"}, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 15, stages[0].EndVUs)
	assert.Equal(t, 10*time.Second, stages[0].Duration)
}

func TestParseStagesSimpleRampZeroBackref(t *testing.T) {
	stages, err := parseStages([]string{"-15"}, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(stages))
	assert.Equal(t, 0, stages[0].StartVUs)
	assert.Equal(t, 15, stages[0].EndVUs)
	assert.Equal(t, 10*time.Second, stages[0].Duration)
}

func TestParseStagesSimpleMulti(t *testing.T) {
	stages, err := parseStages([]string{"10", "15"}, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 10, stages[0].EndVUs)
	assert.Equal(t, 5*time.Second, stages[0].Duration)
	assert.Equal(t, 15, stages[1].StartVUs)
	assert.Equal(t, 15, stages[1].EndVUs)
	assert.Equal(t, 5*time.Second, stages[1].Duration)
}

func TestParseStagesSimpleMultiRamp(t *testing.T) {
	stages, err := parseStages([]string{"10-15", "15-20"}, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 15, stages[0].EndVUs)
	assert.Equal(t, 5*time.Second, stages[0].Duration)
	assert.Equal(t, 15, stages[1].StartVUs)
	assert.Equal(t, 20, stages[1].EndVUs)
	assert.Equal(t, 5*time.Second, stages[1].Duration)
}

func TestParseStagesSimpleMultiRampBackref(t *testing.T) {
	stages, err := parseStages([]string{"10-15", "-20"}, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 15, stages[0].EndVUs)
	assert.Equal(t, 5*time.Second, stages[0].Duration)
	assert.Equal(t, 15, stages[1].StartVUs)
	assert.Equal(t, 20, stages[1].EndVUs)
	assert.Equal(t, 5*time.Second, stages[1].Duration)
}

func TestParseStagesFixed(t *testing.T) {
	stages, err := parseStages([]string{"10:15s"}, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 10, stages[0].EndVUs)
	assert.Equal(t, 15*time.Second, stages[0].Duration)
}

func TestParseStagesFixedFluid(t *testing.T) {
	stages, err := parseStages([]string{"10:5s", "15"}, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 10, stages[0].EndVUs)
	assert.Equal(t, 5*time.Second, stages[0].Duration)
	assert.Equal(t, 15, stages[1].StartVUs)
	assert.Equal(t, 15, stages[1].EndVUs)
	assert.Equal(t, 5*time.Second, stages[1].Duration)
}

func TestParseStagesFixedFluidNoTimeLeft(t *testing.T) {
	stages, err := parseStages([]string{"10:10s", "15"}, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 10, stages[0].EndVUs)
	assert.Equal(t, 10*time.Second, stages[0].Duration)
	assert.Equal(t, 15, stages[1].StartVUs)
	assert.Equal(t, 15, stages[1].EndVUs)
	assert.Equal(t, 0*time.Second, stages[1].Duration)
}

func TestParseStagesInvalid(t *testing.T) {
	_, err := parseStages([]string{"a"}, 10*time.Second)
	assert.Error(t, err)
}

func TestParseStagesInvalidStart(t *testing.T) {
	_, err := parseStages([]string{"a-15"}, 10*time.Second)
	assert.Error(t, err)
}

func TestParseStagesInvalidEnd(t *testing.T) {
	_, err := parseStages([]string{"15-a"}, 10*time.Second)
	assert.Error(t, err)
}

func TestParseStagesInvalidTime(t *testing.T) {
	_, err := parseStages([]string{"15:a"}, 10*time.Second)
	assert.Error(t, err)
}

func TestParseStagesInvalidTimeMissingUnit(t *testing.T) {
	_, err := parseStages([]string{"15:10"}, 10*time.Second)
	assert.Error(t, err)
}

func TestGuessTypeURL(t *testing.T) {
	assert.Equal(t, typeURL, guessType("http://example.com/"))
}

func TestGuessTypeJS(t *testing.T) {
	assert.Equal(t, typeJS, guessType("script.js"))
}

func TestGuessTypeUnknown(t *testing.T) {
	assert.Equal(t, "", guessType("script.txt"))
}
