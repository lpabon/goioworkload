//
// Copyright (c) 2014 The goioworkload Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package spc1

// #cgo LDFLAGS: -lm
// #include "spc1.h"
import "C"

import (
	"fmt"
	"time"
)

type Spc1Io struct {
	// Which asu?
	Asu uint32

	// Read or write
	Isread bool

	// Length of transfer in units of 4KB
	Blocks uint32

	// Which stream in the bsu
	Stream uint32

	// Offset in units of 4KB
	Offset uint32

	// When to do this I/O from the start of the run
	When time.Duration

	spc1_ios C.spc1_ios_t
	context  int
}

// bsus: Number of BSUs
// contexts: Number of contexts
// asuXsize: Size in 4k blocks
func Spc1Init(bsus, contexts int,
	asu1size, asu2size, asu3size uint32) {

	C.spc1_init(C.CString("gospc1"),
		C.int(bsus),
		C.uint(asu1size),
		C.uint(asu2size),
		C.uint(asu3size),
		C.int(contexts),
		nil,
		0)
}

// Must have called Spc1Init() to initalize
// the workload generator
func NewSpc1Io(context int) *Spc1Io {
	return &Spc1Io{
		context: context,
	}
}

func (s *Spc1Io) Generate() {
	C.spc1_next_op(&s.spc1_ios, C.int(s.context))
	s.Asu = uint32(s.spc1_ios.asu)
	s.Blocks = uint32(s.spc1_ios.len)
	s.Isread = s.spc1_ios.dir == 0
	s.Stream = uint32(s.spc1_ios.stream)
	s.Offset = uint32(s.spc1_ios.pos)
	s.When = time.Millisecond / 10 * time.Duration(s.spc1_ios.when)
}

func (s *Spc1Io) String() string {
	return fmt.Sprintf("asu=%v:"+
		"rw=%v:"+
		"blocks=%v:"+
		"stream=%v:"+
		"offset=%v:"+
		"when=%v\n",
		s.Asu,
		s.Isread,
		s.Blocks,
		s.Stream,
		s.Offset,
		s.When)

}
