/*
 * This file from
 *  https://github.com/addthis/stream-lib/blob/master/src/main/java/com/clearspring/analytics/stream/cardinality/RegisterSet.java
 *
 *  This class modified by Scouter-Project *   - original package :  com.clearspring.analytics.stream.cardinality
 *
 *  ====================================
 *
 * Copyright (C) 2012 Clearspring Technologies, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hll

import (
//	"math"
//	"strconv"
)

const (
	LOG2_BITS_PER_WORD = 6
	REGISTER_SIZE      = 5
)

type RegisterSet struct {
	Count int
	Size  int
	M     []uint32
}

func NewRegisterSetInit(count int, initialValues []uint32) *RegisterSet {
	p := new(RegisterSet)
	p.Count = count
	if initialValues == nil {
		p.M = make([]uint32, getSizeForCount(count))
	} else {
		p.M = initialValues
	}
	p.Size = len(p.M)

	return p
}

func NewRegisterSet(count int) *RegisterSet {
	return NewRegisterSetInit(count, nil)
}

func (this *RegisterSet) Set(position, value uint32) {
	bucketPos := position / LOG2_BITS_PER_WORD
	shift := REGISTER_SIZE * (position - (bucketPos * LOG2_BITS_PER_WORD))
	//this.M[bucketPos] = (this.M[bucketPos] & ~(0x1f << shift)) | (value << shift)
	this.M[bucketPos] = (this.M[bucketPos] & ^ (0x1f << uint32(shift)) | (value << uint32(shift)))
}
func (this *RegisterSet) Get(position int) uint32 {
	bucketPos := position / LOG2_BITS_PER_WORD
	shift := REGISTER_SIZE * (position - (bucketPos * LOG2_BITS_PER_WORD))
	//return (this.M[bucketPos] & (0x1f << shift)) >>> shift;
	return (this.M[bucketPos]&(0x1f<<uint32(shift))) >> uint32(shift)
}

func (this *RegisterSet) UpdateIfGreater(position, value uint32) bool {
	bucket := position / LOG2_BITS_PER_WORD
	shift := REGISTER_SIZE * (position - (bucket * LOG2_BITS_PER_WORD))
	mask := uint32(0x1f) << uint32(shift)
	// Use long to avoid sign issues with the left-most shift
	//long curVal = this.M[bucket] & mask;
	//long newVal = value << shift;
	curVal := uint64(this.M[bucket] & mask)
	newVal := uint64(value) << uint32(shift)
	if curVal < newVal {
		//this.M[bucket] = (int) ((this.M[bucket] & ~mask) | newVal)
		this.M[bucket] = uint32(uint64(this.M[bucket] & ^mask) | newVal)
		return true
	} else {
		return false
	}
}
func (this *RegisterSet) Merge(that *RegisterSet) {
	for bucket := 0; bucket < len(this.M); bucket++ {
		word := uint32(0)
		for j := 0; j < LOG2_BITS_PER_WORD; j++ {
			mask := uint32(0x1f << uint32(REGISTER_SIZE * j))
			thisVal := (this.M[bucket] & mask)
			thatVal := (that.M[bucket] & mask)
			if thisVal < thatVal {
				word |= thatVal
			} else {
				word |= thisVal
			}
		}
		this.M[bucket] = word
	}
}
func (this *RegisterSet) ReadOnlyBits() []uint32 {
	return this.M
}

func (this *RegisterSet) Bits() []uint32 {
	copy := this.M[0:]
	return copy
}

func getBits(count int) int {
	return count / LOG2_BITS_PER_WORD
}

func getSizeForCount(count int) int {
	bits := getBits(count)
	if bits == 0 {
		return 1
	} else if bits%32 == 0 {
		return bits
	} else {
		return bits + 1
	}
}
