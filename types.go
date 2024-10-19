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

type Source struct {
	Cask    string `json:"cask"`
	Version string `json:"version"`
	Size    int    `json:"size"`
}
