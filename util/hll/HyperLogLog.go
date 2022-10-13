/*
 * This file from
 *   https://github.com/addthis/stream-lib/blob/master/src/main/java/com/clearspring/analytics/stream/cardinality/HyperLogLog.java
 *
 *   This class modified by Scouter-Project *   - original package :  com.clearspring.analytics.stream.cardinality
 *   - remove implements : ICardinality, Serializable
 *   - add method : public boolean offer(long o)
 *   - remove classes : Builder,  enum Format, HyperLogLogPlusMergeException, SerializationHolder
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

/**
 * Java implementation of HyperLogLog (HLL) algorithm from this paper:
 * <p/>
 * http://algo.inria.fr/flajolet/Publications/FlFuGaMe07.pdf
 * <p/>
 * HLL is an improved version of LogLog that is capable of estimating the
 * cardinality of a set with accuracy = 1.04/sqrt(m) where m = 2^b. So we can
 * control accuracy vs space usage by increasing or decreasing b.
 * <p/>
 * The main benefit of using HLL over LL is that it only requires 64% of the
 * space that LL does to get the same accuracy.
 * <p/>
 * This implementation implements a single counter. If a large (millions) number
 * of counters are required you may want to refer to:
 * <p/>
 * http://dsiutils.dsi.unimi.it/
 * <p/>
 * It has a more complex implementation of HLL that supports multiple counters
 * in a single object, drastically reducing the java overhead from creating a
 * large number of objects.
 * <p/>
 * This implementation leveraged a javascript implementation that Yammer has
 * been working on:
 * <p/>
 * https://github.com/yammer/probablyjs
 * <p>
 * Note that this implementation does not include the long range correction
 * function defined in the original paper. Empirical evidence shows that the
 * correction function causes more harm than good.
 * </p>
 * <p/>
 * <p>
 * Users have different motivations to use different types of hashing functions.
 * Rather than try to keep up with all available hash functions and to remove
 * the concern of causing future binary incompatibilities this class allows
 * clients to offer the value in hashed int or long form. This way clients are
 * free to change their hash function on their own time line. We recommend using
 * Google's Guava Murmur3_128 implementation as it provides good performance and
 * speed when high precision is required. In our tests the 32bit MurmurHash
 * function included in this pcode is faster and produces better results than
 * the 32 bit murmur3 implementation google provides.
 * </p>
 *
 */

package hll

import (
	//"log"
	"math"

	"github.com/whatap/golib/io"
)

type HyperLogLog struct {
	dirty       bool
	registerSet *RegisterSet
	log2m       uint32
	alphaMM     float64
}

/**
 * Creates a new HyperLogLog instance using the given registers. Used for
 * unmarshalling a serialized instance and for merging multiple counters
 * together.
 *
 * @param registerSet
 *            - the initial values for the register set
 */
//	@Deprecated
func NewHyperLogLog(log2m uint32, registerSet *RegisterSet) *HyperLogLog {
	p := new(HyperLogLog)
	if !validateLog2m(log2m) {
		return nil
	}
	p.registerSet = registerSet
	p.log2m = log2m
	m := uint32(1) << p.log2m
	p.alphaMM = p.getAlphaMM(log2m, m)

	return p
}

/**
 * Create a new HyperLogLog instance. The log2m parameter defines the
 * accuracy of the counter. The larger the log2m the better the accuracy.
 * <p/>
 * accuracy = 1.04/sqrt(2^log2m)
 *
 * @param log2m
 *            - the number of bits to use as the basis for the HLL instance
 */
func NewHyperLogLogInt(log2m uint32) *HyperLogLog {
	p := NewHyperLogLog(log2m, NewRegisterSet(1<<log2m))
	return p
}

/**
 * Create a new HyperLogLog instance using the specified standard deviation.
 *
 * @param rsd
 *            - the relative standard deviation for the counter. smaller
 *            values create counters that require more space.
 */
func NewHyperLogLogFloat(Rsd float64) *HyperLogLog {
	//this(log2m(rsd));
	p := NewHyperLogLogInt(_log2m(Rsd))

	return p
}

func NewHyperLogLogDefault() *HyperLogLog {
	p := NewHyperLogLogInt(10)

	return p
}

func (this *HyperLogLog) offerHashedLong(hashedValue uint64) bool {
	// j becomes the binary address determined by the first b log2m of x
	// j will be between 0 and 2^log2m
	//j := int(uint64(hashedValue) >> (strconv. Long.SIZE - log2m))
	j := uint32(hashedValue >> (64 - this.log2m))

	// TODO numberOfLeadingZeros
	//r := int32(1) //Long.numberOfLeadingZeros((hashedValue << this.log2m) | (1 << (this.log2m - 1)) + 1) + 1;
	r := uint32(clz64((hashedValue<<this.log2m)|(1<<(this.log2m-1))+1)) + 1

	return this.registerSet.UpdateIfGreater(j, r)
}
func (this *HyperLogLog) offerHashed(hashedValue uint32) bool {
	// j becomes the binary address determined by the first b log2m of x
	// j will be between 0 and 2^log2m
	//final int j = hashedValue >>> (Integer.SIZE - log2m);
	j := hashedValue >> (32 - this.log2m)

	// TODO numberOfLeadingZeros
	//final int r = Integer.numberOfLeadingZeros((hashedValue << this.log2m) | (1 << (this.log2m - 1)) + 1) + 1;
	//r := int32(1) //Integer.numberOfLeadingZeros((hashedValue << this.log2m) | (1 << (this.log2m - 1)) + 1) + 1
	r := uint32(clz32((hashedValue<<this.log2m)|(1<<(this.log2m-1))+1) + 1)

	//fmt.Println("HyperLogLog offerHashed hash=", hashedValue , ",j=", j, ",r=", r)

	return this.registerSet.UpdateIfGreater(j, r)
}

