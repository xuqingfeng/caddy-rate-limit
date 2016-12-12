// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

// +build amd64, !gccgo, !appengine

package poly1305

// Sum generates an authenticator for msg using a one-time key and puts the
// 16-byte result into out. Authenticating two different messages with the same
// key allows an attacker to forge messages at will.
func Sum(out *[TagSize]byte, msg []byte, key *[32]byte) {
	if len(msg) == 0 {
		msg = []byte{}
	}

	var state [7]uint64 // := uint64{ h0, h1, h2, r0, r1, pad0, pad1 }
	initialize(&state, key)
	core(&state, msg)
	finalize(out, &state)
}

// New returns a hash.Hash computing the poly1305 sum.
// Notice that Poly1305 is inseure if one key is used twice.
func New(key *[32]byte) *Hash {
	p := new(Hash)
	initialize(&(p.state), key)
	return p
}

// Hash implements a Poly1305 writer interface.
// Poly1305 cannot be used like common hash.Hash implementations,
// beause of using a Poly1305 key twice breaks its security.
// So poly1305.Hash does not support some kind of reset.
type Hash struct {
	state [7]uint64 // := uint64{ h0, h1, h2, r0, r1, pad0, pad1 }

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
			core(&(p.state), p.buf[:])
			p.off = 0
		} else {
			p.off += copy(p.buf[p.off:], msg)
			return n, nil
		}
	}

	if nn := len(msg) & (^(TagSize - 1)); nn > 0 {
		core(&(p.state), msg[:nn])
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
	state := p.state

	if p.off > 0 {
		core(&state, p.buf[:p.off])
	}

	finalize(out, &state)
	p.done = true
}

//go:noescape
func initialize(state *[7]uint64, key *[32]byte)

//go:noescape
func core(state *[7]uint64, msg []byte)

//go:noescape
func finalize(tag *[TagSize]byte, state *[7]uint64)
