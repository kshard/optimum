//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package surface

import (
	"time"

	"github.com/kshard/optimum"
)

// Vector representing point in k-dimension space
type Vector struct {
	UniqueKey []uint8   `json:"id,omitempty"`
	SortKey   []uint8   `json:"sk,omitempty"`
	Vector    []float32 `json:"v"`
}

// Query points in k-dimension space
type Query struct {
	K        int       `json:"k,omitempty"`
	EfSearch int       `json:"efSearch,omitempty"`
	Distance float32   `json:"distance,omitempty"`
	Query    []float32 `json:"query"`
}

// Results from query
type Result struct {
	Took   time.Duration  `json:"took,omitempty"`
	Source optimum.Source `json:"source,omitempty"`
	Hits   []Hit          `json:"hits,omitempty"`
}

type Hit struct {
	UniqueKey []uint8 `json:"key,omitempty"`
	SortKey   []uint8 `json:"sort,omitempty"`
	Rank      float32 `json:"rank"`
}
