package log_test

import (
	"errors"
	"testing"

	"github.com/redpkg/formula/log"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	err := log.Init(newConfig())
	if !assert.NoError(err) {
		return
	}

	log.Debug().Msg("debug message")
	log.Info().Msg("info message")
	log.Warn().Msg("warn message")
	log.Error().Msg("error message")
	log.Err(errors.New("foo")).Msg("error message")
	// Fatal().Msg("fatal message")
	// Panic().Msg("panic message")
}

func newConfig() log.Config {
	return log.Config{
		Level:   "debug",
		Console: true,
	}
}