//	public boolean offer(Object o) {
//			final int x = MurmurHash.hash(o);
//			return offerHashed(x);
//		}
func (this *HyperLogLog) Offer(o uint32) bool {
	x := MurmurHash(o)
	return this.offerHashed(x)
}

func (this *HyperLogLog) OfferLong(o uint64) bool {
	x := MurmurHashLong(o)
	return this.offerHashed(x)
}

func (this *HyperLogLog) Cardinality() uint64 {
	registerSum := float64(0)
	count := this.registerSet.Count
	zeros := float64(0.0)
	for j := 0; j < this.registerSet.Count; j++ {
		val := this.registerSet.Get(j)
		registerSum += float64(1.0) / float64(int(1)<<uint(val))
		if val == 0 {
			zeros++
		}
	}
	estimate := this.alphaMM * (float64(1) / registerSum)
	if estimate <= (float64(5.0)/float64(2.0))*float64(count) {
		// Small Range Estimate
		return uint64(Round(linearCounting(count, zeros)))
	} else {
		return uint64(Round(estimate))
	}
}
func (this *HyperLogLog) Sizeof() int {
	return this.registerSet.Size * 4
}

/*
 * This method is modified by Souter-pcode
 *
 */
func (this *HyperLogLog) GetBytes() []byte {
	out := io.NewDataOutputX()
	out.WriteInt(int32(this.log2m))
	out.WriteInt(int32(this.registerSet.Size))
	for _, x := range this.registerSet.ReadOnlyBits() {
		out.WriteInt(int32(x))
	}
	return out.ToByteArray()
}

/**
 * Add all the elements of the other set to this set.
 * <p/>
 * This operation does not imply a loss of precision.
 *
 * @param other
 *            A compatible Hyperloglog instance (same log2m)
 * @throws CardinalityMergeException
 *             if other is not compatible
 */
func (this *HyperLogLog) AddAll(other *HyperLogLog) {

	if this.Sizeof() != other.Sizeof() {
		//throw new RuntimeException("Cannot merge estimators of different sizes");
		panic("AddAll Cannot merge estimators of different sizes")
	}

	this.registerSet.Merge(other.registerSet)
}

func (this *HyperLogLog) Merge(estimators ...*HyperLogLog) *HyperLogLog {
	merged := NewHyperLogLog(this.log2m, NewRegisterSet(this.registerSet.Count))
	merged.AddAll(this)

	if estimators == nil {
		return merged
	}
	for _, estimator := range estimators {
		hll := estimator
		merged.AddAll(hll)
	}
	return merged
}

/*
 * Initial code from HyperLogLog.Builder.build()
 * by Scouter-Project	 */
func BuildHyperLogLog(bytes []byte) *HyperLogLog {
	in := io.NewDataInputX(bytes)
	log2m := uint32(in.ReadInt())
	n := in.ReadInt()
	ints := make([]uint32, n)
	for i := 0; i < int(n); i++ {
		ints[i] = uint32(in.ReadInt())
	}
	return NewHyperLogLog(log2m, NewRegisterSetInit(int(1<<log2m), ints))
}

func (this *HyperLogLog) getAlphaMM(p uint32, m uint32) float64 {
	// See the paper.
	switch p {
	case 4:
		return 0.673 * float64(m) * float64(m)
	case 5:
		return 0.697 * float64(m) * float64(m)
	case 6:
		return 0.709 * float64(m) * float64(m)
	default:
		return (0.7213 / (1 + 1.079/float64(m))) * float64(m) * float64(m)
	}
}

func linearCounting(m int, V float64) float64 {
	return float64(m) * math.Log(float64(m)/V)
}

func _log2m(rsd float64) uint32 {
	return uint32(math.Log((1.106/rsd)*(1.106/rsd)) / math.Log(2))
}

func rsd(log2m uint32) float64 {
	return 1.106 / math.Sqrt(math.Exp(float64(log2m)*math.Log(2)))
}

func validateLog2m(log2m uint32) bool {
	if log2m < 0 || log2m > 30 {
		//throw new IllegalArgumentException("log2m argument is " + log2m + " and is outpcodee the range [0, 30]");
		return false
	}
	return true

}

func Round(val float64) int64 {
	if val < 0 {
		return int64(val - 0.5)
	}
	return int64(val + 0.5)
}

var clzLookup = []uint8{
	32, 31, 30, 30, 29, 29, 29, 29, 28, 28, 28, 28, 28, 28, 28, 28,
}

// This optimized clz32 algorithm is from:
//
//	http://embeddedgurus.com/state-space/2014/09/
//			fast-deterministic-and-portable-counting-leading-zeros/
func clz32(x uint32) uint8 {
	var n uint8

	if x >= (1 << 16) {
		if x >= (1 << 24) {
			if x >= (1 << 28) {
				n = 28
			} else {
				n = 24
			}
		} else {
			if x >= (1 << 20) {
				n = 20
			} else {
				n = 16
			}
		}
	} else {
		if x >= (1 << 8) {
			if x >= (1 << 12) {
				n = 12
			} else {
				n = 8
			}
		} else {
			if x >= (1 << 4) {
				n = 4
			} else {
				n = 0
			}
		}
	}
	return clzLookup[x>>n] - n
}

func clz64(x uint64) uint8 {
	var c uint8
	for m := uint64(1 << 63); m&x == 0 && m != 0; m >>= 1 {
		c++
	}
	return c
}
