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
	"testing"
	"time"
)

func TestSpc1Init(t *testing.T) {
	asu1, asu2 := uint32(45), uint32(45)
	asu3 := uint32(10)
	Spc1Init(50, 1, asu1, asu2, asu3)
}

func TestNewSpc1Io(t *testing.T) {
	asu1, asu2 := uint32(45), uint32(45)
	asu3 := uint32(10)
	Spc1Init(50, 1, asu1, asu2, asu3)

	s := NewSpc1Io(1)
	s.Generate()
	if s.Asu < 1 || s.Asu > 3 {
		t.Errorf("Illegal value of s.Asu: %d\n", s.Asu)
	}
	if s.Stream < 1 || s.Stream > 7 {
		t.Errorf("Illegal value of s.Stream: %d\n", s.Stream)
	}
	if s.Offset >= 45 {
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

		/*
			fmt.Printf("%d:asu=%v:"+
				"rw=%v:"+
				"blocks=%v:"+
				"stream=%v:"+
				"offset=%v:"+
				"when=%v\n",
				i,
				s.Asu,
				s.Isread,
				s.Blocks,
				s.Stream,
				s.Offset,
				s.When)
		*/
	}
	end := time.Now()
	iops := float64(ios) / end.Sub(start).Seconds()
	if iops < 2450 || iops > 2550 {
		t.Errorf("Incorrect number of Iops: %.2f\n", iops)
	}
}

/*
var ios_sent uint64
var blocks uint64

func workload(wg *sync.WaitGroup, context, ios int) {
	var s C.spc1_ios_t

	defer wg.Done()

	start := time.Now()
	for i := 0; i < ios; i++ {
		C.spc1_next_op(&s, C.int(context))
		a := time.Millisecond / 10 * time.Duration(s.when)
		b := time.Now().Sub(start)
		sleep_time := a - b
		//fmt.Printf("a=%v:b=%v:", a, b)
		if sleep_time > 0 {
			time.Sleep(sleep_time)
			//fmt.Printf("%v\n", sleep_time)
		} else {
			//fmt.Print("_\n")
		}
		atomic.AddUint64(&ios_sent, uint64(1))
		atomic.AddUint64(&blocks, uint64(s.len))

		// send the io
		time.Sleep(time.Millisecond * 20)
			fmt.Printf("%d:%d:asu=%v:"+
				"rw=%v:"+
				"len=%v:"+
				"stream=%v:"+
				"bsu=%v:"+
				"offset=%v:"+
				"when=%v\n",
				i, context,
				s.asu,
				s.dir,
				s.len,
				s.stream,
				s.bsu,
				s.pos,
				s.when)
	}
}

func main() {
	var wg sync.WaitGroup

	contexts := 10
	C.spc1_init(C.CString("test"),
		50,              // bsu
		45*1024*1024/4,  // 45G as 4k blocks
		45*1024*1024/4,  // 45G as 4k blocks
		10*1024*1024/4,  // 10G as 4k blocks
		C.int(contexts), // contexts
		nil,             // version
		0)

	ios := 1000
	start := time.Now()
	for context := 0; context < contexts; context++ {
		wg.Add(1)
		go workload(&wg, context, ios)
	}
	wg.Wait()
	end := time.Now()
	fmt.Printf("IOPS = %v\n",
		float64(ios_sent)/end.Sub(start).Seconds())
	fmt.Printf("Bw = %.2f MB/s\n",
		float64(blocks*4*1024)/end.Sub(start).Seconds())
}
*/
