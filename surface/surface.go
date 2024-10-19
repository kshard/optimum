//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package surface

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/fogfish/curie"
	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
	"github.com/kshard/wreck"
)

// Client for reading/writing Graph-based Nearest Neighbor Surface.
type Client struct {
	http.Stack
	host ø.Authority
}

// Creates the client for reading/writing Graph-based Nearest Neighbor Surface.
func New(stack http.Stack, host string) *Client {
	return &Client{
		Stack: stack,
		host:  ø.Authority(host),
	}
}

// Write vector(s)
func (api *Client) Write(ctx context.Context, cask curie.IRI, bag []Vector) error {
	if len(bag) == 0 {
		return nil
	}

	var buf bytes.Buffer
	out := wreck.NewWriterJSON(&buf, false)
	seq := wreck.NewWriter[float32](out)

	for _, vec := range bag {
		if err := seq.Write(vec.UniqueKey, vec.SortKey, vec.Vector); err != nil {
			return err
		}
	}

	if err := out.Close(); err != nil {
		return err
	}

	return api.IO(ctx,
		http.POST(
			ø.URI("%s/ds/%s/%s/object", api.host, curie.Prefix(cask), curie.Reference(cask)),
			ø.Accept.JSON,
			ø.ContentType.JSON,
			ø.Send(struct {
				V json.RawMessage `json:"object"`
			}{
				V: buf.Bytes(),
			}),

			ƒ.Status.Accepted,
		),
	)
}

// Query nearest neighbor points to the given vector
func (api *Client) Query(ctx context.Context, cask curie.IRI, q Query) (*Result, error) {
	return http.IO[Result](
		api.WithContext(ctx),
		http.GET(
			ø.URI("%s/ds/%s/%s", api.host, curie.Prefix(cask), curie.Reference(cask)),
			ø.Accept.JSON,
			ø.ContentType.JSON,
			ø.Send(struct {
				Q Query `json:"query"`
			}{
				Q: q,
			}),

			ƒ.Status.OK,
		),
	)
}
