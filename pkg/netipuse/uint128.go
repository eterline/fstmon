// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package netipuse // import "go4.org/netipx"

import (
	"encoding/binary"
	"math/bits"
	"net/netip"
)

type vec128 struct {
	hiBit uint64
	loBit uint64
}

func vec128From16(a [16]byte) vec128 {
	return vec128{
		binary.BigEndian.Uint64(a[:8]),
		binary.BigEndian.Uint64(a[8:]),
	}
}

func (u vec128) IPv6() netip.Addr {
	var a [16]byte
	binary.BigEndian.PutUint64(a[:8], u.hiBit)
	binary.BigEndian.PutUint64(a[8:], u.loBit)
	return netip.AddrFrom16(a)
}

func (u vec128) IPv4() netip.Addr {
	var a [8]byte
	binary.BigEndian.PutUint64(a[:], u.loBit)
	return netip.AddrFrom4([4]byte{a[4], a[5], a[6], a[7]})
}

// isZero reports whether u == 0.
//
// It's faster than u == (vec128{}) because the compiler (as of Go
// 1.15/1.16b1) doesn't do this trick and instead inserts a branch in
// its eq alg's generated code.
func (u vec128) isZero() bool { return u.hiBit|u.loBit == 0 }

// and returns the bitwise AND of u and m (u&m).
func (u vec128) and(m vec128) vec128 {
	return vec128{u.hiBit & m.hiBit, u.loBit & m.loBit}
}

// xor returns the bitwise XOR of u and m (u^m).
func (u vec128) xor(m vec128) vec128 {
	return vec128{u.hiBit ^ m.hiBit, u.loBit ^ m.loBit}
}

// or returns the bitwise OR of u and m (u|m).
func (u vec128) or(m vec128) vec128 {
	return vec128{u.hiBit | m.hiBit, u.loBit | m.loBit}
}

// not returns the bitwise NOT of u.
func (u vec128) not() vec128 {
	return vec128{^u.hiBit, ^u.loBit}
}

// subOne returns u - 1.
func (u vec128) subOne() vec128 {
	lo, borrow := bits.Sub64(u.loBit, 1, 0)
	return vec128{u.hiBit - borrow, lo}
}

// addOne returns u + 1.
func (u vec128) addOne() vec128 {
	lo, carry := bits.Add64(u.loBit, 1, 0)
	return vec128{u.hiBit + carry, lo}
}

func u64CommonPrefixLen(a, b uint64) uint8 {
	return uint8(bits.LeadingZeros64(a ^ b))
}

func (u vec128) commonPrefixLen(v vec128) (n uint8) {
	if n = u64CommonPrefixLen(u.hiBit, v.hiBit); n == 64 {
		n += u64CommonPrefixLen(u.loBit, v.loBit)
	}
	return
}

// func (u *vec128) halves() [2]*uint64 {
// 	return [2]*uint64{&u.hiBit, &u.loBit}
// }

// bitsSetFrom returns a copy of u with the given bit
// and all subsequent ones set.
func (u vec128) bitsSetFrom(bit uint8) vec128 {
	return u.or(mask6[bit].not())
}

// bitsClearedFrom returns a copy of u with the given bit
// and all subsequent ones cleared.
func (u vec128) bitsClearedFrom(bit uint8) vec128 {
	return u.and(mask6[bit])
}
