/*
 * Copyright (c) 2021 Tobias Briones. All rights reserved.
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * This file is part of Example Project: Cryptosystems.
 *
 * This source code is licensed under the BSD-3-Clause License found in the
 * LICENSE file in the root directory of this source tree or at
 * https://opensource.org/licenses/BSD-3-Clause.
 */

// Package affine provides the data structures and implementation of the classical Affine cipher algorithm.
//
// Author: Tobias Briones
package affine

import (
	"errors"
	"github.com/tobiasbriones/ep-cryptosystems/algorithm"
	"strings"
)

type E struct {
	Pair     Pair
	alphabet algorithm.Alphabet
}

func (e *E) Set(pair Pair) error {
	if algorithm.Gcd(int(pair.A), 26) != 1 {
		return errors.New("A and 26 must be relatively prime")
	}
	e.Pair = pair
	e.alphabet = algorithm.GetAlphabet()
	return nil
}

func (e *E) Apply(x byte) byte {
	return (e.Pair.A*x + e.Pair.B) % 26
}

func (e *E) Inverse() D {
	return D{e.Pair, algorithm.GetAlphabet()}
}

type D struct {
	Pair     Pair
	alphabet algorithm.Alphabet
}

func (d *D) Apply(y byte) byte {
	var aInverse = modInverseIfCoprime(int(d.Pair.A), 26)
	var yInt = int(y)
	var bInt = int(d.Pair.B)
	return byte(norm(aInverse * (yInt - bInt) % 26))
}

type Pair struct {
	A byte
	B byte
}

func Encrypt(msg string, fn E) string {
	var enc = ""
	var input = strings.ToUpper(msg)
	input = strings.ReplaceAll(input, " ", "")

	for _, ch := range input {
		var x = byte(fn.alphabet.CanonicalPositionOf(byte(ch)))
		var y = fn.Apply(x)
		enc += string(fn.alphabet.Chars[y])
	}
	enc = strings.ToLower(enc)
	return enc
}

func Decrypt(enc string, fn D) string {
	var msg = ""
	var input = strings.ToUpper(enc)

	for _, ch := range input {
		var y = byte(fn.alphabet.CanonicalPositionOf(byte(ch)))
		var x = fn.Apply(y)
		msg += string(fn.alphabet.Chars[x])
	}
	msg = strings.ToLower(msg)
	return msg
}

func modInverseIfCoprime(a int, m int) int {
	var x int
	var y int
	gcd(a, m, &x, &y) // gcd must be 1 by hypothesis
	return (x%m + m) % m
}

func gcd(a int, b int, x *int, y *int) int {
	if a == 0 {
		*x = 0
		*y = 1
		return b
	}
	var x1 int
	var y1 int
	var gcd = gcd(b%a, a, &x1, &y1)
	*x = y1 - (b/a)*x1
	*y = x1
	return gcd
}

// Used to fix negative results, I took this function from the shift algorithm
// since I solved that problem there before
func norm(newPos int) int {
	if newPos < 0 {
		newPos *= -1
		newPos %= 26
		newPos = 26 - newPos
	}
	return newPos
}
