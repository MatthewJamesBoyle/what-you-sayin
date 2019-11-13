package reddit_test

import (
	"github.com/matthewjamesboyle/whattheysayingbot/internal/reddit"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSafeResponse(t *testing.T) {
	t.Run("returns a safe response", func(t *testing.T) {
		sr := reddit.NewSafeResponse()
		assert.IsType(t, reddit.SafeResponse{}, *sr)
	})
}

func TestNewConfig(t *testing.T) {
	t.Run("returns an ErrMalformedParam given an empty clientID", func(t *testing.T) {
		c, err := reddit.NewConfig(
			"",
			"",
			"",
			"",
			"",
		)
		require.Nil(t, c)
		assert.Equal(t, reddit.ErrMalformedParam, errors.Cause(err))
		assert.Contains(t, err.Error(), "clientID")
	})

	//just doing a sample of these as it should provide enough evidence it works.
	t.Run("returns an ErrMalformedParam given an empty username", func(t *testing.T) {
		c, err := reddit.NewConfig(
			"someuser",
			"something",
			"something",
			"",
			"",
		)
		require.Nil(t, c)
		assert.Equal(t, reddit.ErrMalformedParam, errors.Cause(err))
		assert.Contains(t, err.Error(), "username")
	})
	t.Run("Returns a config and no error given a valid set of params", func(t *testing.T) {
		c,err := reddit.NewConfig(
			"someuser",
			"something",
			"something",
			"some",
			"some",
		)

		assert.NotNil(t,c)
		assert.NoError(t,err)
	})
}

func TestNewInteractor(t *testing.T) {

}
