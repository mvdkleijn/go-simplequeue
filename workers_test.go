/*
	Copyright (C) 2023  Martijn van der Kleijn

	This file is part of the go-simplequeue library.

    This Source Code Form is subject to the terms of the Mozilla Public
  	License, v. 2.0. If a copy of the MPL was not distributed with this
  	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

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
