//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package encoding

import (
	"bufio"
	"encoding/hex"
	"io"
	"strconv"
	"strings"
)

type Scanner struct {
	r         *bufio.Scanner
	err       error
	uniqueKey []byte
	vector    []float32
}

func New(r io.Reader) *Scanner {
	return &Scanner{
		r: bufio.NewScanner(r),
	}
}

func (s *Scanner) Err() error        { return s.r.Err() }
func (s *Scanner) UniqueKey() []byte { return s.uniqueKey }
func (s *Scanner) Vector() []float32 { return s.vector }

func (s *Scanner) Scan() bool {
	if !s.r.Scan() {
		return false
	}

	seq := strings.Split(s.r.Text(), " ")

	f32 := make([]float32, len(seq)-1)
	for i := 1; i < len(seq); i++ {
		v, err := strconv.ParseFloat(seq[i], 32)
		if err != nil {
			s.err = err
			return false
		}
		f32[i-1] = float32(v)
	}
	s.vector = f32

	if strings.HasPrefix(seq[0], "0x") {
		x, err := hex.DecodeString(seq[0][2:])
		if err != nil {
			s.err = err
			return false
		}
		s.uniqueKey = x
	} else {
		s.uniqueKey = []byte(seq[0])
	}

	return true
}
