package simplequeue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	w := Worker{
		id:          120,
		jobsHandled: 55,
	}

	actual := w.ID()

	assert.Equal(t, int64(120), actual)
}

func TestHandled(t *testing.T) {
	w := Worker{
		id:          120,
		jobsHandled: 55,
	}

	actual := w.Handled()

	assert.Equal(t, int64(55), actual)
}
