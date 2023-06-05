/*
	Copyright (C) 2023  Martijn van der Kleijn

	This file is part of the go-simplequeue library.

    This Source Code Form is subject to the terms of the Mozilla Public
  	License, v. 2.0. If a copy of the MPL was not distributed with this
  	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package simplequeue

// Job should be implemented by any struct intended for queueing
type Job interface {
	ID() int64 // Returns the job's ID
	Do()       // Actually executes the job
}
