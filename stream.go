//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package optimum

import (
	"context"

	"github.com/fogfish/curie"
	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)

// Stream client
type Stream struct {
	http.Stack

	host ø.Authority
	cask curie.IRI

	chunk    int
	segement segment
}

type segment struct {
	Bag []Vector `json:"bag"`
}

func NewStream(stack http.Stack, host string, cask curie.IRI, chunk int) *Stream {
	return &Stream{
		Stack: stack,
		host:  ø.Authority(host),
		cask:  cask,

		chunk:    chunk,
		segement: segment{Bag: make([]Vector, 0)},
	}
}

// Write vector
func (api *Stream) Write(ctx context.Context, v Vector) error {
	api.segement.Bag = append(api.segement.Bag, v)

	if len(api.segement.Bag) >= api.chunk {
		return api.Sync(ctx)
	}

	return nil
}

// Commit vectors
func (api *Stream) Sync(ctx context.Context) (err error) {
	if len(api.segement.Bag) == 0 {
		return nil
	}

	return api.Stack.IO(ctx,
		http.PUT(
			ø.URI("%s/ds/%s/%s", api.host, curie.Prefix(api.cask), curie.Reference(api.cask)),
			ø.Accept.JSON,
			ø.ContentType.JSON,
			ø.Send(api.segement),

			ƒ.Status.Accepted,
			func(ctx *http.Context) error {
				api.segement = segment{Bag: make([]Vector, 0)}
				return nil
			},
		),
	)
}
