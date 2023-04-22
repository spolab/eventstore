package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"toremo.com/petclinic/eventstore/pkg/errors"
)

func TestInvalidStreamVersionError(t *testing.T) {
	assert.Error(t, errors.InvalidStreamVersionError())
}
