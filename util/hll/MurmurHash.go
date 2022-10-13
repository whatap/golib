/*
 * This file from
 *  https://github.com/addthis/stream-lib/blob/master/src/main/java/com/clearspring/analytics/hash/MurmurHash.java
 *
 *    This class modified by Scouter-Project *   - original  package :  com.clearspring.analytics.hash
 *
 */
/**
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements. See the NOTICE file distributed with this
 * work for additional information regarding copyright ownership. The ASF
 * licenses this file to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 *
 */
/**
 * This is a very fast, non-cryptographic hash suitable for general hash-based
 * lookup. See http://murmurhash.googlepages.com/ for more details.
 * <p/>
 * <p>
 * The C version of MurmurHash 2.0 found at that site was ported to Java by
 * Andrzej Bialecki (ab at getopt org).
 * </p>
 */
package hll

import ()

//func MurmurHash(Object o) int {
//        if (o == null) {
//            return 0;
//        }
//        if (o instanceof Long) {
//            return hashLong((Long) o);
//        }
//        if (o instanceof Integer) {
//            return hashLong((Integer) o);
//        }
//        if (o instanceof Double) {
//            return hashLong(Double.doubleToRawLongBits((Double) o));
//        }
//        if (o instanceof Float) {
//            return hashLong(Float.floatToRawIntBits((Float) o));
//        }
//        if (o instanceof String) {
//            return hash(((String) o).getBytes());
//        }
//        if (o instanceof byte[]) {
//            return hash((byte[]) o);
//        }
//
//        return hash(o.toString());
//    }
func MurmurHash(o uint32) uint32 {
	return MurmurHashLong(uint64(o))
}

func MurmurHashByte(data []byte) uint32 {
	return murmurHash(data, int32(len(data)), 0xe17a1465)
}
func MurmurHashByteSeed(data []byte, seed uint32) uint32 {
	return murmurHash(data, int32(len(data)), seed)
}
func murmurHash(data []byte, length int32, seed uint32) uint32 {
	m := uint32(0x5bd1e995)
	r := uint32(24)
	h := seed ^ uint32(length)
	len_4 := uint32(length) >> 2

	for i := 0; i < int(len_4); i++ {
		i_4 := i << 2
		k := uint32(data[i_4+3])
		k = k << 8
		k = k | uint32(data[i_4+2]) & 0xff
		k = k << 8
		k = k | uint32(data[i_4+1]) & 0xff
		k = k << 8
		k = k | uint32(data[i_4+0]) & 0xff
		k *= m
		//k ^= k >>> r
		k ^= k >> r
		k *= m
		h *= m
		h ^= k
	}
	// avoid calculating modulo
	len_m := len_4 << 2
	left := uint32(length) - len_m
	if left != 0 {
		if left >= 3 {
			h ^= uint32(data[length-3]) << 16
		}
		if left >= 2 {
			h ^= uint32(data[length-2]) << 8
		}
		if left >= 1 {
			h ^= uint32(data[length-1])
		}
		h *= m
	}
	h ^= h >> 13
	h *= m
	h ^= h >> 15

	return h
}

func MurmurHashLong(data uint64) uint32 {
	m := uint32(0x5bd1e995)
	r := uint32(24)
	h := uint32(0)
	k := uint32(data * uint64(m))
	//k ^= k >>> r
	k ^= k >> r
	h ^= k * m
	k = uint32((data >> 32) * uint64(m))
	//k ^= k >>> r;
	k ^= k >> r
	h *= m
	h ^= k * m
	//h ^= h >>> 13;
	h ^= h >> 13
	h *= m
	//h ^= h >>> 15;
	h ^= h >> 15

	return h
}

//func MurmurHash64(Object o) int64 {
//        if (o == null) {
//            return 0l;
//        } else if (o instanceof String) {
//            final byte[] bytes = ((String) o).getBytes();
//            return hash64(bytes, bytes.length);
//        } else if (o instanceof byte[]) {
//            final byte[] bytes = (byte[]) o;
//            return hash64(bytes, bytes.length);
//        }
//        return hash64(o.toString());
//    }
// 64 bit implementation copied from here:  https://github.com/tnm/murmurhash-java
/**
 * Generates 64 bit hash from byte array with default seed value.
 *
 * @param data   byte array to hash
 * @param length length of the array to hash
 * @return 64 bit hash of the given string
 */
func MurmurHashLongByte(data []byte, length int32) uint64 {
	return murmurHashLong(data, length, 0xe17a1465)
}

/**
 * Generates 64 bit hash from byte array of the given length and seed.
 *
 * @param data   byte array to hash
 * @param length length of the array to hash
 * @param seed   initial seed value
 * @return 64 bit hash of the given array
 */
func murmurHashLong(data []byte, length int32, seed uint32) uint64 {
	//m := 0xc6a4a7935bd1e995L
	m := uint64(0xc6a4a7935bd1e995)
	r := uint64(47)
	//h := (seed & 0xffffffffl) ^ (length * m)
	h := uint64(seed & 0xffffffff) ^ (uint64(length) * m)
	length8 := length / 8
	for i := 0; i < int(length8); i++ {
		i8 := i * 8
		k := (uint64(data[i8+0]) & 0xff) + ((uint64(data[i8+1]) & 0xff) << 8) +
		((uint64(data[i8+2]) & 0xff) << 16) + ((uint64(data[i8+3]) & 0xff) << 24) +
		((uint64(data[i8+4]) & 0xff) << 32) + ((uint64(data[i8+5]) & 0xff) << 40) +
		((uint64(data[i8+6]) & 0xff) << 48) + ((uint64(data[i8+7]) & 0xff) << 56)
		k *= m
		//k ^= k >>> r
		k ^= k >> r
		k *= m
		h ^= k
		h *= m
	}
	switch length % 8 {
	case 7:
		//h ^= (long) (data[(length & ~7) + 6] & 0xff) << 48;
		h ^= (uint64(data[(length & ^7)+6])&0xff) << 48
		fallthrough
	case 6:
		h ^= uint64(data[(length & ^7)+5]&0xff) << 40
		fallthrough
	case 5:
		h ^= uint64(data[(length & ^7)+4]&0xff) << 32
		fallthrough
	case 4:
		h ^= uint64(data[(length & ^7)+3]&0xff) << 24
		fallthrough
	case 3:
		h ^= uint64(data[(length & ^7)+2]&0xff) << 16
		fallthrough
	case 2:
		h ^= uint64(data[(length & ^7)+1]&0xff) << 8
		fallthrough
	case 1:
		h ^= uint64(data[length & ^7] & 0xff)
		h *= m
	}
	//h ^= h >>> r;
	h ^= h >> r
	h *= m
	//h ^= h >>> r;
	h ^= h >> r
	return h
}
