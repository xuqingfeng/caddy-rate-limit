// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

// +build !amd64

package poly1305

const (
	msgBlock   = uint32(1 << 24)
	finalBlock = uint32(0)
)

// Sum generates an authenticator for msg using a one-time key and puts the
// 16-byte result into out. Authenticating two different messages with the same
// key allows an attacker to forge messages at will.
func Sum(out *[TagSize]byte, msg []byte, key *[32]byte) {
	var (
		h, r  [5]uint32
		pad   [4]uint32
		block [TagSize]byte
		off   int
	)

	initialize(&r, &pad, key)

	n := len(msg) & (^(TagSize - 1))
	if n > 0 {
		core(msg[:n], msgBlock, &h, &r)
		msg = msg[n:]
	}

	if len(msg) > 0 {
		off += copy(block[:], msg)
		block[off] = 1
		core(block[:], finalBlock, &h, &r)
	}

	finalize(out, &h, &pad)
}

// New returns a hash.Hash computing the poly1305 sum.
// Notice that Poly1305 is inseure if one key is used twice.
func New(key *[32]byte) *Hash {
	p := new(Hash)
	initialize(&(p.r), &(p.pad), key)
	return p
}

// Hash implements a Poly1305 writer interface.
// Poly1305 cannot be used like common hash.Hash implementations,
// beause of using a Poly1305 key twice breaks its security.
// So poly1305.Hash does not support some kind of reset.
type Hash struct {
	h, r [5]uint32
	pad  [4]uint32

	buf  [TagSize]byte
	off  int
	done bool
}

// Write adds more data to the running Poly1305 hash.
// This function returns an non-nil error, if a call
// to Write happens after the hash's Sum function was
// called. So it's not possible to compute the checksum
// and than add more data.
func (p *Hash) Write(msg []byte) (int, error) {
	if p.done {
		return 0, errWriteAfterSum
	}
	n := len(msg)

	if p.off > 0 {
		dif := TagSize - p.off
		if n > dif {
			p.off += copy(p.buf[p.off:], msg[:dif])
			msg = msg[dif:]
			core(p.buf[:], msgBlock, &(p.h), &(p.r))
			p.off = 0
		} else {
			p.off += copy(p.buf[p.off:], msg)
			return n, nil
		}
	}

	if nn := len(msg) & (^(TagSize - 1)); nn > 0 {
		core(msg[:nn], msgBlock, &(p.h), &(p.r))
		msg = msg[nn:]
	}

	if len(msg) > 0 {
		p.off += copy(p.buf[p.off:], msg)
	}

	return n, nil
}

// Sum computes the Poly1305 checksum of the prevouisly
// processed data and writes it to out. It is legal to
// call this function more than one time.
func (p *Hash) Sum(out *[TagSize]byte) {
	h, r := p.h, p.r
	pad := p.pad

	if p.off > 0 {
		var buf [TagSize]byte
		copy(buf[:], p.buf[:p.off])
		buf[p.off] = 1 // invariant: p.off < TagSize

		core(buf[:], finalBlock, &h, &r)
	}

	finalize(out, &h, &pad)
	p.done = true
}

func initialize(r *[5]uint32, pad *[4]uint32, key *[32]byte) {
	r[0] = (uint32(key[0]) | uint32(key[1])<<8 | uint32(key[2])<<16 | uint32(key[3])<<24) & 0x3ffffff
	r[1] = ((uint32(key[3]) | uint32(key[4])<<8 | uint32(key[5])<<16 | uint32(key[6])<<24) >> 2) & 0x3ffff03
	r[2] = ((uint32(key[6]) | uint32(key[7])<<8 | uint32(key[8])<<16 | uint32(key[9])<<24) >> 4) & 0x3ffc0ff
	r[3] = ((uint32(key[9]) | uint32(key[10])<<8 | uint32(key[11])<<16 | uint32(key[12])<<24) >> 6) & 0x3f03fff
	r[4] = ((uint32(key[12]) | uint32(key[13])<<8 | uint32(key[14])<<16 | uint32(key[15])<<24) >> 8) & 0x00fffff

	pad[0] = (uint32(key[16]) | uint32(key[17])<<8 | uint32(key[18])<<16 | uint32(key[19])<<24)
	pad[1] = (uint32(key[20]) | uint32(key[21])<<8 | uint32(key[22])<<16 | uint32(key[23])<<24)
	pad[2] = (uint32(key[24]) | uint32(key[25])<<8 | uint32(key[26])<<16 | uint32(key[27])<<24)
	pad[3] = (uint32(key[28]) | uint32(key[29])<<8 | uint32(key[30])<<16 | uint32(key[31])<<24)
}

