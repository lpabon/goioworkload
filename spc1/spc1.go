//
// Copyright (c) 2014 The zipfworkload Authors
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
package main

// #cgo LDFLAGS: -lm
// #include "spc1.h"
import "C"

import (
	"fmt"
)

func workload(context, ios int) {
	var s C.ssss

	for i:=0; i<ios; i++ {
		C.spc1_next_op(&s, C.int(context))
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
	contexts := 1
	C.spc1_init(C.CString("test"),
		100, // bsu
		45*1024*1024 / 4, // 45G as 4k blocks
		45*1024*1024 / 4, // 45G as 4k blocks
		10*1024*1024 /4, // 10G as 4k blocks
		C.int(contexts), // contexts
		nil, // version
		0)

	ios := 1000
	workload(1, ios)
}