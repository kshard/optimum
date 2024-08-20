//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package optimum

import (
	"time"

	"github.com/fogfish/curie"
	"github.com/fogfish/schemaorg"
)

type Instances struct {
	Items []Instance `json:"items,omitempty"`
}

type Instance struct {
	ID      curie.IRI `json:"id"`
	Opts    string    `json:"opts"`
	Status  string    `json:"status"`
	Updated time.Time `json:"updated"`
	Version string    `json:"version"`
	Pending string    `json:"pending"`
}

type create struct {
	Name string         `json:"name"`
	Opts map[string]any `json:"opts"`
}

type Created struct {
	Version string        `json:"version,omitempty"`
	Job     schemaorg.Url `json:"job"`
}

type commit struct {
	Cursor string `json:"cursor"`
}

type Committed struct {
	Version string        `json:"version,omitempty"`
	Job     schemaorg.Url `json:"job"`
}

type JobStatus struct {
	Status  string `json:"status,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Created string `json:"created,omitempty"`
	Started string `json:"started,omitempty"`
	Stopped string `json:"stopped,omitempty"`
}

// Vector format
type Vector struct {
	UniqueKey []uint8   `json:"id,omitempty"`
	SortKey   []uint8   `json:"sk,omitempty"`
	Vec       []float32 `json:"v"`
}

type Query struct {
	K        int       `json:"k,omitempty"`
	EfSearch int       `json:"efSearch,omitempty"`
	Distance float32   `json:"distance,omitempty"`
	Query    []float32 `json:"query"`
}

type Result struct {
	Took    time.Duration `json:"took,omitempty"`
	Version Version       `json:"version,omitempty"`
	Hits    []Hit         `json:"hits,omitempty"`
}

type Hit struct {
	UniqueKey []uint8 `json:"key,omitempty"`
	SortKey   []uint8 `json:"sort,omitempty"`
	Rank      float32 `json:"rank"`
}

type Version struct {
	Cask    string `json:"cask"`
	Version string `json:"version"`
	Size    int    `json:"size"`
}