func core(msg []byte, flag uint32, h, r *[5]uint32) {
	h0, h1, h2, h3, h4 := h[0], h[1], h[2], h[3], h[4]
	r0, r1, r2, r3, r4 := uint64(r[0]), uint64(r[1]), uint64(r[2]), uint64(r[3]), uint64(r[4])
	s1, s2, s3, s4 := r1*5, r2*5, r3*5, r4*5

	var d0, d1, d2, d3, d4 uint64

	length := len(msg) & (^(TagSize - 1))
	for i := 0; i < length; i += TagSize {
		// h += m
		h0 += (uint32(msg[i+0]) | uint32(msg[i+1])<<8 | uint32(msg[i+2])<<16 | uint32(msg[i+3])<<24) & 0x3ffffff
		h1 += ((uint32(msg[i+3]) | uint32(msg[i+4])<<8 | uint32(msg[i+5])<<16 | uint32(msg[i+6])<<24) >> 2) & 0x3ffffff
		h2 += ((uint32(msg[i+6]) | uint32(msg[i+7])<<8 | uint32(msg[i+8])<<16 | uint32(msg[i+9])<<24) >> 4) & 0x3ffffff
		h3 += ((uint32(msg[i+9]) | uint32(msg[i+10])<<8 | uint32(msg[i+11])<<16 | uint32(msg[i+12])<<24) >> 6) & 0x3ffffff
		h4 += ((uint32(msg[i+12]) | uint32(msg[i+13])<<8 | uint32(msg[i+14])<<16 | uint32(msg[i+15])<<24) >> 8) | flag

		// h *= r
		d0 = (uint64(h0) * r0) + (uint64(h1) * s4) + (uint64(h2) * s3) + (uint64(h3) * s2) + (uint64(h4) * s1)
		d1 = (d0 >> 26) + (uint64(h0) * r1) + (uint64(h1) * r0) + (uint64(h2) * s4) + (uint64(h3) * s3) + (uint64(h4) * s2)
		d2 = (d1 >> 26) + (uint64(h0) * r2) + (uint64(h1) * r1) + (uint64(h2) * r0) + (uint64(h3) * s4) + (uint64(h4) * s3)
		d3 = (d2 >> 26) + (uint64(h0) * r3) + (uint64(h1) * r2) + (uint64(h2) * r1) + (uint64(h3) * r0) + (uint64(h4) * s4)
		d4 = (d3 >> 26) + (uint64(h0) * r4) + (uint64(h1) * r3) + (uint64(h2) * r2) + (uint64(h3) * r1) + (uint64(h4) * r0)

		// h %= p
		h0 = uint32(d0) & 0x3ffffff
		h1 = uint32(d1) & 0x3ffffff
		h2 = uint32(d2) & 0x3ffffff
		h3 = uint32(d3) & 0x3ffffff
		h4 = uint32(d4) & 0x3ffffff

		h0 += uint32(d4>>26) * 5
		h1 += h0 >> 26
		h0 = h0 & 0x3ffffff
	}
	h[0], h[1], h[2], h[3], h[4] = h0, h1, h2, h3, h4
}

func finalize(tag *[TagSize]byte, h *[5]uint32, pad *[4]uint32) {
	var g0, g1, g2, g3, g4 uint32

	// fully carry h
	h0, h1, h2, h3, h4 := h[0], h[1], h[2], h[3], h[4]

	h2 += h1 >> 26
	h1 &= 0x3ffffff
	h3 += h2 >> 26
	h2 &= 0x3ffffff
	h4 += h3 >> 26
	h3 &= 0x3ffffff
	h0 += 5 * (h4 >> 26)
	h4 &= 0x3ffffff
	h1 += h0 >> 26
	h0 &= 0x3ffffff

	// h + -p
	g0 = h0 + 5

	g1 = h1 + (g0 >> 26)
	g0 &= 0x3ffffff
	g2 = h2 + (g1 >> 26)
	g1 &= 0x3ffffff
	g3 = h3 + (g2 >> 26)
	g2 &= 0x3ffffff
	g4 = h4 + (g3 >> 26) - (1 << 26)
	g3 &= 0x3ffffff

	// select h if h < p else h + -p
	mask := (g4 >> (32 - 1)) - 1
	g0 &= mask
	g1 &= mask
	g2 &= mask
	g3 &= mask
	g4 &= mask
	mask = ^mask
	h0 = (h0 & mask) | g0
	h1 = (h1 & mask) | g1
	h2 = (h2 & mask) | g2
	h3 = (h3 & mask) | g3
	h4 = (h4 & mask) | g4

	// h %= 2^128
	h0 |= h1 << 26
	h1 = ((h1 >> 6) | (h2 << 20))
	h2 = ((h2 >> 12) | (h3 << 14))
	h3 = ((h3 >> 18) | (h4 << 8))

	// tag = (h + pad) % (2^128)
	f := uint64(h0) + uint64(pad[0])
	h0 = uint32(f)
	f = uint64(h1) + uint64(pad[1]) + (f >> 32)
	h1 = uint32(f)
	f = uint64(h2) + uint64(pad[2]) + (f >> 32)
	h2 = uint32(f)
	f = uint64(h3) + uint64(pad[3]) + (f >> 32)
	h3 = uint32(f)

	tag[0] = byte(h0)
	tag[1] = byte(h0 >> 8)
	tag[2] = byte(h0 >> 16)
	tag[3] = byte(h0 >> 24)
	tag[4] = byte(h1)
	tag[5] = byte(h1 >> 8)
	tag[6] = byte(h1 >> 16)
	tag[7] = byte(h1 >> 24)
	tag[8] = byte(h2)
	tag[9] = byte(h2 >> 8)
	tag[10] = byte(h2 >> 16)
	tag[11] = byte(h2 >> 24)
	tag[12] = byte(h3)
	tag[13] = byte(h3 >> 8)
	tag[14] = byte(h3 >> 16)
	tag[15] = byte(h3 >> 24)
}
