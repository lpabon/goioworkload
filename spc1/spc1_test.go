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

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type Meter struct {
	reads, writes uint64
}

var (
	asu1m, asu2m Meter
)

func TestSpc1Init(t *testing.T) {
	asu1, asu2 := uint32(4500), uint32(4500)
	asu3 := uint32(1000)
	Spc1Init(50, 1, asu1, asu2, asu3)
}

func TestNewSpc1Io(t *testing.T) {
	asu1, asu2 := uint32(4500), uint32(4500)
	asu3 := uint32(1000)
	Spc1Init(50, 1, asu1, asu2, asu3)

	s := NewSpc1Io(1)
	s.Generate()
	if s.Asu < 1 || s.Asu > 3 {
		t.Errorf("Illegal value of s.Asu: %d\n", s.Asu)
	}
	if s.Stream < 0 || s.Stream > 7 {
		t.Errorf("Illegal value of s.Stream: %d\n", s.Stream)
	}
	if s.Offset >= 4500 {
		t.Errorf("Offset out of bounds: %d\n", s.Offset)
	}
}

func TestSpc1Generate(t *testing.T) {
	asu1, asu2 := uint32(45*1024*1024/4), uint32(45*1024*1024/4)
	asu3 := uint32(10 * 1024 * 1024 / 4)

	// 50 BSUs, each BSU doing 50 Iops
	// Total IOPs should be ~2500
	Spc1Init(50, 1, asu1, asu2, asu3)

	s := NewSpc1Io(1)
	ios := 10000
	start := time.Now()
	lastiotime := start
	for i := 0; i < ios; i++ {
		s.Generate()
		sleep_time := start.Add(s.When).Sub(lastiotime)
		if sleep_time > 0 {
			time.Sleep(sleep_time)
		}
		lastiotime = time.Now()
	}
	end := time.Now()
	iops := float64(ios) / end.Sub(start).Seconds()
	if iops < 2450 || iops > 2550 {
		t.Errorf("Incorrect number of Iops: %.2f\n", iops)
	}
}

func TestSpc1Contexts(t *testing.T) {
	asu1, asu2 := uint32(4500), uint32(4500)
	asu3 := uint32(1000)
	contexts := 4

	// 50 BSUs, each BSU doing 50 Iops
	// Total IOPs should be ~2500
	Spc1Init(50, contexts, asu1, asu2, asu3)

	for context := 1; context <= contexts; context++ {
		s := NewSpc1Io(context)
		s.Generate()
		fmt.Print(s)
		if s.Asu < 1 || s.Asu > 3 {
			t.Errorf("Illegal value of s.Asu: %d\n", s.Asu)
		}
		if s.Stream < 0 || s.Stream > 7 {
			t.Errorf("Illegal value of s.Stream: %d\n", s.Stream)
		}
		if s.Offset >= 4500 {
			t.Errorf("Offset out of bounds: %d\n", s.Offset)
		}
	}
}

func context_tester(wg *sync.WaitGroup, context int, t *testing.T) {
	defer wg.Done()

	start := time.Now()
	lastiotime := start
	for io := 1; io < 10000; io++ {
		s := NewSpc1Io(context)
		err := s.Generate()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		if s.Asu < 1 || s.Asu > 3 {
			t.Errorf("Illegal value of s.Asu: %d\n", s.Asu)
			t.FailNow()
		}
		if s.Stream < 0 || s.Stream > 7 {
			t.Errorf("Illegal value of s.Stream: %d\n", s.Stream)
			t.FailNow()
		}

		if s.Asu == 1 {
			if s.Isread {
				atomic.AddUint64(&asu1m.reads, 1)
			} else {
				atomic.AddUint64(&asu1m.writes, 1)
			}
		} else if s.Asu == 2 {
			if s.Isread {
				atomic.AddUint64(&asu2m.reads, 1)
			} else {
				atomic.AddUint64(&asu2m.writes, 1)
			}
		}
		// Check how much time we should wait
		sleep_time := start.Add(s.When).Sub(lastiotime)
		if sleep_time > 0 {
			time.Sleep(sleep_time)
		}
		lastiotime = time.Now()
	}
}

func TestSpc1ConcurrentContexts(t *testing.T) {
	asu1, asu2 := uint32(45000), uint32(45000)
	asu3 := uint32(10000)
	bsu := 200
	contexts := int((bsu + 99) / 100)

	// 200 BSUs, each BSU doing 50 Iops
	// Total IOPs should be ~5k
	err := Spc1Init(bsu, contexts, asu1, asu2, asu3)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var wg sync.WaitGroup
	start := time.Now()
	for context := 0; context < contexts; context++ {
		wg.Add(1)
		go context_tester(&wg, context, t)
	}
	wg.Wait()

	end := time.Now()

	iops := int(float64(10000*contexts) / end.Sub(start).Seconds())

	if iops < 9500 || iops > 10500 {
		t.Errorf("Incorrect number of iops")
	}

	fmt.Printf("ASU1 Read Rate = %.4f\n"+
		"ASU2 Read Rate = %.4f\n",
		float64(asu1m.reads)/float64(asu1m.reads+asu1m.writes),
		float64(asu2m.reads)/float64(asu2m.reads+asu2m.writes))
}
