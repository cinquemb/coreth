// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package aclock contains an adjustable clock implementation. The time can be
// adjusted by adding an offset.
package aclock

import (
	"errors"
	"time"
)

var offset time.Duration

var (
	// MaxTime was taken from https://stackoverflow.com/questions/25065055/what-is-the-maximum-time-time-in-go/32620397#32620397
	MaxTime = time.Unix(1<<63-62135596801, 0) // 0 is used because we drop the nano-seconds
)

// Clock acts as a thin wrapper around global time that allows for easy testing
type Clock struct {
	faked bool
	time  time.Time
}

// Set the time on the clock
func (c *Clock) Set(time time.Time) { c.faked = true; c.time = time }

// Sync this clock with global time
func (c *Clock) Sync() { c.faked = false }

// Time returns the time on this clock
func (c *Clock) Time() time.Time {
	if c.faked {
		return c.time
	}
	return time.Now().Add(offset)
}

// Unix returns the unix time on this clock.
func (c *Clock) Unix() uint64 {
	unix := c.Time().Unix()
	if unix < 0 {
		unix = 0
	}
	return uint64(unix)
}


// AddOffset adds offset d to the current offset. d cannot be negative, and an
// error will be returned if the resulting offset overflows time.Duration.
func AddOffset(d time.Duration) (time.Duration, error) {
	if d < 0 {
		return 0, errors.New("aclock: duration for offset cannot be negative")
	}
	if offset+d < offset {
		return 0, errors.New("aclock: offset overflow")
	}
	offset += d
	return offset, nil
}

// Now returns the current time with the offset applied.
func Now() time.Time {
	return time.Now().Add(offset)
}

// NowWithOffset returns the current time with the offset applied and the offset
// itself.
func NowWithOffset() (time.Time, time.Duration) {
	return time.Now().Add(offset), offset
}