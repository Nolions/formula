package log_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/redpkg/formula/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	assert := assert.New(t)

	err := log.Init(log.Config{
		Level:   "trace",
		Console: true,
	})
	if !assert.NoError(err) {
		return
	}

	log.Trace().Msg("trace message")
	log.Debug().Msg("debug message")
	log.Info().Msg("info message")
	log.Warn().Msg("warn message")
	log.Error().Msg("error message")
	log.Trace().Stack().Err(errors.New("foobar")).Msg("err message")
	log.Error().Stack().Err(outer()).Msg("err message")
	log.Debug().Err(errors.New("foobar")).Msg("err message")
	log.Error().Err(outer()).Msg("err message")
	// log.Fatal().Msg("fatal")
	// log.Panic().Msg("panic")
}

func inner() error {
	return errors.New("seems we have an error here")
}

func middle() error {
	err := inner()
	if err != nil {
		return err
	}
	return nil
}

func outer() error {
	err := middle()
	if err != nil {
		return err
	}
	return nil
}
